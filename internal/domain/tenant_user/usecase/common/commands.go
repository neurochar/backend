package account

import (
	"context"
	"net"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/domain/tenant_user/entity"
	"github.com/neurochar/backend/internal/domain/tenant_user/usecase"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) CreateUser(
	ctx context.Context,
	tenantID uuid.UUID,
	in usecase.CreateAccountDataInput,
	author *usecase.AccountDTO,
	requestIP net.IP,
) (*usecase.AccountDTO, error) {
	const op = "CreateUser"

	var accountDTO *usecase.AccountDTO

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		var err error
		var activationCode *entity.AccountCode

		accountDTO, activationCode, err = uc.accountUC.CreateAccountByDTO(ctx, tenantID, in, author, requestIP)
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
