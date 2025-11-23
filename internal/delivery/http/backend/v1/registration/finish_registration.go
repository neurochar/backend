package registration

import (
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/httperrs"
	"github.com/neurochar/backend/internal/delivery/http/middleware"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/validation"
)

type FinishRegistrationHandlerIn struct {
	RegistrationID string `json:"registrationID" validate:"omitempty,uuid"`
	TenantTextID   string `json:"tenantTextID" validate:"required"`
	ProfileName    string `json:"profileName" validate:"required,min=1,max=150"`
	ProfileSurname string `json:"profileSurname" validate:"required,min=1,max=150"`
}

func (ctrl *Controller) FinishRegistrationHandler(c *fiber.Ctx) error {
	const op = "FinishRegistrationHandler"

	in := &FinishRegistrationHandlerIn{}

	if err := c.BodyParser(in); err != nil {
		return appErrors.Chainf(httperrs.ErrCantParseBody, "%s.%s", ctrl.pkg, op)
	}

	if err := ctrl.vldtr.Struct(in); err != nil {
		return appErrors.Chainf(
			httperrs.ErrValidation.WithHints(validation.FormatErrors(err)...),
			"%s.%s", ctrl.pkg, op,
		)
	}

	registrationID, err := uuid.Parse(in.RegistrationID)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	ip := middleware.GetRealIP(c)
	requestIP := net.ParseIP(ip)

	tenant, err := ctrl.tenantFacade.Registration.FinishByDTO(
		c.Context(),
		registrationID,
		tenantUC.FinishRegistrationIn{
			TenantTextID:   in.TenantTextID,
			ProfileName:    in.ProfileName,
			ProfileSurname: in.ProfileSurname,
		},
		&requestIP,
	)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return c.JSON(OutTenant{
		ID:     tenant.ID,
		TextID: tenant.TextID,
		URL:    tenant.GetUrl(ctrl.cfg.Global.TenantMainDomain, ctrl.cfg.Global.TenantUrlScheme),
	})
}
