package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"

	appErrors "github.com/neurochar/backend/internal/app/errors"

	"github.com/neurochar/backend/internal/domain/user/usecase"
)

func (uc *UsecaseImpl) FindOneByEmail(
	ctx context.Context,
	email string,
	queryParams *uctypes.QueryGetOneParams,
) (*userEntity.Account, error) {
	const op = "FindOneByEmail"

	item, err := uc.repoAccount.FindOneByEmail(ctx, email, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return item, nil
}

func (uc *UsecaseImpl) FindOneByID(
	ctx context.Context,
	id uuid.UUID,
	queryParams *uctypes.QueryGetOneParams,
) (*userEntity.Account, error) {
	const op = "FindOneByID"

	item, err := uc.repoAccount.FindOneByID(ctx, id, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return item, nil
}

func (uc *UsecaseImpl) FindList(
	ctx context.Context,
	listOptions *usecase.AccountListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*userEntity.Account, error) {
	const op = "FindList"

	items, err := uc.repoAccount.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return items, nil
}

func (uc *UsecaseImpl) FindListInMap(
	ctx context.Context,
	listOptions *usecase.AccountListOptions,
	queryParams *uctypes.QueryGetListParams,
) (map[uuid.UUID]*userEntity.Account, error) {
	const op = "FindListInMap"

	items, err := uc.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out := make(map[uuid.UUID]*userEntity.Account, len(items))
	for _, item := range items {
		out[item.ID] = item
	}

	return out, nil
}
