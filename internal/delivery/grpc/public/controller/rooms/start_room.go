package rooms

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/common/tools"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	"github.com/neurochar/backend/pkg/auth"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/rooms/v1"
)

func (ctrl *Controller) StartRoom(
	ctx context.Context,
	req *desc.StartRoomRequest,
) (*desc.StartRoomResponse, error) {
	const op = "StartRoom"

	ctx = auth.WithoutCheckTenantAccess(ctx)

	ip := tools.GetRealIP(ctx)

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	room, err := ctrl.testingFacade.Room.Start(ctx, id, ip)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	tenant, err := ctrl.tenantFacade.Tenant.FindOneByID(auth.WithoutCheckTenantAccess(ctx), room.Room.TenantID, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	outRoom, err := mapper.RoomToPb(room, tenant)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.StartRoomResponse{
		Room: outRoom,
	}, nil
}
