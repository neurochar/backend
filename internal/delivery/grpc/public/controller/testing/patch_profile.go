package testing

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	testingEntity "github.com/neurochar/backend/internal/domain/testing/entity"
	testingUC "github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/neurochar/backend/pkg/auth"
	typesv1 "github.com/neurochar/backend/pkg/proto_pb/common/types"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/testing/v1"
	"github.com/samber/lo"
)

func (ctrl *Controller) PatchProfile(
	ctx context.Context,
	req *desc.PatchProfileRequest,
) (*desc.PatchProfileResponse, error) {
	const op = "PatchProfile"

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

	usecaseInput := testingUC.PatchProfileDataInput{
		Version: req.Version,

		Name:        req.Payload.Name,
		Description: req.Payload.Description,
	}

	if req.Payload.PersonalityTraits != nil {
		personalityTraitsMap := lo.MapValues(
			req.Payload.PersonalityTraits.Map,
			func(v *typesv1.ProfilePersonalityTraitsMapItem, _ uint64) testingEntity.ProfilePersonalityTraitsMapItem {
				return mapper.TestingPersonalityTraitsMapItemPbToEntity(v)
			},
		)

		usecaseInput.PersonalityTraitsMap = lo.ToPtr(testingEntity.ProfilePersonalityTraitsMap(personalityTraitsMap))
	}

	err = ctrl.testingFacade.Profile.PatchByDTO(
		ctx,
		id,
		usecaseInput,
		req.SkipVersionCheck,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.PatchProfileResponse{}, nil
}
