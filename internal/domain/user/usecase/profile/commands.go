package profile

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
	"github.com/neurochar/backend/internal/domain/user/usecase"
	"github.com/neurochar/backend/internal/infra/imageproc"
)

func (uc *UsecaseImpl) UploadProfileImageFile(
	ctx context.Context,
	fileName string,
	fileData []byte,
) (usecase.UploadFileOut, error) {
	const op = "UploadProfileImageFile"

	filesMap, _, err := uc.fileUC.UploadAndCreateFiles(ctx, fileUC.UploadFilesIn{
		{
			FileData: fileData,
			Target:   string(usecase.FileTargetProfilePhoto100x100),
			FileName: fileName,
			Process: func(fd []byte) ([]byte, *appErrors.AppError) {
				return uc.imageProc.ScaleAndCrop(fileData, 100, 100, imageproc.WithAllowedFormats(imageproc.FormatJPEG))
			},
		},
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	result := make(usecase.UploadFileOut, 0, len(filesMap))

	for _, file := range filesMap {
		result = append(result, file)
	}

	return result, nil
}

func (uc *UsecaseImpl) CreateByDTO(
	ctx context.Context,
	account *userEntity.Account,
	in usecase.ProfileDataInput,
) (*usecase.FullProfileDTO, error) {
	const op = "CreateByDTO"

	profile, err := userEntity.NewProfile(account.ID, in.Name, in.Surname)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		_, err = uc.fileUC.ProcessFilesToTarget(ctx, []fileUC.ProcessFileToTargetIn{
			{
				CurrentFileID: profile.Photo100x100FileID,
				NewFileID:     in.Photo100x100FileID,
				Target:        string(usecase.FileTargetProfilePhoto100x100),
			},
		})
		if err != nil {
			return err
		}

		profile.Photo100x100FileID = in.Photo100x100FileID

		err = uc.repo.Create(ctx, profile)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	res, err := uc.FindFullList(ctx, &usecase.ProfileListOptions{
		AccountID: &account.ID,
	}, &uctypes.QueryGetListParams{
		Limit: 1,
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(res) == 0 {
		return nil, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	return res[0], nil
}

func (uc *UsecaseImpl) UpdateByDTO(ctx context.Context, id uint64, in usecase.ProfileDataInput, skipVersionCheck bool) error {
	const op = "UpdateByDTO"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		profile, err := uc.repo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if !skipVersionCheck && profile.Version() != in.Version {
			return appErrors.ErrVersionConflict.
				WithDetail("last_version", false, profile.Version()).
				WithDetail("last_updated_at", false, profile.UpdatedAt)
		}

		_, err = uc.fileUC.ProcessFilesToTarget(ctx, []fileUC.ProcessFileToTargetIn{
			{
				CurrentFileID: profile.Photo100x100FileID,
				NewFileID:     in.Photo100x100FileID,
				Target:        string(usecase.FileTargetProfilePhoto100x100),
			},
		})
		if err != nil {
			return err
		}

		err = profile.SetName(in.Name)
		if err != nil {
			return err
		}

		err = profile.SetSurname(in.Surname)
		if err != nil {
			return err
		}

		profile.Photo100x100FileID = in.Photo100x100FileID

		err = uc.repo.Update(ctx, profile)
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
