package testing

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	testingEntity "github.com/neurochar/backend/internal/domain/testing/entity"
	"github.com/neurochar/backend/pkg/auth"
	typesv1 "github.com/neurochar/backend/pkg/proto_pb/common/types"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/testing/v1"
	"github.com/samber/lo"
)

func (ctrl *Controller) GetPersonalityTraits(
	ctx context.Context,
	req *desc.GetPersonalityTraitsRequest,
) (*desc.GetPersonalityTraitsResponse, error) {
	const op = "GetPersonalityTraits"

	ctx = auth.WithCheckTenantAccess(ctx)

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	items, err := ctrl.testingFacade.PersonalityTrait.FindList(
		ctx,
		nil,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.GetPersonalityTraitsResponse{
		Items: lo.Map(items, func(i testingEntity.PersonalityTrait, _ int) *typesv1.PersonalityTrait {
			return mapper.TestingPersonalityTraitToPb(i)
		}),
	}, nil
}
