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

func (uc *UsecaseImpl) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	const op = "DeleteProfile"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		candidate, err := uc.profileRepo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
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

		err = uc.profileRepo.Update(ctx, candidate)
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

func (uc *UsecaseImpl) DeleteRoom(ctx context.Context, id uuid.UUID) error {
	const op = "DeleteRoom"

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		room, err := uc.roomRepo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if auth.IsNeedToCheckTenantAccess(ctx) {
			authData := auth.GetAuthData(ctx)
			if authData == nil || !authData.IsTenantUser() || authData.TenantUserClaims().TenantID != room.TenantID {
				return appErrors.ErrForbidden
			}
		}

		room.DeletedAt = lo.ToPtr(time.Now())

		err = uc.roomRepo.Update(ctx, room)
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
