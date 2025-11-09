package user

import (
	"context"

	"github.com/google/uuid"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/tenant/entity"
	"github.com/neurochar/backend/internal/domain/tenant/usecase"
)

func (uc *UsecaseImpl) FindOneByID(
	ctx context.Context,
	id uuid.UUID,
	queryParams *uctypes.QueryGetOneParams,
) (*entity.Tenant, error) {
	const op = "FindOneByID"

	item, err := uc.repo.FindOneByID(ctx, id, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return item, nil
}

func (uc *UsecaseImpl) FindOneByTextID(
	ctx context.Context,
	textID string,
	queryParams *uctypes.QueryGetOneParams,
) (*entity.Tenant, error) {
	const op = "FindOneByID"

	item, err := uc.repo.FindOneByTextID(ctx, textID, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return item, nil
}

func (uc *UsecaseImpl) FindList(
	ctx context.Context,
	listOptions *usecase.TenantListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*entity.Tenant, error) {
	const op = "FindList"

	items, err := uc.repo.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return items, nil
}

func (uc *UsecaseImpl) FindPagedList(
	ctx context.Context,
	listOptions *usecase.TenantListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*entity.Tenant, uint64, error) {
	const op = "FindPagedList"

	items, total, err := uc.repo.FindPagedList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, 0, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return items, total, nil
}
