package candidate_resume

import (
	"context"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/crm/entity"
	"github.com/neurochar/backend/internal/domain/crm/usecase"
	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	"github.com/neurochar/backend/pkg/auth"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) UploadResumeFile(
	ctx context.Context,
	fileName string,
	fileData []byte,
) ([]*fileEntity.File, error) {
	const op = "UploadResumeFile"

	_, _, mimetype, _ := uc.storageClient.FileMetaByBytes(ctx, fileName, fileData)

	_, ok := usecase.MimetypesForFileTypeOnCandidateResume[mimetype]
	if !ok {
		return nil, appErrors.Chainf(usecase.ErrInvalidCandidateResumeFileType, "%s.%s", uc.pkg, op)
	}

	filesMap, _, err := uc.fileUC.UploadAndCreateFiles(ctx, fileUC.UploadFilesIn{
		{
			FileData: fileData,
			Target:   string(usecase.FileTargetCandidateResumeFile),
			FileName: fileName,
		},
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	result := make([]*fileEntity.File, 0, len(filesMap))

	for _, file := range filesMap {
		result = append(result, file)
	}

	return result, nil
}

func (uc *UsecaseImpl) Create(ctx context.Context, item *entity.CandidateResume) error {
	const op = "Create"

	err := uc.repo.Create(ctx, item)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) Update(ctx context.Context, item *entity.CandidateResume) error {
	const op = "Update"

	err := uc.repo.Update(ctx, item)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) PatchByDTO(
	ctx context.Context,
	id uuid.UUID,
	in usecase.PatchCandidateResumeDataInput,
	skipVersionCheck bool,
) error {
	const op = "PatchByDTO"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		item, err := uc.repo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if auth.IsNeedToCheckTenantAccess(ctx) {
			authData := auth.GetAuthData(ctx)
			if authData == nil || !authData.IsTenantUser() || authData.TenantUserClaims().TenantID != item.TenantID {
				return appErrors.ErrForbidden
			}
		}

		if !skipVersionCheck && item.Version() != in.Version {
			return appErrors.ErrVersionConflict.
				WithDetail("last_version", false, item.Version()).
				WithDetail("last_updated_at", false, item.UpdatedAt)
		}

		if in.Status != nil {
			err := item.SetStatus(*in.Status)
			if err != nil {
				return err
			}
		}

		if in.AnalyzeData != nil {
			err := item.SetAnalyzeData(*in.AnalyzeData)
			if err != nil {
				return err
			}
		}

		if in.ErrorText != nil {
			err := item.SetErrorText(*in.ErrorText)
			if err != nil {
				return err
			}
		}

		err = uc.repo.Update(ctx, item)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) Delete(ctx context.Context, selectReq *usecase.CandidateResumeListOptions) error {
	const op = "Delete"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		items, err := uc.repo.FindList(
			ctx,
			selectReq,
			&uctypes.QueryGetListParams{
				ForUpdate: true,
			})
		if err != nil {
			return err
		}

		if len(items) == 0 {
			return nil
		}

		err = uc.fileUC.DeleteByIDs(ctx, lo.Map(items, func(i *entity.CandidateResume, _ int) uuid.UUID {
			return i.FileID
		}))
		if err != nil {
			return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		tnow := time.Now()

		for _, item := range items {
			item.DeletedAt = &tnow
			err = uc.repo.Update(ctx, item)
			if err != nil {
				return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
			}
		}

		return nil
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}
