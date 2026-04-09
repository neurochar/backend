package tenants

import (
	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	backendMiddleware "github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/pkg/validation"

	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
)

type PatchTenantHandlerIn struct {
	Version          int64 `json:"_version"`
	SkipVersionCheck bool  `json:"_skipVersionCheck"`

	Name *string `json:"name" validate:"omitempty,min=1,max=150"`
}

func (ctrl *Controller) PatchTenantHandler(c *fiber.Ctx) error {
	const op = "PatchTenantHandler"

	in := &PatchTenantHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(
			httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...),
			"%s.%s", ctrl.pkg, op,
		)
	}

	auth := backendMiddleware.GetAuthData(c)
	if auth == nil || !auth.IsTenantUser() {
		return appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	usecaseInput := tenantUC.PatchTenantDataInput{
		Version: in.Version,

		Name: in.Name,
	}

	err := ctrl.tenantFacade.Cross.PatchTenantByDTO(
		c.Context(),
		auth.TenantUserClaims().TenantID,
		usecaseInput,
		in.SkipVersionCheck,
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return nil
}
