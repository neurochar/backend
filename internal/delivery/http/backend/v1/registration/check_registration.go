package registration

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (ctrl *Controller) CheckRegistrationHandler(c *fiber.Ctx) error {
	const op = "CheckRegistrationHandler"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	registration, err := ctrl.tenantFacade.Registration.FindOneByID(c.Context(), id, nil)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	if registration.IsFinished {
		return appErrors.Chainf(appErrors.ErrNotFound, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
