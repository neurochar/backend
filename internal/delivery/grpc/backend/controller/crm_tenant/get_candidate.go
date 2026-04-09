package crm_tenant

import (
	"context"
	"errors"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/grpc/backend/mapper"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/crm/v1"
)

func (ctrl *controller) GetCandidate(ctx context.Context, req *desc.GetCandidateRequest) (*desc.GetCandidateResponse, error) {
	const op = "GetCandidate"

	ctx = auth.WithCheckTenantAccess(ctx)

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	candidateDTO, err := ctrl.crmFacade.Candidate.FindOneByID(ctx, id, nil, nil)
	if err != nil {
		if errors.Is(err, appErrors.ErrForbidden) {
			return nil, appErrors.Chainf(appErrors.ErrNotFound, "%s.%s", ctrl.pkg, op)
		}
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.GetCandidateResponse{
		Item: mapper.CandidateDTOToPb(candidateDTO),
	}, nil
}
