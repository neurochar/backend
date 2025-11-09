package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/tenant/entity"
)

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
