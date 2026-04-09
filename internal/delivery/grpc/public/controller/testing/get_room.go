package testing

import (
	"context"
	"errors"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/testing/v1"
)

func (ctrl *Controller) GetRoom(ctx context.Context, req *desc.GetRoomRequest) (*desc.GetRoomResponse, error) {
	const op = "GetRoom"

	ctx = auth.WithCheckTenantAccess(ctx)

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	roomDTO, err := ctrl.testingFacade.Room.FindOneByID(ctx, id, nil, nil)
	if err != nil {
		if errors.Is(err, appErrors.ErrForbidden) {
			return nil, appErrors.Chainf(appErrors.ErrNotFound, "%s.%s", ctrl.pkg, op)
		}
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.GetRoomResponse{
		Item: mapper.TestingRoomDTOToPb(roomDTO),
	}, nil
}
