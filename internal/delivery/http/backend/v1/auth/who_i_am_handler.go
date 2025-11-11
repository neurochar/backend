package auth

import (
	"github.com/gofiber/fiber/v2"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	tenantUserUC "github.com/neurochar/backend/internal/domain/tenant_user/usecase"
)

func (ctrl *Controller) WhoIAmHandler(c *fiber.Ctx) error {
	const op = "WhoIAmHandler"

	auth := middleware.GetAuthData(c)
	if auth == nil {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	isRevoked, err := ctrl.tenantUserFacade.Auth.IsSessionRevoked(c.Context(), auth.SessionID)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	if isRevoked {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	accountDTO, err := ctrl.tenantUserFacade.Account.FindOneByID(
		c.Context(),
		auth.AccountID,
		nil,
		&tenantUserUC.AccountDTOOptions{
			FetchTenant: true,
		},
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	out, err := OutWhoIAmDTO(c, ctrl.fileUC, accountDTO)
	if err != nil {
		return err
	}

	return c.JSON(out)
}
