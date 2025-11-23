package tenants

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (ctrl *Controller) GetIsExistsTenantHandler(c *fiber.Ctx) error {
	const op = "GetIsExistsTenantHandler"

	textIDStr := c.Params("text_id")

	_, err := ctrl.tenantFacade.Tenant.FindOneByTextID(c.Context(), textIDStr, nil)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
