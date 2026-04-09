package rooms

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/pkg/auth"
)

func (ctrl *Controller) GetRoomHandler(c *fiber.Ctx) error {
	const op = "GetRoomHandler"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	roomDTO, err := ctrl.testingFacade.Room.FindOneByID(auth.WithoutCheckTenantAccess(c.Context()), id, nil, nil)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	tenant, err := ctrl.tenantFacade.Tenant.FindOneByID(auth.WithoutCheckTenantAccess(c.Context()), roomDTO.Room.TenantID, nil)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out, err := OutRoomDTO(
		c,
		roomDTO,
		tenant,
	)
	if err != nil {
		return err
	}

	return c.JSON(out)
}
