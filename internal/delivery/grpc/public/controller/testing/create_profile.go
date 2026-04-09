package testing

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	testingEntity "github.com/neurochar/backend/internal/domain/testing/entity"
	testingUC "github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/neurochar/backend/pkg/auth"
	typesv1 "github.com/neurochar/backend/pkg/proto_pb/common/types"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/testing/v1"
	"github.com/samber/lo"
)

func (ctrl *Controller) CreateProfile(
	ctx context.Context,
	req *desc.CreateProfileRequest,
) (*desc.CreateProfileResponse, error) {
	const op = "CreateProfile"

	if req.Payload == nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", ctrl.pkg, op)
	}

	ctx = auth.WithCheckTenantAccess(ctx)

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	personalityTraitsMap := lo.MapValues(
		req.Payload.PersonalityTraits.Map,
		func(v *typesv1.ProfilePersonalityTraitsMapItem, _ uint64) testingEntity.ProfilePersonalityTraitsMapItem {
			return mapper.TestingPersonalityTraitsMapItemPbToEntity(v)
		},
	)

	profileDTO, err := ctrl.testingFacade.Profile.CreateByDTO(
		ctx,
		authData.TenantUserClaims().TenantID,
		testingUC.CreateProfileDataInput{
			Name:                 req.Payload.Name,
			PersonalityTraitsMap: personalityTraitsMap,
			CreatedBy:            &authData.TenantUserClaims().AccountID,
		},
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.CreateProfileResponse{
		Item: mapper.TestingProfileDTOToPb(profileDTO),
	}, nil
}
