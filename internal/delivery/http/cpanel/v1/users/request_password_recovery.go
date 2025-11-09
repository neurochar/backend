package users

import (
	"fmt"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/internal/delivery/http/middleware"
	"github.com/neurochar/backend/pkg/validation"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

type RequestPasswordRecoveryHandlerIn struct {
	Email string `json:"email" validate:"required,email"`
}

type RequestPasswordRecoveryHandlerOut struct {
	CodeID uuid.UUID `json:"codeID"`
}

func (ctrl *Controller) RequestPasswordRecoveryHandler(c *fiber.Ctx) error {
	const op = "RequestPasswordRecoveryHandler"

	in := &RequestPasswordRecoveryHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...), "%s.%s", ctrl.pkg, op)
	}

	ip := middleware.GetRealIP(c)

	backoffSession, ok := ctrl.backoff.GetIfExists(ip, backoffConfigPasswordRecoveryGroupID)

	if ok && !backoffSession.IsAllowed() {
		tryAfter := backoffSession.NextAllowedUntilSeconds()

		c.Set("Retry-After", fmt.Sprintf("%d", tryAfter))

		return appErrors.Chainf(
			appErrors.ErrBackoff.WithDetail("try_after_sec", false, tryAfter),
			"%s.%s", ctrl.pkg, op,
		)
	}

	if !ok {
		backoffSession = ctrl.backoff.Get(ip, backoffConfigPasswordRecoveryGroupID)
	}

	requestIP := net.ParseIP(ip)

	_ = backoffSession.AddBackoff()
	c.Set("Retry-After", fmt.Sprintf("%d", backoffSession.NextAllowedUntilSeconds()))

	code, err := ctrl.userFacade.Account.RequestPasswordRecoveryByEmail(c.Context(), in.Email, requestIP)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return c.JSON(RequestPasswordRecoveryHandlerOut{
		CodeID: code.ID,
	})
}
