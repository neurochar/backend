package tenant

import (
	"context"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/auth"
)

func (uc *UsecaseImpl) PatchTenantByDTO(
	ctx context.Context,
	id uuid.UUID,
	in usecase.PatchTenantDataInput,
	skipVersionCheck bool,
) error {
	const op = "PatchTenantByDTO"

	if auth.IsNeedToCheckRights(ctx) {
		authData := auth.GetAuthData(ctx)
		if authData == nil || authData.TenantID != id {
			return appErrors.Chainf(appErrors.ErrForbidden, "%s.%s", uc.pkg, op)
		}

		authorAccount, err := uc.accoutUC.FindOneByID(ctx, authData.AccountID, nil, &usecase.AccountDTOOptions{})
		if err != nil {
			return err
		}

		if authorAccount.Role.Rank < 1 {
			return appErrors.Chainf(appErrors.ErrBadRequest.WithHints("forbidden"), "%s.%s", uc.pkg, op)
		}
	}

	err := uc.tenantUC.PatchByDTO(ctx, id, in, skipVersionCheck)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}
