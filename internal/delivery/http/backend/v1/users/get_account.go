package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	backendMiddleware "github.com/neurochar/backend/internal/delivery/http/backend/middleware"
)

func (ctrl *Controller) GetAccountHandler(c *fiber.Ctx) error {
	const op = "GetAccountHandler"

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	auth := backendMiddleware.GetAuthData(c)
	if auth == nil {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	isConfirmed, err := ctrl.tenantUserFacade.Auth.IsSessionConfirmed(c.Context(), auth.SessionID)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	if !isConfirmed {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	accountDTO, err := ctrl.tenantUserFacade.Account.FindOneByID(c.Context(), id, nil, nil)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrInternal.WithParent(err), "%s.%s", ctrl.pkg, op)
	}

	if accountDTO.Account.TenantID != auth.TenantID {
		return appErrors.Chainf(appErrors.ErrForbidden, "%s.%s", ctrl.pkg, op)
	}

	out, err := OutAccountDTO(
		c,
		true,
		ctrl.fileUC,
		accountDTO,
	)
	if err != nil {
		return err
	}

	return c.JSON(out)
}
