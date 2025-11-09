package users

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (ctrl *Controller) GetRoleHandler(c *fiber.Ctx) error {
	const op = "GetRoleHandler"

	id, err := c.ParamsInt("id", 0)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	role, err := ctrl.userFacade.Role.GetRoleByID(c.Context(), uint64(id))
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out := OutAccountRole{
		Version:  role.Role.Version(),
		ID:       role.Role.ID,
		Name:     role.Role.Name,
		IsSuper:  role.Role.IsSuper,
		IsSystem: role.Role.IsSystem,
		Rights:   make(map[string]OutAccountRoleRight, len(role.Rights)),
	}

	if role.Rights != nil {
		for _, right := range role.Rights {
			out.Rights[right.Right.Key] = OutAccountRoleRight{
				ID:    right.Right.ID,
				Key:   right.Right.Key,
				Type:  right.Right.Type.String(),
				Value: right.Value,
			}
		}
	}

	return c.JSON(out)
}
