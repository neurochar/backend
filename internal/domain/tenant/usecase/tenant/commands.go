package tenant

import (
	"context"
	"errors"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/tenant/entity"
	"github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/auth"
)

func (uc *UsecaseImpl) CreateByDTO(
	ctx context.Context,
	in usecase.CreateTenantIn,
) (*entity.Tenant, error) {
	const op = "CreateByDTO"

	tenant, err := entity.NewTenant(
		in.TextID,
		in.Name,
		in.IsDemo,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	tenant.SetIsActive(in.IsActive)

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		_, err := uc.repo.FindOneByTextID(ctx, in.TextID, nil)
		if err != nil && !errors.Is(err, appErrors.ErrNotFound) {
			return err
		}
		if err == nil {
			return usecase.ErrTenantAlreadyExists
		}

		err = uc.repo.Create(ctx, tenant)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return tenant, nil
}

func (uc *UsecaseImpl) PatchByDTO(
	ctx context.Context,
	id uuid.UUID,
	in usecase.PatchTenantDataInput,
	skipVersionCheck bool,
) error {
	const op = "PatchByDTO"

	if auth.IsNeedToCheckRights(ctx) {
		authData := auth.GetAuthData(ctx)
		if authData == nil || authData.TenantID != id {
			return appErrors.Chainf(appErrors.ErrForbidden, "%s.%s", uc.pkg, op)
		}
	}

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		tenant, err := uc.repo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if in.Name != nil {
			err = tenant.SetName(*in.Name)
			if err != nil {
				return err
			}
		}

		err = uc.repo.Update(ctx, tenant)
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
