package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/tenant/entity"
	"github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/auth"
	"github.com/samber/lo"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (uc *UsecaseImpl) FindOneByEmail(
	ctx context.Context,
	tenantID uuid.UUID,
	email string,
	queryParams *uctypes.QueryGetOneParams,
	dtoOpts *usecase.AccountDTOOptions,
) (*usecase.AccountDTO, error) {
	const op = "FindOneByEmail"

	item, err := uc.repoAccount.FindOneByEmail(ctx, tenantID, email, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if auth.IsNeedToCheckRights(ctx) {
		authData := auth.GetAuthData(ctx)
		if authData == nil || authData.TenantID != item.TenantID {
			return nil, appErrors.Chainf(appErrors.ErrForbidden, "%s.%s", uc.pkg, op)
		}
	}

	dto, err := uc.entitiesToDTO(ctx, []*entity.Account{item}, dtoOpts)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(dto) == 0 {
		return nil, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	return dto[0], nil
}

func (uc *UsecaseImpl) FindOneByID(
	ctx context.Context,
	id uuid.UUID,
	queryParams *uctypes.QueryGetOneParams,
	dtoOpts *usecase.AccountDTOOptions,
) (*usecase.AccountDTO, error) {
	const op = "FindOneByID"

	item, err := uc.repoAccount.FindOneByID(ctx, id, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if auth.IsNeedToCheckRights(ctx) {
		authData := auth.GetAuthData(ctx)
		if authData == nil || authData.TenantID != item.TenantID {
			return nil, appErrors.Chainf(appErrors.ErrForbidden, "%s.%s", uc.pkg, op)
		}
	}

	dto, err := uc.entitiesToDTO(ctx, []*entity.Account{item}, dtoOpts)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(dto) == 0 {
		return nil, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	return dto[0], nil
}

func (uc *UsecaseImpl) FindList(
	ctx context.Context,
	listOptions *usecase.AccountListOptions,
	queryParams *uctypes.QueryGetListParams,
	dtoOpts *usecase.AccountDTOOptions,
) ([]*usecase.AccountDTO, error) {
	const op = "FindList"

	items, err := uc.repoAccount.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out, err := uc.entitiesToDTO(ctx, items, dtoOpts)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return out, nil
}

func (uc *UsecaseImpl) FindPagedList(
	ctx context.Context,
	listOptions *usecase.AccountListOptions,
	queryParams *uctypes.QueryGetListParams,
	dtoOpts *usecase.AccountDTOOptions,
) ([]*usecase.AccountDTO, uint64, error) {
	const op = "FindPagedList"

	items, total, err := uc.repoAccount.FindPagedList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, 0, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out, err := uc.entitiesToDTO(ctx, items, dtoOpts)
	if err != nil {
		return nil, 0, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return out, total, nil
}

func (uc *UsecaseImpl) FindListInMap(
	ctx context.Context,
	listOptions *usecase.AccountListOptions,
	queryParams *uctypes.QueryGetListParams,
	dtoOpts *usecase.AccountDTOOptions,
) (map[uuid.UUID]*usecase.AccountDTO, error) {
	const op = "FindListInMap"

	items, err := uc.repoAccount.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out, err := uc.entitiesToDTO(ctx, items, dtoOpts)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	result := lo.SliceToMap(out, func(item *usecase.AccountDTO) (uuid.UUID, *usecase.AccountDTO) {
		return item.Account.ID, item
	})

	return result, nil
}
