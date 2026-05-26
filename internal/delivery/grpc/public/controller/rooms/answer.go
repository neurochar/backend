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

func (ctrl *Controller) Answer(
	ctx context.Context,
	req *desc.AnswerRequest,
) (*desc.AnswerResponse, error) {
	const op = "Answer"

	ctx = auth.WithoutCheckTenantAccess(ctx)

	ip := tools.GetRealIP(ctx)

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	room, err := ctrl.testingFacade.Room.Answer(ctx, id, req.QuestionIndex, mapper.ParseAnswerValue(req.AnswerValue), ip)
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

	return &desc.AnswerResponse{
		Room: outRoom,
	}, nil
}
