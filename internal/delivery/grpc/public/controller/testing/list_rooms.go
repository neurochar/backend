package testing

import (
	"context"

	"github.com/google/uuid"
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
	}

	if req.Sort != nil {
		sortField := testingUC.RoomListOptionsSortFieldCreatedAt

		switch *req.Sort {
		case desc.ListRoomsSort_LIST_ROOM_SORT_CREATED_AT:
			sortField = testingUC.RoomListOptionsSortFieldCreatedAt
		case desc.ListRoomsSort_LIST_ROOM_SORT_FINISHED_AT:
			sortField = testingUC.RoomListOptionsSortFieldFinishedAt
		case desc.ListRoomsSort_LIST_ROOM_SORT_RESULT_INDEX:
			sortField = testingUC.RoomListOptionsSortFieldResultIndex
		}

		listOptions.Sort = []uctypes.SortOption[testingUC.RoomListOptionsSortField]{
			{
				Field:  sortField,
				IsDesc: true,
			},
		}
	}

	if req.FilterCandidateId != nil {
		id, err := uuid.Parse(*req.FilterCandidateId)
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
		}

		listOptions.FilterCandidateID = lo.ToPtr(id)
	}

	if req.FilterProfileId != nil {
		id, err := uuid.Parse(*req.FilterProfileId)
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", ctrl.pkg, op)
		}

		listOptions.FilterProfileID = lo.ToPtr(id)
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
