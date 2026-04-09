package auth

import (
	"errors"
	"fmt"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/internal/delivery/http/middleware"
	"github.com/neurochar/backend/pkg/validation"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

type LoginHandlerIn struct {
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"min=0"`
	TenantTextID string `json:"tenantTextID" validate:"required"`
}

func (ctrl *Controller) LoginHandler(c *fiber.Ctx) error {
	const op = "LoginHandler"

	in := &LoginHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...), "%s.%s", ctrl.pkg, op)
	}

	tenant, err := ctrl.tenantFacade.Tenant.FindOneByTextID(c.Context(), in.TenantTextID, nil)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return appErrors.Chainf(appErrors.ErrBadRequest.WithHints("tenant not found"), "%s.%s", ctrl.pkg, op)
		}
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
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

	authDTO, err := ctrl.tenantFacade.Auth.LoginByEmail(c.Context(), tenant.ID, in.Email, in.Password, requestIP)
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

	accessJWT, err := ctrl.tenantFacade.Auth.IssueAccessJWT(authDTO.AccessClaims)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	refreshJWT, err := ctrl.tenantFacade.Auth.IssueRefreshJWT(authDTO.RefreshClaims)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out, err := OutLoginDTO(
		c,
		ctrl.fileUC,
		authDTO.AccountDTO,
		refreshJWT,
		uint64(ctrl.cfg.Auth.RefreshTokenLifetimeHrs)*3600,
		accessJWT,
	)
	if err != nil {
		return err
	}

	return c.JSON(out)
}
