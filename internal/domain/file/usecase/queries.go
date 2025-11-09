package usecase

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
)

func (uc *UsecaseImpl) FindOneByID(
	ctx context.Context,
	id uuid.UUID,
	queryParams *uctypes.QueryGetOneParams,
) (*fileEntity.File, error) {
	const op = "FindOneByID"

	item, err := uc.repo.FindOneByID(ctx, id, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return item, nil
}

func (uc *UsecaseImpl) FindList(
	ctx context.Context,
	listOptions *ListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*fileEntity.File, error) {
	const op = "FindList"

	items, err := uc.repo.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return items, nil
}

func (uc *UsecaseImpl) FindListInMap(
	ctx context.Context,
	listOptions *ListOptions,
	queryParams *uctypes.QueryGetListParams,
) (map[uuid.UUID]*fileEntity.File, error) {
	const op = "FindListInMap"

	items, err := uc.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out := make(map[uuid.UUID]*fileEntity.File, len(items))
	for _, item := range items {
		out[item.ID] = item
	}

	return out, nil
}
