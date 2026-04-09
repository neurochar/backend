package crm

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper/helpers"
	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/crm/v1"
	"github.com/samber/lo"
)

func (ctrl *Controller) CreateCandidate(
	ctx context.Context,
	req *desc.CreateCandidateRequest,
) (*desc.CreateCandidateResponse, error) {
	const op = "CreateCandidate"

	if req.Payload == nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", ctrl.pkg, op)
	}

	ctx = auth.WithCheckTenantAccess(ctx)

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	candidateGender, err := crmEntity.CandidateGenderFromUint8(uint8(req.Payload.Gender))
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	candidateDTO, err := ctrl.crmFacade.Candidate.CreateByDTO(
		ctx,
		authData.TenantUserClaims().TenantID,
		crmUC.CreateCandidateDataInput{
			CandidateName:     req.Payload.Name,
			CandidateSurname:  req.Payload.Surname,
			CandidateGender:   candidateGender,
			CandidateBirthday: helpers.PbDateToTimePtr(lo.FromPtr(req.Payload.Birthday).Date),
			CreatedBy:         &authData.TenantUserClaims().AccountID,
		},
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.CreateCandidateResponse{
		Item: mapper.CandidateDTOToPb(candidateDTO),
	}, nil
}
