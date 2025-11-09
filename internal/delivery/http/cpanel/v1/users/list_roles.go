package users

import (
	"github.com/gofiber/fiber/v2"
)

type ListRolesHandlerOut struct {
	Items []OutAccountRole `json:"items"`
}

func (ctrl *Controller) ListRolesHandler(c *fiber.Ctx) error {
	roles := ctrl.userFacade.Role.GetRolesList(c.Context())

	out := ListRolesHandlerOut{
		Items: make([]OutAccountRole, 0, len(roles)),
	}

	for _, role := range roles {
		outItem := OutAccountRole{
			Version:  role.Role.Version(),
			ID:       role.Role.ID,
			Name:     role.Role.Name,
			IsSuper:  role.Role.IsSuper,
			IsSystem: role.Role.IsSystem,
		}

		out.Items = append(out.Items, outItem)
	}

	return c.JSON(out)
}
