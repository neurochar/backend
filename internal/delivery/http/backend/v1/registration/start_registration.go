package registration

import (
	"fmt"
	"net"

	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	v1 "github.com/neurochar/backend/internal/delivery/http/backend/v1"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/internal/delivery/http/middleware"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/validation"
)

type StartRegistrationHandlerIn struct {
	Email string `json:"email"`
}

func (ctrl *Controller) StartRegistrationHandler(c *fiber.Ctx) error {
	const op = "StartRegistrationHandler"

	in := &StartRegistrationHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(
			httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...),
			"%s.%s", ctrl.pkg, op,
		)
	}

	ip := middleware.GetRealIP(c)

	backoffSession, ok := ctrl.backoff.GetIfExists(ip, v1.BackoffDefaultGroupID)

	if ok && !backoffSession.IsAllowed() {
		tryAfter := backoffSession.NextAllowedUntilSeconds()

		c.Set("Retry-After", fmt.Sprintf("%d", tryAfter))

		return appErrors.Chainf(
			appErrors.ErrBackoff.WithDetail("try_after_sec", false, tryAfter),
			"%s.%s", ctrl.pkg, op,
		)
	}

	requestIP := net.ParseIP(ip)

	_, err := ctrl.tenantFacade.Registration.CreateByDTO(
		c.Context(),
		tenantUC.CreateRegistrationIn{
			Email:     in.Email,
			RequestIP: requestIP,
		},
		&requestIP,
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	backoffSession = ctrl.backoff.Get(ip, v1.BackoffDefaultGroupID)
	backoffSession.AddBackoff()

	return nil
}
