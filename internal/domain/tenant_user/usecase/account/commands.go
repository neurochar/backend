package account

import (
	"context"
	"errors"
	"net"
	"strings"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	"github.com/neurochar/backend/internal/domain/tenant_user/constants"
	entity "github.com/neurochar/backend/internal/domain/tenant_user/entity"
	"github.com/neurochar/backend/internal/domain/tenant_user/usecase"
	"github.com/neurochar/backend/internal/infra/imageproc"
	"github.com/neurochar/backend/pkg/emailnormalize"
)

func (uc *UsecaseImpl) UploadProfileImageFile(
	ctx context.Context,
	fileName string,
	fileData []byte,
) ([]*fileEntity.File, error) {
	const op = "UploadProfileImageFile"

	filesMap, _, err := uc.fileUC.UploadAndCreateFiles(ctx, fileUC.UploadFilesIn{
		{
			FileData: fileData,
			Target:   string(usecase.FileTargetProfilePhotoOriginal),
			FileName: fileName,
			Process: func(fd []byte) ([]byte, *appErrors.AppError) {
				return uc.imageProc.DownscaleIfLarger(fileData, 1920, 1920, imageproc.WithAllowedFormats(imageproc.FormatJPEG))
			},
		},
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

	result := make([]*fileEntity.File, 0, len(filesMap))

	for _, file := range filesMap {
		result = append(result, file)
	}

	return result, nil
}

// CreateAccountByDTO TODO
func (uc *UsecaseImpl) CreateAccountByDTO(
	ctx context.Context,
	tenantID uuid.UUID,
	in usecase.CreateAccountDataInput,
	requestIP net.IP,
) (*entity.Account, *entity.AccountCode, error) {
	// const op = "CreateAccountByDTO"

	/*
		var code *entity.AccountCode

		account, err := entity.NewAccount(in.Email, in.Password, in.RoleID, in.IsEmailVerified)
		if err != nil {
			return nil, nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
			_, err := uc.repoAccount.FindOneByEmail(ctx, account.Email, nil)
			if err == nil {
				return appErrors.ErrUniqueViolation.WithDetail("field", false, "email")
			} else if !errors.Is(err, appErrors.ErrNotFound) {
				return err
			}

			_, err = uc.roleUC.GetRoleByID(ctx, account.RoleID)
			if err != nil {
				if errors.Is(err, appErrors.ErrNotFound) {
					return usecase.ErrRoleNotFound
				}
				return err
			}

			err = uc.repoAccount.Create(ctx, account)
			if err != nil {
				return err
			}

			if !in.IsEmailVerified {
				code, err = uc.createAccountEmailVerificationCode(ctx, account, requestIP)
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return nil, nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		return account, code, nil
	*/

	return nil, nil, nil
}

func (uc *UsecaseImpl) PatchAccountByDTO(
	ctx context.Context,
	id uuid.UUID,
	in usecase.PatchAccountDataInput,
	skipVersionCheck bool,
) error {
	const op = "PatchAccountByDTO"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		account, err := uc.repoAccount.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		accountDTOEntities, err := uc.entitiesToDTO(ctx, []*entity.Account{account}, &usecase.AccountDTOOptions{
			FetchTenant: true,
		})
		if err != nil {
			return err
		}

		accountDTO := accountDTOEntities[0]

		if !skipVersionCheck && account.Version() != in.Version {
			return appErrors.ErrVersionConflict.
				WithDetail("last_version", false, account.Version()).
				WithDetail("last_updated_at", false, account.UpdatedAt)
		}

		if in.Email != nil {
			err = account.SetEmail(*in.Email)
			if err != nil {
				return err
			}

			checkAccount, err := uc.repoAccount.FindOneByEmail(ctx, accountDTO.Tenant.ID, account.Email, nil)
			if err == nil && checkAccount.ID != account.ID {
				return appErrors.ErrUniqueViolation.WithDetail("field", false, "email")
			} else if err != nil && !errors.Is(err, appErrors.ErrNotFound) {
				return err
			}
		}

		if in.Password != nil {
			err = account.SetPassword(*in.Password)
			if err != nil {
				return err
			}
		}

		if in.IsEmailVerified != nil {
			account.IsEmailVerified = *in.IsEmailVerified
		}

		if in.IsBlocked != nil {
			account.IsBlocked = *in.IsBlocked
		}

		if in.RoleID != nil {
			if _, ok := constants.RolesMap[*in.RoleID]; !ok {
				return usecase.ErrRoleNotFound
			}

			err := account.SetRoleID(*in.RoleID)
			if err != nil {
				return err
			}
		}

		if in.ProfileName != nil {
			err := account.SetProfileName(*in.ProfileName)
			if err != nil {
				return err
			}
		}

		if in.ProfileSurname != nil {
			err := account.SetProfileSurname(*in.ProfileSurname)
			if err != nil {
				return err
			}
		}

		if in.ProfilePhotos != nil {
			_, err = uc.fileUC.ProcessFilesToTarget(ctx, []fileUC.ProcessFileToTargetIn{
				{
					CurrentFileID: account.ProfilePhotoOriginalFileID,
					NewFileID:     in.ProfilePhotos.PhotoOriginalFileID,
					Target:        string(usecase.FileTargetProfilePhotoOriginal),
					Group:         "profile_photo",
				},
				{
					CurrentFileID: account.ProfilePhoto100x100FileID,
					NewFileID:     in.ProfilePhotos.Photo100x100FileID,
					Target:        string(usecase.FileTargetProfilePhoto100x100),
					Group:         "profile_photo",
				},
			})
			if err != nil {
				return err
			}

			account.ProfilePhotoOriginalFileID = in.ProfilePhotos.PhotoOriginalFileID
			account.ProfilePhoto100x100FileID = in.ProfilePhotos.Photo100x100FileID
		}

		err = uc.repoAccount.Update(ctx, account)
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

func (uc *UsecaseImpl) UpdateAccount(ctx context.Context, item *entity.Account) error {
	const op = "UpdateAccount"

	err := uc.repoAccount.Update(ctx, item)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) VerifyAccountEmailByCode(ctx context.Context, codeID uuid.UUID, codeValue string) error {
	const op = "VerifyAccountEmailByCode"

	code, err := uc.checkCodeByID(ctx, codeID, codeValue)
	if err != nil {
		return err
	}

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		account, err := uc.repoAccount.FindOneByID(ctx, code.AccountID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return appErrors.ErrInternal.WithParent(err)
		}

		if account.IsEmailVerified {
			return appErrors.ErrBadRequest.WithHints("account email already verified")
		}

		account.IsConfirmed = true
		account.IsEmailVerified = true

		err = uc.repoAccount.Update(ctx, account)
		if err != nil {
			return err
		}

		code.Deactivate()

		err = uc.repoAccountCode.Update(ctx, code)
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

func (uc *UsecaseImpl) RequestPasswordRecoveryByEmail(
	ctx context.Context,
	tenantID uuid.UUID,
	email string,
	requestIP net.IP,
) (*entity.AccountCode, error) {
	const op = "RequestPasswordRecoveryByEmail"

	var code *entity.AccountCode

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		res, err := emailnormalize.Normalize(strings.TrimSpace(email))
		if err != nil {
			return appErrors.ErrBadRequest
		}

		email = res.NormalizedAddress

		account, err := uc.repoAccount.FindOneByEmail(ctx, tenantID, email, nil)
		if err != nil {
			return err
		}

		accountDTOEntities, err := uc.entitiesToDTO(ctx, []*entity.Account{account}, &usecase.AccountDTOOptions{
			FetchTenant: true,
		})
		if err != nil {
			return err
		}

		accountDTO := accountDTOEntities[0]

		code, err = uc.createAccountPasswordRecoveryCode(ctx, account, requestIP)
		if err != nil {
			return err
		}

		err = uc.sendRecoveryCodeEmailToUser(ctx, accountDTO, code, requestIP)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return code, nil
}

func (uc *UsecaseImpl) CheckCode(
	ctx context.Context,
	codeID uuid.UUID,
	codeValue string,
) (*entity.AccountCode, error) {
	const op = "CheckCode"

	code, err := uc.checkCodeByID(ctx, codeID, codeValue)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return code, nil
}

func (uc *UsecaseImpl) UpdateCode(
	ctx context.Context,
	code *entity.AccountCode,
) error {
	const op = "UpdateCode"

	err := uc.repoAccountCode.Update(ctx, code)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) UpdatePasswordByRecoveryCode(
	ctx context.Context,
	codeID uuid.UUID,
	codeValue string,
	newPassword string,
) (*entity.AccountCode, error) {
	const op = "UpdatePasswordByRecoveryCode"

	code, err := uc.checkCodeByID(ctx, codeID, codeValue)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		account, err := uc.repoAccount.FindOneByID(ctx, code.AccountID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return appErrors.ErrInternal.WithParent(err)
		}

		err = account.SetPassword(newPassword)
		if err != nil {
			return err
		}

		err = uc.repoAccount.Update(ctx, account)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return code, nil
}
