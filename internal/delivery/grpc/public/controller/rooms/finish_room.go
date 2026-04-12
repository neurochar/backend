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

func (ctrl *Controller) FinishRoom(
	ctx context.Context,
	req *desc.FinishRoomRequest,
) (*desc.FinishRoomResponse, error) {
	const op = "FinishRoom"

	ctx = auth.WithoutCheckTenantAccess(ctx)

	ip := tools.GetRealIP(ctx)

	if req.Payload == nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", ctrl.pkg, op)
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithWrap(err), "%s.%s", ctrl.pkg, op)
	}

	answerData := make(map[uint64]any, len(req.Payload.GetAnswers().Data))

	if req.Payload.Answers != nil {
		for i, v := range req.Payload.Answers.Data {
			answerData[i] = mapper.ParseAnswerValue(v)
		}
	}

	err = ctrl.testingFacade.Room.Finish(ctx, id, answerData, ip)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.FinishRoomResponse{}, nil
}
