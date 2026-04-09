package crm

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper/helpers"
	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/crm/v1"
	"github.com/samber/lo"
)

func (ctrl *Controller) PatchCandidate(
	ctx context.Context,
	req *desc.PatchCandidateRequest,
) (*desc.PatchCandidateResponse, error) {
	const op = "PatchCandidate"

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

	usecaseInput := crmUC.PatchCandidateDataInput{
		Version: req.Version,

		CandidateName:    req.Payload.Name,
		CandidateSurname: req.Payload.Surname,
	}

	if req.Payload.Gender != nil {
		candidateGender, err := crmEntity.CandidateGenderFromUint8(uint8(req.Payload.GetGender()))
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
		}

		usecaseInput.CandidateGender = &candidateGender
	}

	if req.Payload.Birthday != nil {
		usecaseInput.CandidateBirthday = lo.ToPtr(helpers.PbDateToTimePtr(req.Payload.Birthday.Date))
	}

	err = ctrl.crmFacade.Candidate.PatchByDTO(
		ctx,
		id,
		usecaseInput,
		req.SkipVersionCheck,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.PatchCandidateResponse{}, nil
}
