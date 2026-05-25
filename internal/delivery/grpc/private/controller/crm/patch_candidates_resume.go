package crm

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	crmEntity "github.com/neurochar/backend/internal/domain/crm/entity"
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	desc "github.com/neurochar/backend/pkg/proto_pb/private/crm/v1"
	"github.com/samber/lo"
)

func (ctrl *Controller) PatchCandidatesResume(
	ctx context.Context,
	req *desc.PatchCandidatesResumeRequest,
) (*desc.PatchCandidatesResumeResponse, error) {
	const op = "PatchCandidatesResume"

	if req.Payload == nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", ctrl.pkg, op)
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	usecaseInput := crmUC.PatchCandidateResumeDataInput{
		Version: req.Version,
	}

	if req.Payload.Status != nil {
		usecaseInput.Status = lo.ToPtr(mapper.CandidateResumeStatusPbToEntity(*req.Payload.Status))
	}

	if req.Payload.AnalyzeData != nil {
		var analyzeData *crmEntity.CandidateResumeAnalyzeData
		if req.Payload.AnalyzeData.Data != nil {
			analyzeData = &crmEntity.CandidateResumeAnalyzeData{
				AnonymizedText: req.Payload.AnalyzeData.Data.AnonymizedText,
				DataVersion:    req.Payload.AnalyzeData.Data.DataVersion,
			}
		}

		usecaseInput.AnalyzeData = &analyzeData
	}

	if req.Payload.ErrorText != nil {
		var errorText *string
		if req.Payload.ErrorText.Text != nil {
			errorText = req.Payload.ErrorText.Text
		}

		usecaseInput.ErrorText = &errorText
	}

	err = ctrl.crmFacade.CandidateResume.PatchByDTO(
		ctx,
		id,
		usecaseInput,
		req.SkipVersionCheck,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.PatchCandidatesResumeResponse{}, nil
}
