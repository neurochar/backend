package auth

import (
	"errors"
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
	Email        string `json:"email" validate:"required,email"`
	TenantTextID string `json:"tenantTextID" validate:"required"`
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

	tenant, err := ctrl.tenantFacade.Tenant.FindOneByTextID(c.Context(), in.TenantTextID, nil)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return appErrors.Chainf(appErrors.ErrBadRequest.WithHints("tenant not found"), "%s.%s", ctrl.pkg, op)
		}
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
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

	code, err := ctrl.tenantUserFacade.Account.RequestPasswordRecoveryByEmail(c.Context(), tenant.ID, in.Email, requestIP)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return c.JSON(RequestPasswordRecoveryHandlerOut{
		CodeID: code.ID,
	})
}
