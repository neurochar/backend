package user

import (
	"context"
	"errors"
	"net"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/user/usecase"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) CreateUser(
	ctx context.Context,
	in usecase.CreateUserInput,
	requestIP net.IP,
	isAdminMode bool,
) (*usecase.UserDTO, error) {
	const op = "CreateUser"

	var role *usecase.RoleDTO

	out := &usecase.UserDTO{}

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		var err error

		role, err = uc.roleUC.GetRoleByID(ctx, in.Account.RoleID)
		if err != nil {
			if errors.Is(err, appErrors.ErrNotFound) {
				return usecase.ErrRoleNotFound
			}
			return err
		}

		out.Role = role

		account, activationCode, err := uc.accountUC.CreateAccountByDTO(ctx, in.Account, requestIP)
		if err != nil {
			return err
		}

		out.Account = account

		profileDTO, err := uc.profileUC.CreateByDTO(ctx, account, in.Profile)
		if err != nil {
			return err
		}

		out.ProfileDTO = profileDTO

		err = uc.sendStartEmailToUser(
			ctx,
			isAdminMode,
			account,
			activationCode,
			profileDTO,
			in.IsSendPassword,
			in.Account.Password,
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return out, nil
}

func (uc *UsecaseImpl) PatchAccountByDTO(
	ctx context.Context,
	id uuid.UUID,
	in usecase.PatchAccountDataInput,
	removeSessions bool,
	skipVersionCheck bool,
) error {
	const op = "PatchAccountByDTO"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		err := uc.accountUC.PatchAccountByDTO(ctx, id, in, skipVersionCheck)
		if err != nil {
			return err
		}

		if removeSessions {
			err = uc.DeleteAccountActiveSessions(ctx, id)
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

func (uc *UsecaseImpl) DeleteAccountActiveSessions(
	ctx context.Context,
	accountID uuid.UUID,
) (resErr error) {
	const op = "DeleteAccountActiveSessions"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		err := uc.adminAuthUC.DeleteActiveSessionsByAccountID(ctx, accountID)
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

func (uc *UsecaseImpl) DeleteRole(ctx context.Context, roleID uint64) error {
	const op = "DeleteRole"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		check, err := uc.accountUC.FindList(ctx, &usecase.AccountListOptions{
			RoleID: &roleID,
		}, &uctypes.QueryGetListParams{
			Limit: 1,
		})
		if err != nil {
			return err
		}

		if len(check) > 0 {
			return usecase.ErrCantDeleteRoleAccountsExists
		}

		err = uc.roleUC.DeleteRole(ctx, roleID)
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

func (uc *UsecaseImpl) UpdatePasswordByRecoveryCode(
	ctx context.Context,
	codeID uuid.UUID,
	codeValue string,
	newPassword string,
	removeSessions bool,
) error {
	const op = "UpdatePasswordByRecoveryCode"

	code, err := uc.accountUC.CheckCode(ctx, codeID, codeValue)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		err := uc.accountUC.PatchAccountByDTO(ctx, code.AccountID, usecase.PatchAccountDataInput{
			Password: lo.ToPtr(newPassword),
		}, true)
		if err != nil {
			return err
		}

		code.Deactivate()

		err = uc.accountUC.UpdateCode(ctx, code)
		if err != nil {
			return err
		}

		if removeSessions {
			err = uc.DeleteAccountActiveSessions(ctx, code.AccountID)
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
