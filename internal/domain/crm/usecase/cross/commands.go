package cross

import (
	"context"
	"time"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/pkg/auth"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "Delete"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		candidate, err := uc.candidateRepo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if auth.IsNeedToCheckTenantAccess(ctx) {
			authData := auth.GetAuthData(ctx)
			if authData == nil || !authData.IsTenantUser() || authData.TenantUserClaims().TenantID != candidate.TenantID {
				return appErrors.ErrForbidden
			}
		}

		candidate.DeletedAt = lo.ToPtr(time.Now())

		err = uc.candidateRepo.Update(ctx, candidate)
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
