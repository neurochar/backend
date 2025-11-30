package rooms

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (ctrl *Controller) GetRoomHandler(c *fiber.Ctx) error {
	const op = "GetRoomHandler"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	profileDTO, err := ctrl.testingFacade.Room.FindOneByID(c.Context(), id, nil, nil)
	if err != nil {
		if errors.Is(err, appErrors.ErrForbidden) {
			return appErrors.Chainf(appErrors.ErrNotFound, "%s.%s", ctrl.pkg, op)
		}
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out, err := OutRoomDTO(
		c,
		profileDTO,
	)
	if err != nil {
		return err
	}

	return c.JSON(out)
}
