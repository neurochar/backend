package account

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/domain/tenant_user/usecase"
	"github.com/samber/lo"
)

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
