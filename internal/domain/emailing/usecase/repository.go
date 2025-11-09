package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/emailing/entity"
)

type ItemRepository interface {
	FindList(
		ctx context.Context,
		listOptions *ListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*entity.Item, resErr error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (resFile *entity.Item, resErr error)

	Create(ctx context.Context, item *entity.Item) (resErr error)

	Update(ctx context.Context, item *entity.Item) (resErr error)

	DeleteByID(ctx context.Context, id uuid.UUID) (resErr error)
}
