package candidate

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/domain/crm/entity"
	"github.com/neurochar/backend/internal/domain/crm/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) entitiesToDTO(
	ctx context.Context,
	items []*entity.Candidate,
	dtoOpts *usecase.CandidateDTOOptions,
) ([]*usecase.CandidateDTO, error) {
	const op = "entitiesToDTO"

	tenantAccountsMap := make(map[uuid.UUID]*tenantUC.AccountDTO, 0)

	tenantAccountsIDs := make([]uuid.UUID, 0)

	for _, item := range items {
		if item.CreatedBy != nil {
			tenantAccountsIDs = append(tenantAccountsIDs, *item.CreatedBy)
		}
	}

	if (dtoOpts == nil || dtoOpts.FetchCreatedBy) && len(tenantAccountsIDs) > 0 {
		accountsList, err := uc.tenantAccountUC.FindList(ctx, &tenantUC.AccountListOptions{
			FilterIDs: &tenantAccountsIDs,
		}, nil, nil)
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		tenantAccountsMap = lo.SliceToMap(accountsList, func(item *tenantUC.AccountDTO) (uuid.UUID, *tenantUC.AccountDTO) {
			return item.Account.ID, item
		})
	}

	out := make([]*usecase.CandidateDTO, 0, len(items))

	for _, item := range items {
		resItem := &usecase.CandidateDTO{
			Candidate: item,
		}

		if (dtoOpts == nil || dtoOpts.FetchCreatedBy) && item.CreatedBy != nil {
			account, ok := tenantAccountsMap[*item.CreatedBy]
			if !ok {
				return nil, appErrors.Chainf(appErrors.ErrInternal.Extend("account not fetched"), "%s.%s", uc.pkg, op)
			}

			resItem.CreatedBy = account
		}

		out = append(out, resItem)
	}

	return out, nil
}
