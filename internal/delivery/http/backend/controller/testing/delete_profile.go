package testing

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (ctrl *Controller) DeleteProfileHandler(c *fiber.Ctx) error {
	const op = "DeleteProfileHandler"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	err = ctrl.testingFacade.Cross.DeleteProfile(
		c.Context(),
		id,
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
