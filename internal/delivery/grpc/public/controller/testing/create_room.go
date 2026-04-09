package testing

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	testingUC "github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/testing/v1"
)

func (ctrl *Controller) CreateRoom(
	ctx context.Context,
	req *desc.CreateRoomRequest,
) (*desc.CreateRoomResponse, error) {
	const op = "CreateRoom"

	if req.Payload == nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", ctrl.pkg, op)
	}

	ctx = auth.WithCheckTenantAccess(ctx)

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	candidateId, err := uuid.Parse(req.Payload.CandidateId)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	profileId, err := uuid.Parse(req.Payload.ProfileId)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	usecaseInput := testingUC.CreateRoomDataInput{
		CandidateID: candidateId,
		ProfileID:   profileId,
		CreatedBy:   &authData.TenantUserClaims().AccountID,
	}

	roomDTO, err := ctrl.testingFacade.Room.CreateByDTO(
		ctx,
		authData.TenantUserClaims().TenantID,
		usecaseInput,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.CreateRoomResponse{
		Item: mapper.TestingRoomDTOToPb(roomDTO),
	}, nil
}
