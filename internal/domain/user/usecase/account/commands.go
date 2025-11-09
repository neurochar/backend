package account

import (
	"context"
	"errors"
	"net"
	"strings"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
	"github.com/neurochar/backend/internal/domain/user/usecase"
	"github.com/neurochar/backend/pkg/emailnormalize"
)

func (uc *UsecaseImpl) CreateAccountByDTO(
	ctx context.Context,
	in usecase.AccountDataInput,
	requestIP net.IP,
) (*userEntity.Account, *userEntity.AccountCode, error) {
	const op = "CreateAccountByDTO"

	var code *userEntity.AccountCode

	account, err := userEntity.NewAccount(in.Email, in.Password, in.RoleID, in.IsEmailVerified)
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

			checkAccount, err := uc.repoAccount.FindOneByEmail(ctx, account.Email, nil)
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
			_, err := uc.roleUC.GetRoleByID(ctx, *in.RoleID)
			if err != nil {
				if errors.Is(err, appErrors.ErrNotFound) {
					return usecase.ErrRoleNotFound
				}
				return err
			}

			account.RoleID = *in.RoleID
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

func (uc *UsecaseImpl) UpdateAccount(ctx context.Context, item *userEntity.Account) error {
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
	email string,
	requestIP net.IP,
) (*userEntity.AccountCode, error) {
	const op = "RequestPasswordRecoveryByEmail"

	var code *userEntity.AccountCode

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		res, err := emailnormalize.Normalize(strings.TrimSpace(email))
		if err != nil {
			return appErrors.ErrBadRequest
		}

		email = res.NormalizedAddress

		account, err := uc.repoAccount.FindOneByEmail(ctx, email, nil)
		if err != nil {
			return err
		}

		code, err = uc.createAccountPasswordRecoveryCode(ctx, account, requestIP)
		if err != nil {
			return err
		}

		err = uc.sendRecoveryCodeEmailToUser(ctx, account, code, requestIP)
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
) (*userEntity.AccountCode, error) {
	const op = "CheckCode"

	code, err := uc.checkCodeByID(ctx, codeID, codeValue)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return code, nil
}

func (uc *UsecaseImpl) UpdateCode(
	ctx context.Context,
	code *userEntity.AccountCode,
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
) (*userEntity.AccountCode, error) {
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
