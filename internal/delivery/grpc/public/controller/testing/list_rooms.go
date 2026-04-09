package testing

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/delivery/grpc/mapper"
	testingUC "github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/neurochar/backend/pkg/auth"
	typesv1 "github.com/neurochar/backend/pkg/proto_pb/common/types"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/testing/v1"
	"github.com/samber/lo"
)

func (ctrl *Controller) ListRooms(
	ctx context.Context,
	req *desc.ListRoomsRequest,
) (*desc.ListRoomsResponse, error) {
	const op = "ListRooms"

	ctx = auth.WithCheckTenantAccess(ctx)

	authData := auth.GetAuthData(ctx)
	if authData == nil || !authData.IsTenantUser() {
		return nil, appErrors.Chainf(appErrors.ErrUnauthorized, "%s.%s", ctrl.pkg, op)
	}

	limit := req.GetLimit()
	if limit == 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	} else if limit < 1 {
		limit = 1
	}

	offset := req.GetOffset()

	listOptions := &testingUC.RoomListOptions{
		FilterTenantID: &authData.TenantUserClaims().TenantID,
		Sort: []uctypes.SortOption[testingUC.RoomListOptionsSortField]{
			{
				Field:  testingUC.RoomListOptionsSortFieldCreatedAt,
				IsDesc: true,
			},
		},
	}

	listParams := &uctypes.QueryGetListParams{
		Limit:  limit,
		Offset: offset,
	}

	items, total, err := ctrl.testingFacade.Room.FindPagedList(
		ctx,
		listOptions,
		listParams,
		&testingUC.RoomDTOOptions{
			FetchCandidate: true,
			FetchProfile:   true,
		},
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
	}

	return &desc.ListRoomsResponse{
		Items: lo.Map(items, func(item *testingUC.RoomDTO, _ int) *typesv1.TestingListRoom {
			return mapper.TestingRoomDTOToListPb(item)
		}),
		Total: int32(total),
	}, nil
}
