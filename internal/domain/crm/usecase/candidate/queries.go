package candidate

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/crm/entity"
	"github.com/neurochar/backend/internal/domain/crm/usecase"
	"github.com/neurochar/backend/pkg/auth"
	"github.com/samber/lo"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (uc *UsecaseImpl) FindOneByID(
	ctx context.Context,
	id uuid.UUID,
	queryParams *uctypes.QueryGetOneParams,
	dtoOpts *usecase.CandidateDTOOptions,
) (*usecase.CandidateDTO, error) {
	const op = "FindOneByID"

	item, err := uc.repo.FindOneByID(ctx, id, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if auth.IsNeedToCheckTenantAccess(ctx) {
		authData := auth.GetAuthData(ctx)
		if authData == nil || !authData.IsTenantUser() || authData.TenantUserClaims().TenantID != item.TenantID {
			return nil, appErrors.Chainf(appErrors.ErrForbidden, "%s.%s", uc.pkg, op)
		}
	}

	dto, err := uc.entitiesToDTO(ctx, []*entity.Candidate{item}, dtoOpts)
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
	listOptions *usecase.CandidateListOptions,
	queryParams *uctypes.QueryGetListParams,
	dtoOpts *usecase.CandidateDTOOptions,
) ([]*usecase.CandidateDTO, error) {
	const op = "FindList"

	items, err := uc.repo.FindList(ctx, listOptions, queryParams)
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
	listOptions *usecase.CandidateListOptions,
	queryParams *uctypes.QueryGetListParams,
	dtoOpts *usecase.CandidateDTOOptions,
) ([]*usecase.CandidateDTO, uint64, error) {
	const op = "FindPagedList"

	items, total, err := uc.repo.FindPagedList(ctx, listOptions, queryParams)
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
	listOptions *usecase.CandidateListOptions,
	queryParams *uctypes.QueryGetListParams,
	dtoOpts *usecase.CandidateDTOOptions,
) (map[uuid.UUID]*usecase.CandidateDTO, error) {
	const op = "FindListInMap"

	items, err := uc.repo.FindList(ctx, listOptions, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	out, err := uc.entitiesToDTO(ctx, items, dtoOpts)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	result := lo.SliceToMap(out, func(item *usecase.CandidateDTO) (uuid.UUID, *usecase.CandidateDTO) {
		return item.Candidate.ID, item
	})

	return result, nil
}
