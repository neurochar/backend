package rooms

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/rooms/v1"
)

func (ctrl *Controller) GetRoom(
	ctx context.Context,
	req *desc.GetRoomRequest,
) (*desc.GetRoomResponse, error) {
	const op = "GetRoom"

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	roomDTO, err := ctrl.testingFacade.Room.FindOneByID(auth.WithoutCheckTenantAccess(ctx), id, nil, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	tenant, err := ctrl.tenantFacade.Tenant.FindOneByID(auth.WithoutCheckTenantAccess(ctx), roomDTO.Room.TenantID, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	outRoom, err := mapper.RoomToPb(roomDTO, tenant)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.GetRoomResponse{
		Room: outRoom,
	}, nil
}
