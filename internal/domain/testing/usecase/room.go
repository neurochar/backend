package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"
	crmUsecase "github.com/neurochar/backend/internal/domain/crm/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/domain/testing/entity"
)

type RoomListOptions struct {
	FilterTenantID *uuid.UUID
	Sort           []uctypes.SortOption[RoomListOptionsSortField]
}

type RoomListOptionsSortField string

const (
	RoomListOptionsSortFieldCreatedAt RoomListOptionsSortField = "created_at"
)

type RoomDTOOptions struct {
	FetchCreatedBy bool
	FetchCandidate bool
	FetchProfile   bool
}

type RoomDTO struct {
	Room         *entity.Room
	CandidateDTO *crmUsecase.CandidateDTO
	ProfileDTO   *ProfileDTO
	CreatedBy    *tenantUC.AccountDTO
}

type CreateRoomDataInput struct {
	CandidateID uuid.UUID
	ProfileID   uuid.UUID
	CreatedBy   *uuid.UUID
}

type PatchRoomDataInput struct {
	Version int64
}

type RoomUsecase interface {
	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
		dtoOpts *RoomDTOOptions,
	) (resRoom *RoomDTO, resErr error)

	FindList(
		ctx context.Context,
		listOptions *RoomListOptions,
		queryParams *uctypes.QueryGetListParams,
		dtoOpts *RoomDTOOptions,
	) (resItems []*RoomDTO, resErr error)

	FindPagedList(
		ctx context.Context,
		listOptions *RoomListOptions,
		queryParams *uctypes.QueryGetListParams,
		dtoOpts *RoomDTOOptions,
	) (resItems []*RoomDTO, total uint64, resErr error)

	FindListInMap(
		ctx context.Context,
		listOptions *RoomListOptions,
		queryParams *uctypes.QueryGetListParams,
		dtoOpts *RoomDTOOptions,
	) (resItems map[uuid.UUID]*RoomDTO, resErr error)

	CreateByDTO(
		ctx context.Context,
		tenantID uuid.UUID,
		in CreateRoomDataInput,
	) (resRoomDTO *RoomDTO, resErr error)

	PatchByDTO(
		ctx context.Context,
		id uuid.UUID,
		in PatchRoomDataInput,
		skipVersionCheck bool,
	) (resErr error)

	Update(ctx context.Context, item *entity.Room) (resErr error)
}

type RoomRepository interface {
	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (account *entity.Room, err error)

	FindList(
		ctx context.Context,
		listOptions *RoomListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Room, err error)

	FindPagedList(
		ctx context.Context,
		listOptions *RoomListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Room, total uint64, err error)

	Create(ctx context.Context, item *entity.Room) (err error)

	Update(ctx context.Context, item *entity.Room) (err error)
}
