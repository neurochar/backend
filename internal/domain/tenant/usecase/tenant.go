package usecase

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/tenant/entity"
)

var ErrTenantAlreadyExists = appErrors.ErrBadRequest.WithTextCode("ALREADY_EXISTS")

type TenantListOptions struct {
	FilterIDs *[]uuid.UUID
}

type CreateTenantIn struct {
	Name     string
	TextID   string
	IsDemo   bool
	IsActive bool
}

type PatchTenantDataInput struct {
	Version int64

	Name *string
}

type TenantUsecase interface {
	FindList(
		ctx context.Context,
		listOptions *TenantListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*entity.Tenant, resErr error)

	FindPagedList(
		ctx context.Context,
		listOptions *TenantListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*entity.Tenant, total uint64, resErr error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (resItem *entity.Tenant, resErr error)

	FindOneByTextID(
		ctx context.Context,
		textID string,
		queryParams *uctypes.QueryGetOneParams,
	) (resItem *entity.Tenant, resErr error)

	CreateByDTO(ctx context.Context, in CreateTenantIn) (resItem *entity.Tenant, resErr error)

	PatchByDTO(
		ctx context.Context,
		id uuid.UUID,
		in PatchTenantDataInput,
		skipVersionCheck bool,
	) (resErr error)
}

type TenantRepository interface {
	FindList(
		ctx context.Context,
		listOptions *TenantListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*entity.Tenant, resErr error)

	FindPagedList(
		ctx context.Context,
		listOptions *TenantListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*entity.Tenant, total uint64, resErr error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (resFile *entity.Tenant, resErr error)

	FindOneByTextID(
		ctx context.Context,
		textID string,
		queryParams *uctypes.QueryGetOneParams,
	) (resFile *entity.Tenant, resErr error)

	Create(ctx context.Context, item *entity.Tenant) (resErr error)

	Update(ctx context.Context, item *entity.Tenant) (resErr error)
}
