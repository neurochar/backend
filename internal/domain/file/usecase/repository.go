package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"

	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
)

type FileRepository interface {
	FindList(
		ctx context.Context,
		listOptions *ListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*fileEntity.File, resErr error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (resFile *fileEntity.File, resErr error)

	Create(ctx context.Context, item *fileEntity.File) (resErr error)

	Update(ctx context.Context, item *fileEntity.File) (resErr error)
}
