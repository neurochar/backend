package users

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (ctrl *Controller) GetUserHandler(c *fiber.Ctx) error {
	const op = "GetUserHandler"

	id, err := c.ParamsInt("profile_id", 0)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	item, err := ctrl.userFacade.Common.FindOneByProfileID(c.Context(), uint64(id))
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out, err := OutUserDTO(c, ctrl.fileUC, item.Account, item.ProfileDTO, item.Role, false)
	if err != nil {
		return err
	}

	return c.JSON(out)
}
