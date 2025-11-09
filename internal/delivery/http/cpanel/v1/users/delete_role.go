package users

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (ctrl *Controller) DeleteRoleHandler(c *fiber.Ctx) error {
	const op = "DeleteRoleHandler"

	id, err := c.ParamsInt("id", 0)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	err = ctrl.userFacade.Common.DeleteRole(c.Context(), uint64(id))
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
