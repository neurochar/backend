package account

import (
	"context"
	"net"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/tenant_user/constants"
	"github.com/neurochar/backend/internal/domain/tenant_user/entity"
	"github.com/neurochar/backend/internal/domain/tenant_user/usecase"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) CreateUser(
	ctx context.Context,
	tenantID uuid.UUID,
	in usecase.CreateAccountDataInput,
	authorID uuid.UUID,
	requestIP net.IP,
) (*usecase.AccountDTO, error) {
	const op = "CreateUser"

	targetRole, ok := constants.RolesMap[in.RoleID]
	if !ok {
		return nil, appErrors.Chainf(appErrors.ErrBadRequest.WithHints("roleID invalid"), "%s.%s", uc.pkg, op)
	}

	var accountDTO *usecase.AccountDTO

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		author, err := uc.accountUC.FindOneByID(ctx, authorID, nil, &usecase.AccountDTOOptions{})
		if err != nil {
			return err
		}

		if targetRole.Rank <= author.Role.Rank {
			return appErrors.Chainf(appErrors.ErrBadRequest.WithHints("roleID value forbidden"), "%s.%s", uc.pkg, op)
		}

		var activationCode *entity.AccountCode

		accountDTO, activationCode, err = uc.accountUC.CreateAccountByDTO(ctx, tenantID, in, requestIP)
		if err != nil {
			return err
		}

		err = uc.sendStartEmailToUser(
			ctx,
			accountDTO,
			activationCode,
			true,
			in.Password,
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return accountDTO, nil
}

func (uc *UsecaseImpl) PatchAccountByDTO(
	ctx context.Context,
	id uuid.UUID,
	in usecase.PatchAccountDataInput,
	authorID uuid.UUID,
	skipVersionCheck bool,
) error {
	const op = "PatchAccountByDTO"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		targetAccount, err := uc.accountUC.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		}, &usecase.AccountDTOOptions{})
		if err != nil {
			return err
		}

		author, err := uc.accountUC.FindOneByID(ctx, authorID, nil, &usecase.AccountDTOOptions{})
		if err != nil {
			return err
		}

		if targetAccount.Role.Rank <= author.Role.Rank {
			return appErrors.Chainf(appErrors.ErrBadRequest.WithHints("roleID value forbidden"), "%s.%s", uc.pkg, op)
		}

		err = uc.accountUC.PatchAccountByDTO(ctx, id, in, skipVersionCheck)
		if err != nil {
			return err
		}

		if in.Password != nil {
			err = uc.authUC.RevokeSessionsByAccountID(ctx, id)
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
			err = uc.authUC.RevokeSessionsByAccountID(ctx, code.AccountID)
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
