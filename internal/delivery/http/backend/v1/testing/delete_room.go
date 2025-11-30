package testing

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (ctrl *Controller) DeleteRoomHandler(c *fiber.Ctx) error {
	const op = "DeleteRoomHandler"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	err = ctrl.testingFacade.Cross.DeleteRoom(
		c.Context(),
		id,
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
