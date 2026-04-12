package candidate

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/crm/entity"
	"github.com/neurochar/backend/internal/domain/crm/usecase"
	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
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

	resumesMap := make(map[uuid.UUID]*usecase.CandidateResumeDTO, 0)
	candidatesIDs := make([]uuid.UUID, 0)

	for _, item := range items {
		candidatesIDs = append(candidatesIDs, item.ID)

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

	if (dtoOpts == nil || dtoOpts.FetchResume) && len(candidatesIDs) > 0 {
		resumeList, err := uc.candidateResumeUC.FindList(ctx, &usecase.CandidateResumeListOptions{
			FilterCandidatesIDs: &candidatesIDs,
		}, nil, nil)
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		resumesMap = lo.SliceToMap(resumeList, func(item *usecase.CandidateResumeDTO) (uuid.UUID, *usecase.CandidateResumeDTO) {
			return *item.Resume.CandidateID, item
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

		if dtoOpts == nil || dtoOpts.FetchResume {
			resume, ok := resumesMap[item.ID]
			if ok {
				resItem.Resume = resume
			}
		}

		out = append(out, resItem)
	}

	return out, nil
}

func (uc *UsecaseImpl) processResumeFilesForCandidateInTx(
	ctx context.Context,
	candidate *entity.Candidate,
	newFileID *uuid.UUID,
) error {
	const op = "processResumeFilesForCandidateInTx"

	var newFile *fileEntity.File
	if newFileID != nil {
		var err error
		newFile, err = uc.fileUC.FindOneByID(ctx, *newFileID, nil)
		if err != nil {
			return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}
	}

	currentResumes, err := uc.candidateResumeUC.FindList(ctx, &usecase.CandidateResumeListOptions{
		FilterCandidatesIDs: lo.ToPtr([]uuid.UUID{candidate.ID}),
	}, &uctypes.QueryGetListParams{ForUpdate: true}, nil)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	toDelete := []uuid.UUID{}

	if len(currentResumes) > 0 {
		if newFile == nil {
			toDelete = lo.Map(currentResumes, func(item *usecase.CandidateResumeDTO, _ int) uuid.UUID {
				return item.Resume.ID
			})
		} else {
			currentResume := currentResumes[0]

			if currentResume.File.ID == newFile.ID {
				return nil
			}

			toDelete = append(toDelete, currentResume.Resume.ID)
		}
	}

	if newFile != nil {
		if newFile.FileMimetype == nil {
			return appErrors.Chainf(usecase.ErrInvalidCandidateResumeFileType, "%s.%s", uc.pkg, op)
		}

		newFileType, ok := usecase.MimetypesForFileTypeOnCandidateResume[*newFile.FileMimetype]
		if !ok {
			return appErrors.Chainf(usecase.ErrInvalidCandidateResumeFileType, "%s.%s", uc.pkg, op)
		}

		if newFile.FileHash == nil {
			return appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
		}

		newResume, err := entity.NewCandidateResume(candidate.TenantID, newFile.ID, *newFile.FileHash, newFileType)
		if err != nil {
			return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		err = newResume.SetCandidateID(&candidate.ID)
		if err != nil {
			return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		err = uc.candidateResumeUC.Create(ctx, newResume)
		if err != nil {
			return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}
	}

	if len(toDelete) > 0 {
		err := uc.candidateResumeUC.Delete(ctx, &usecase.CandidateResumeListOptions{
			FilterIDs: &toDelete,
		})
		if err != nil {
			return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}
	}

	return nil
}
