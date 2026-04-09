package registration

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/tenant/entity"
	"github.com/neurochar/backend/internal/domain/tenant/roles"
	"github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/pkg/pgclient"
	"github.com/samber/lo"
	"github.com/sethvargo/go-password/password"
)

func (uc *UsecaseImpl) CreateByDTO(
	ctx context.Context,
	in usecase.CreateRegistrationIn,
	requestIP *netip.Addr,
) (*entity.Registration, error) {
	const op = "CreateByDTO"

	registration, err := entity.NewRegistration(
		in.Email,
		in.Tariff,
	)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = registration.SetRequestIP(requestIP)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.dbMasterClient.DoWithIsoLvl(ctx, pgclient.Serializable, func(ctx context.Context) error {
		if in.Tariff == 0 {
			checkOps := &usecase.RegistrationListOptions{
				FilterEmail:      &registration.Email,
				FilterIsFinished: lo.ToPtr(true),
			}

			check, err := uc.repo.FindList(ctx, checkOps, nil)
			if err != nil {
				return err
			}

			if len(check) > 0 {
				foundReg := check[0]
				foundTenantTextID := ""

				if foundReg.TenantID != nil {
					tenant, err := uc.tenantUC.FindOneByID(ctx, *foundReg.TenantID, nil)
					if err != nil {
						return err
					}

					foundTenantTextID = tenant.TextID
				}

				return usecase.ErrRegistrationDemoAlreadyExists.
					WithDetail("tenantTextID", false, foundTenantTextID)
			}
		}

		err := uc.repo.Create(ctx, registration)
		if err != nil {
			return err
		}

		err = uc.sendRegistrationEmailToUser(ctx, registration)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil, nil
}

func (uc *UsecaseImpl) FinishByDTO(
	ctx context.Context,
	id uuid.UUID,
	in usecase.FinishRegistrationIn,
	requestIP *netip.Addr,
) (*entity.Tenant, error) {
	const op = "FinishByDTO"

	pass, err := password.Generate(12, 0, 0, false, false)
	if err != nil {
		return nil, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", uc.pkg, op)
	}

	var tenant *entity.Tenant

	err = uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		registration, err := uc.repo.FindOneByID(ctx, id, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if registration.IsFinished {
			return usecase.ErrRegistrationAlreadyFinished
		}

		registration.IsFinished = true

		isDemo := registration.Tariff == 0

		tenant, err = uc.tenantUC.CreateByDTO(ctx, usecase.CreateTenantIn{
			Name:   fmt.Sprintf("%s %s", in.ProfileName, in.ProfileSurname),
			TextID: in.TenantTextID,
			IsDemo: isDemo,
		})
		if err != nil {
			return err
		}

		registration.TenantID = &tenant.ID

		ownerAccount, _, err := uc.tenantUserAccountUC.CreateAccountByDTO(
			ctx,
			tenant.ID,
			usecase.CreateAccountDataInput{
				Email:             registration.Email,
				Password:          pass,
				SkipPasswordCheck: true,
				RoleID:            roles.RoleCreator.ID,
				IsConfirmed:       true,
				IsEmailVerified:   true,
				ProfileName:       in.ProfileName,
				ProfileSurname:    in.ProfileSurname,
			},
			false,
			requestIP,
		)
		if err != nil {
			return err
		}

		err = uc.repo.Update(ctx, registration)
		if err != nil {
			return err
		}

		_, err = uc.repo.Delete(ctx, &usecase.RegistrationListOptions{
			FilterEmail:      &registration.Email,
			FilterIsFinished: lo.ToPtr(false),
		})
		if err != nil {
			return err
		}

		err = uc.tenantUserAccountUC.SendStartEmailToUser(ctx, ownerAccount, nil, true, pass)
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
