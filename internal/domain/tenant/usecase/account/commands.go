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
	"github.com/neurochar/backend/internal/domain/tenant/entity"
	"github.com/neurochar/backend/internal/domain/tenant/roles"
	"github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/infra/imageproc"
	"github.com/neurochar/backend/internal/infra/loghandler"
	"github.com/neurochar/backend/pkg/auth"
	"github.com/neurochar/backend/pkg/emailnormalize"
	"github.com/samber/lo"
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

func (uc *UsecaseImpl) CreateAccountByDTO(
	ctx context.Context,
	tenantID uuid.UUID,
	in usecase.CreateAccountDataInput,
	sendStartEmail bool,
	requestIP *net.IP,
) (*usecase.AccountDTO, *entity.AccountCode, error) {
	const op = "CreateAccountByDTO"

	targetRole, ok := roles.RolesMap[in.RoleID]
	if !ok {
		return nil, nil, appErrors.Chainf(appErrors.ErrBadRequest.WithHints("roleID invalid"), "%s.%s", uc.pkg, op)
	}

	if auth.IsNeedToCheckTenantAccess(ctx) {
		authData := auth.GetAuthData(ctx)
		if authData == nil || !authData.IsTenantUser() || authData.TenantUserClaims().TenantID != tenantID {
			return nil, nil, appErrors.Chainf(appErrors.ErrForbidden, "%s.%s", uc.pkg, op)
		}

		authorAccount, err := uc.FindOneByID(ctx, authData.TenantUserClaims().AccountID, nil, &usecase.AccountDTOOptions{})
		if err != nil {
			return nil, nil, err
		}

		if targetRole.Rank <= authorAccount.Role.Rank {
			return nil, nil, appErrors.Chainf(appErrors.ErrBadRequest.WithHints("roleID value forbidden"), "%s.%s", uc.pkg, op)
		}
	}

	var code *entity.AccountCode

	account, err := entity.NewAccount(
		tenantID,
		in.Email,
		in.Password,
		in.SkipPasswordCheck,
		in.RoleID,
		in.IsConfirmed,
		in.IsEmailVerified,
	)
	if err != nil {
		return nil, nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = account.SetProfileName(in.ProfileName)
	if err != nil {
		return nil, nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = account.SetProfileSurname(in.ProfileSurname)
	if err != nil {
		return nil, nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		_, err := uc.repoAccount.FindOneByEmail(ctx, tenantID, account.Email, nil)
		if err == nil {
			return appErrors.ErrUniqueViolation.WithDetail("field", false, "email").WithHints("email is already in use")
		} else if !errors.Is(err, appErrors.ErrNotFound) {
			return err
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

	accountDTO, err := uc.entitiesToDTO(ctx, []*entity.Account{account}, nil)
	if err != nil {
		return nil, nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if len(accountDTO) == 0 {
		uc.logger.ErrorContext(loghandler.WithSource(ctx), "unpredicted empty account dto")
		return nil, nil, appErrors.Chainf(appErrors.ErrInternal, "%s.%s", uc.pkg, op)
	}

	if sendStartEmail {
		err = uc.sendStartEmailToUser(ctx, accountDTO[0], code, true, in.Password)
		if err != nil {
			return nil, nil, err
		}
	}

	return accountDTO[0], code, nil
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

		if auth.IsNeedToCheckTenantAccess(ctx) {
			authData := auth.GetAuthData(ctx)
			if authData == nil || !authData.IsTenantUser() || authData.TenantUserClaims().TenantID != account.TenantID {
				return appErrors.ErrForbidden
			}

			authorAccount, err := uc.FindOneByID(ctx, authData.TenantUserClaims().AccountID, nil, &usecase.AccountDTOOptions{})
			if err != nil {
				return err
			}

			if accountDTO.Role.Rank <= authorAccount.Role.Rank {
				return appErrors.ErrBadRequest.WithHints("roleID value forbidden")
			}
		}

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
			err = account.SetPassword(*in.Password, in.SkipPasswordCheck)
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
			if _, ok := roles.RolesMap[*in.RoleID]; !ok {
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
			allNil := lo.EveryBy([]*uuid.UUID{
				in.ProfilePhotos.PhotoOriginalFileID,
				in.ProfilePhotos.Photo100x100FileID,
			}, func(v *uuid.UUID) bool {
				return v == nil
			})

			allNotNil := lo.NoneBy([]*uuid.UUID{
				in.ProfilePhotos.PhotoOriginalFileID,
				in.ProfilePhotos.Photo100x100FileID,
			}, func(v *uuid.UUID) bool {
				return v == nil
			})

			if !allNil && !allNotNil {
				return appErrors.ErrBadRequest
			}

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

		if in.Password != nil {
			err = uc.sessionUC.RevokeSessionsByAccountID(ctx, id)
			if err != nil {
				return err
			}
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
	requestIP *net.IP,
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
	removeSessions bool,
) error {
	const op = "UpdatePasswordByRecoveryCode"

	code, err := uc.checkCodeByID(ctx, codeID, codeValue)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		account, err := uc.repoAccount.FindOneByID(ctx, code.AccountID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return appErrors.ErrInternal.WithParent(err)
		}

		err = account.SetPassword(newPassword, false)
		if err != nil {
			return err
		}

		err = uc.repoAccount.Update(ctx, account)
		if err != nil {
			return err
		}

		code.Deactivate()

		err = uc.repoAccountCode.Update(ctx, code)
		if err != nil {
			return err
		}

		if removeSessions {
			err = uc.sessionUC.RevokeSessionsByAccountID(ctx, code.AccountID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) SendStartEmailToUser(
	ctx context.Context,
	accountDTO *usecase.AccountDTO,
	activationCode *entity.AccountCode,
	isSendPassword bool,
	passwordForSend string,
) error {
	const op = "sendStartEmailToUser"

	err := uc.sendStartEmailToUser(ctx, accountDTO, activationCode, isSendPassword, passwordForSend)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}
