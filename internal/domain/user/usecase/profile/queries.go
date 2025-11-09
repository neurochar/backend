package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"

	appErrors "github.com/neurochar/backend/internal/app/errors"

	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	"github.com/neurochar/backend/internal/domain/user/usecase"
)

func (uc *UsecaseImpl) FindList(
	ctx context.Context,
	listOptions *usecase.ProfileListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*userEntity.Profile, error) {
	const op = "FindList"

	item, err := uc.repo.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return item, nil
}

func (uc *UsecaseImpl) FindFullList(
	ctx context.Context,
	listOptions *usecase.ProfileListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*usecase.FullProfileDTO, error) {
	const op = "FindFullList"

	items, err := uc.repo.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	filesIDs := make([]uuid.UUID, 0, len(items))
	for _, item := range items {
		if item.Photo100x100FileID != nil {
			filesIDs = append(filesIDs, *item.Photo100x100FileID)
		}
	}

	filesMap, err := uc.fileUC.FindListInMap(ctx, &fileUC.ListOptions{
		IDs: &filesIDs,
	}, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out := make([]*usecase.FullProfileDTO, 0, len(items))

	for _, item := range items {
		resItem := &usecase.FullProfileDTO{
			Profile: item,
		}

		if item.Photo100x100FileID != nil {
			file, ok := filesMap[*item.Photo100x100FileID]
			if ok {
				resItem.Photo100x100File = file
			}
		}

		out = append(out, resItem)
	}

	return out, nil
}

func (uc *UsecaseImpl) FindFullOneByID(
	ctx context.Context,
	id uint64,
	queryParams *uctypes.QueryGetOneParams,
) (*usecase.FullProfileDTO, error) {
	const op = "FindFullOneByID"

	item, err := uc.repo.FindOneByID(ctx, id, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	filesIDs := make([]uuid.UUID, 0)
	if item.Photo100x100FileID != nil {
		filesIDs = append(filesIDs, *item.Photo100x100FileID)
	}

	filesMap, err := uc.fileUC.FindListInMap(ctx, &fileUC.ListOptions{
		IDs: &filesIDs,
	}, nil)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	resItem := &usecase.FullProfileDTO{
		Profile: item,
	}

	if item.Photo100x100FileID != nil {
		file, ok := filesMap[*item.Photo100x100FileID]
		if ok {
			resItem.Photo100x100File = file
		}
	}

	return resItem, nil
}
