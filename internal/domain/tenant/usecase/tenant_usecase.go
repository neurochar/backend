package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/tenant/entity"
)

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
}
