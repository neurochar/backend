package candidate_resume

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/crm/entity"
	"github.com/neurochar/backend/internal/domain/crm/usecase"
	"github.com/neurochar/backend/pkg/auth"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (uc *UsecaseImpl) FindOneByID(
	ctx context.Context,
	id uuid.UUID,
	queryParams *uctypes.QueryGetOneParams,
	dtoOpts *usecase.CandidateResumeDTOOptions,
) (*usecase.CandidateResumeDTO, error) {
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

	dto, err := uc.entitiesToDTO(ctx, []*entity.CandidateResume{item}, dtoOpts)
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
	listOptions *usecase.CandidateResumeListOptions,
	queryParams *uctypes.QueryGetListParams,
	dtoOpts *usecase.CandidateResumeDTOOptions,
) ([]*usecase.CandidateResumeDTO, error) {
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
