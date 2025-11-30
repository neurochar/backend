package room

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	candidateUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/domain/testing/entity"
	"github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) entitiesToDTO(
	ctx context.Context,
	items []*entity.Room,
	dtoOpts *usecase.RoomDTOOptions,
) ([]*usecase.RoomDTO, error) {
	const op = "entitiesToDTO"

	tenantAccountsMap := make(map[uuid.UUID]*tenantUC.AccountDTO, 0)
	tenantAccountsIDs := make([]uuid.UUID, 0)

	candidatesMap := make(map[uuid.UUID]*candidateUC.CandidateDTO, 0)
	candidatesIDs := make([]uuid.UUID, 0)

	profilesMap := make(map[uuid.UUID]*usecase.ProfileDTO, 0)
	profilesIDs := make([]uuid.UUID, 0)

	for _, item := range items {
		if item.CreatedBy != nil {
			tenantAccountsIDs = append(tenantAccountsIDs, *item.CreatedBy)
		}

		if item.CandidateID != nil {
			candidatesIDs = append(candidatesIDs, *item.CandidateID)
		}

		if item.ProfileID != nil {
			profilesIDs = append(profilesIDs, *item.ProfileID)
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

	if (dtoOpts == nil || dtoOpts.FetchCandidate) && len(candidatesIDs) > 0 {
		candidatesList, err := uc.candidateUC.FindList(ctx, &candidateUC.CandidateListOptions{
			FilterIDs: &candidatesIDs,
		}, nil, nil)
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		candidatesMap = lo.SliceToMap(
			candidatesList,
			func(item *candidateUC.CandidateDTO) (uuid.UUID, *candidateUC.CandidateDTO) {
				return item.Candidate.ID, item
			},
		)
	}

	if (dtoOpts == nil || dtoOpts.FetchProfile) && len(profilesIDs) > 0 {
		profilesList, err := uc.profileUC.FindList(ctx, &usecase.ProfileListOptions{
			FilterIDs: &profilesIDs,
		}, nil, nil)
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		profilesMap = lo.SliceToMap(
			profilesList,
			func(item *usecase.ProfileDTO) (uuid.UUID, *usecase.ProfileDTO) {
				return item.Profile.ID, item
			},
		)
	}

	out := make([]*usecase.RoomDTO, 0, len(items))

	for _, item := range items {
		resItem := &usecase.RoomDTO{
			Room: item,
		}

		if (dtoOpts == nil || dtoOpts.FetchCreatedBy) && item.CreatedBy != nil {
			account, ok := tenantAccountsMap[*item.CreatedBy]
			if !ok {
				return nil, appErrors.Chainf(appErrors.ErrInternal.Extend("account not fetched"), "%s.%s", uc.pkg, op)
			}

			resItem.CreatedBy = account
		}

		if (dtoOpts == nil || dtoOpts.FetchCandidate) && item.CandidateID != nil {
			candidate, ok := candidatesMap[*item.CandidateID]
			if !ok {
				return nil, appErrors.Chainf(appErrors.ErrInternal.Extend("candidate not fetched"), "%s.%s", uc.pkg, op)
			}

			resItem.CandidateDTO = candidate
		}

		if (dtoOpts == nil || dtoOpts.FetchProfile) && item.ProfileID != nil {
			profile, ok := profilesMap[*item.ProfileID]
			if !ok {
				return nil, appErrors.Chainf(appErrors.ErrInternal.Extend("profile not fetched"), "%s.%s", uc.pkg, op)
			}

			resItem.ProfileDTO = profile
		}

		out = append(out, resItem)
	}

	return out, nil
}
