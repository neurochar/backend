package crm_tenant

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/crm/v1"
)

func (ctrl *controller) DeleteCandidate(
	ctx context.Context,
	req *desc.DeleteCandidateRequest,
) (*desc.DeleteCandidateResponse, error) {
	const op = "DeleteCandidate"

	ctx = auth.WithCheckTenantAccess(ctx)

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	err = ctrl.crmFacade.Cross.Delete(
		ctx,
		id,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.DeleteCandidateResponse{}, nil
}
