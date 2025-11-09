package users

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/internal/delivery/http/middleware"
	"github.com/neurochar/backend/pkg/validation"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

const authCookie = "auth_admin_session"

type LoginHandlerIn struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"min=0"`
}

// LoginHandler - auth handler
func (ctrl *Controller) LoginHandler(c *fiber.Ctx) error {
	const op = "LoginHandler"

	in := &LoginHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...), "%s.%s", ctrl.pkg, op)
	}

	ip := middleware.GetRealIP(c)

	backoffSession, ok := ctrl.backoff.GetIfExists(ip, backoffConfigAuthGroupID)

	if ok && !backoffSession.IsAllowed() {
		tryAfter := backoffSession.NextAllowedUntilSeconds()

		c.Set("Retry-After", fmt.Sprintf("%d", tryAfter))

		return appErrors.Chainf(
			appErrors.ErrBackoff.WithDetail("try_after_sec", false, tryAfter),
			"%s.%s", ctrl.pkg, op,
		)
	}

	requestIP := net.ParseIP(ip)

	session, _, role, err := ctrl.userFacade.AdminAuth.LoginByEmail(c.Context(), in.Email, in.Password, requestIP)
	if err != nil {
		if errors.Is(err, appErrors.ErrUnauthorized) {
			backoffSession = ctrl.backoff.Get(ip, backoffConfigAuthGroupID)
			backoffSession.AddCounter()
			if backoffSession.Counter() > 1 {
				if backoffSession.AddBackoff() {
					c.Set("Retry-After", fmt.Sprintf("%d", backoffSession.NextAllowedUntilSeconds()))
				}
			}
		}

		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	jwt, err := ctrl.userFacade.AdminAuth.SessionToJWT(session)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	user, err := ctrl.userFacade.Common.FindOneByAccountID(c.Context(), session.AccountID)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	cookie := fiber.Cookie{
		Name:     authCookie,
		Value:    jwt,
		Path:     fmt.Sprintf("/%s/v1", ctrl.cfg.CPanelApp.HTTP.Prefix),
		HTTPOnly: true,
		Secure:   ctrl.cfg.CPanelApp.Base.IsProd,
		SameSite: "lax",
		MaxAge:   3600 * 24 * 30,
		Expires:  time.Now().Add(time.Second * 3600 * 24 * 30),
	}

	c.Cookie(&cookie)

	out, err := OutUserDTO(c, ctrl.fileUC, user.Account, user.ProfileDTO, role, true)
	if err != nil {
		return err
	}

	return c.JSON(out)
}
