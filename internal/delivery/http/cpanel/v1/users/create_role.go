package users

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	userUC "github.com/neurochar/backend/internal/domain/user/usecase"
	"github.com/neurochar/backend/pkg/validation"
)

type CreateRoleHandlerIn struct {
	InAccountRole
}

func (ctrl *Controller) CreateRoleHandler(c *fiber.Ctx) error {
	const op = "CreateRoleHandler"

	in := &CreateRoleHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(
			httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...),
			"%s.%s", ctrl.pkg, op,
		)
	}

	role, err := ctrl.userFacade.Role.CreateRole(c.Context(), userUC.CreateRoleInput{
		Name:   in.Name,
		Rights: in.Rights,
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out := OutAccountRole{
		Version: role.Role.Version(),

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
