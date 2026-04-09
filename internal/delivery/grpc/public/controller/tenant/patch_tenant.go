package tenant

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/tenant/v1"
)

func (ctrl *Controller) PatchTenant(
	ctx context.Context,
	req *desc.PatchTenantRequest,
) (*desc.PatchTenantResponse, error) {
	const op = "PatchTenant"

	if req.Payload == nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", ctrl.pkg, op)
	}

	ctx = auth.WithCheckTenantAccess(ctx)

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	usecaseInput := tenantUC.PatchTenantDataInput{
		Version: req.Version,

		Name: req.Payload.Name,
	}

	err = ctrl.tenantFacade.Cross.PatchTenantByDTO(
		ctx,
		id,
		usecaseInput,
		req.SkipVersionCheck,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.PatchTenantResponse{}, nil
}
