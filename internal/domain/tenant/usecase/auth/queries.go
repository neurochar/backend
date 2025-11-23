package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/domain/tenant/usecase"
)

func (uc *UsecaseImpl) IsSessionConfirmed(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	const op = "IsSessionConfirmed"

	session, err := uc.sessionUC.FindOneByID(ctx, id, nil)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return false, nil
		}
		return false, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	accountDTO, err := uc.accountUC.FindOneByID(ctx, session.AccountID, nil, &usecase.AccountDTOOptions{})
	if err != nil {
		return false, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	if !accountDTO.Account.IsConfirmed || accountDTO.Account.IsBlocked {
		return false, nil
	}

	return true, nil
}
