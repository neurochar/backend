package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"
	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
)

type Usecase interface {
	FileUrl(file *fileEntity.File) string

	UploadAndCreateFiles(
		ctx context.Context,
		input UploadFilesIn,
	) (resMap UploadFilesOut, cancelFn func(context.Context) error, resErr error)

	Create(ctx context.Context, item *fileEntity.File) (resErr error)

	Update(ctx context.Context, item *fileEntity.File) (resErr error)

	Delete(ctx context.Context, item *fileEntity.File) (resErr error)

	DeleteByIDs(ctx context.Context, ids []uuid.UUID) (resErr error)

	FindList(
		ctx context.Context,
		listOptions *ListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*fileEntity.File, err error)

	FindListInMap(
		ctx context.Context,
		listOptions *ListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items map[uuid.UUID]*fileEntity.File, err error)

	FindOneByID(ctx context.Context, id uuid.UUID, queryParams *uctypes.QueryGetOneParams) (file *fileEntity.File, err error)

	ProcessFilesToTarget(ctx context.Context, items []ProcessFileToTargetIn) (filesMap map[uuid.UUID]*fileEntity.File, err error)

	ProcessFilesSliceToTarget(
		ctx context.Context,
		oldIDs []uuid.UUID,
		newIDs []uuid.UUID,
		target string,
	) (filesMap map[uuid.UUID]*fileEntity.File, err error)

	JobProcessFilesToDelete(ctx context.Context) (anyJobDone bool, err error)

	JobProcessUnusedFiles(ctx context.Context, unusedTTL time.Duration) (anyJobDone bool, err error)
}
