package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/domain/testing/entity"
)

type ProfileListOptions struct {
	FilterIDs      *[]uuid.UUID
	FilterTenantID *uuid.UUID
	SearchQuery    *string
	Sort           []uctypes.SortOption[ProfileListOptionsSortField]
}

type ProfileListOptionsSortField string

const (
	ProfileListOptionsSortFieldCreatedAt ProfileListOptionsSortField = "created_at"
)

type ProfileDTOOptions struct {
	FetchCreatedBy bool
}

type ProfileDTO struct {
	Profile   *entity.Profile
	CreatedBy *tenantUC.AccountDTO
}

type CreateProfileDataInput struct {
	Name                 string
	Description          string
	PersonalityTraitsMap entity.ProfilePersonalityTraitsMap
	CreatedBy            *uuid.UUID
}

type PatchProfileDataInput struct {
	Version int64

	Name                 *string
	Description          *string
	PersonalityTraitsMap *entity.ProfilePersonalityTraitsMap
}

type ProfileUsecase interface {
	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
		dtoOpts *ProfileDTOOptions,
	) (resProfile *ProfileDTO, resErr error)

	FindList(
		ctx context.Context,
		listOptions *ProfileListOptions,
		queryParams *uctypes.QueryGetListParams,
		dtoOpts *ProfileDTOOptions,
	) (resItems []*ProfileDTO, resErr error)

	FindPagedList(
		ctx context.Context,
		listOptions *ProfileListOptions,
		queryParams *uctypes.QueryGetListParams,
		dtoOpts *ProfileDTOOptions,
	) (resItems []*ProfileDTO, total uint64, resErr error)

	FindListInMap(
		ctx context.Context,
		listOptions *ProfileListOptions,
		queryParams *uctypes.QueryGetListParams,
		dtoOpts *ProfileDTOOptions,
	) (resItems map[uuid.UUID]*ProfileDTO, resErr error)

	CreateByDTO(
		ctx context.Context,
		tenantID uuid.UUID,
		in CreateProfileDataInput,
	) (resProfileDTO *ProfileDTO, resErr error)

	PatchByDTO(
		ctx context.Context,
		id uuid.UUID,
		in PatchProfileDataInput,
		skipVersionCheck bool,
	) (resErr error)

	Update(ctx context.Context, item *entity.Profile) (resErr error)

	GenerateProfileDescriptionByName(ctx context.Context, name string) (string, error)

	GenerateProfileTraitsMapByDescription(
		ctx context.Context,
		req *GenerateProfileTraitsMapByDescriptionRequest,
	) (*GenerateProfileTraitsMapByDescriptionResponse, error)
}

type ProfileRepository interface {
	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (account *entity.Profile, err error)

	FindList(
		ctx context.Context,
		listOptions *ProfileListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Profile, err error)

	FindPagedList(
		ctx context.Context,
		listOptions *ProfileListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*entity.Profile, total uint64, err error)

	Create(ctx context.Context, item *entity.Profile) (err error)

	Update(ctx context.Context, item *entity.Profile) (err error)
}
