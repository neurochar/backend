package registration

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/common/tools"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/registration/v1"
)

func (ctrl *Controller) FinishRegistration(
	ctx context.Context,
	req *desc.FinishRegistrationRequest,
) (*desc.FinishRegistrationResponse, error) {
	const op = "FinishRegistration"

	if req.Payload == nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", ctrl.pkg, op)
	}

	ip := tools.GetRealIP(ctx)

	registrationID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	tenant, err := ctrl.tenantFacade.Registration.FinishByDTO(
		ctx,
		registrationID,
		tenantUC.FinishRegistrationIn{
			TenantTextID:   req.Payload.TenantTextId,
			ProfileName:    req.Payload.ProfileName,
			ProfileSurname: req.Payload.ProfileSurname,
		},
		ip,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.FinishRegistrationResponse{
		Id:     tenant.ID.String(),
		TextId: tenant.TextID,
		Url:    tenant.GetUrl(ctrl.cfg.Global.TenantMainDomain, ctrl.cfg.Global.TenantUrlScheme),
	}, nil
}
