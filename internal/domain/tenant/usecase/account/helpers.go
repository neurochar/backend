package account

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/app/config"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	"github.com/neurochar/backend/internal/domain/tenant/entity"
	"github.com/neurochar/backend/internal/domain/tenant/roles"
	"github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/infra/emailing"
	"github.com/neurochar/backend/internal/infra/loghandler"
	"github.com/samber/lo"
)

func (uc *UsecaseImpl) entitiesToDTO(
	ctx context.Context,
	items []*entity.Account,
	dtoOpts *usecase.AccountDTOOptions,
) ([]*usecase.AccountDTO, error) {
	const op = "entitiesToDTO"

	tenantsMap := make(map[uuid.UUID]*entity.Tenant, 0)
	filesMap := make(map[uuid.UUID]*fileEntity.File, 0)

	tenantsIDs := make([]uuid.UUID, 0)
	filesIDs := make([]uuid.UUID, 0)

	for _, item := range items {
		tenantsIDs = append(tenantsIDs, item.TenantID)
		filesIDs = append(filesIDs, item.FilesIDs()...)
	}

	if (dtoOpts == nil || dtoOpts.FetchTenant) && len(tenantsIDs) > 0 {
		tenantsList, err := uc.tenantUC.FindList(ctx, &usecase.TenantListOptions{
			FilterIDs: &tenantsIDs,
		}, nil)
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}

		tenantsMap = lo.SliceToMap(tenantsList, func(item *entity.Tenant) (uuid.UUID, *entity.Tenant) {
			return item.ID, item
		})
	}

	if (dtoOpts == nil || dtoOpts.FetchPhotoFiles) && len(filesIDs) > 0 {
		var err error
		filesMap, err = uc.fileUC.FindListInMap(ctx, &fileUC.ListOptions{
			IDs: &filesIDs,
		}, nil)
		if err != nil {
			return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
		}
	}

	out := make([]*usecase.AccountDTO, 0, len(items))

	for _, item := range items {
		resItem := &usecase.AccountDTO{
			Account: item,
		}

		role, ok := roles.RolesMap[item.RoleID]
		if !ok {
			return nil, appErrors.Chainf(appErrors.ErrInternal.Extend("unknown role"), "%s.%s", uc.pkg, op)
		}

		resItem.Role = role

		if dtoOpts == nil || dtoOpts.FetchTenant {
			tenant, ok := tenantsMap[item.TenantID]
			if !ok {
				return nil, appErrors.Chainf(appErrors.ErrInternal.Extend("tenant not fetched"), "%s.%s", uc.pkg, op)
			}

			resItem.Tenant = tenant
		}

		if dtoOpts == nil || dtoOpts.FetchPhotoFiles {
			if item.ProfilePhoto100x100FileID != nil {
				file, ok := filesMap[*item.ProfilePhoto100x100FileID]
				if !ok {
					return nil, appErrors.Chainf(appErrors.ErrInternal.Extend("file not fetched"), "%s.%s", uc.pkg, op)
				}

				resItem.ProfilePhoto100x100File = file
			}

			if item.ProfilePhotoOriginalFileID != nil {
				file, ok := filesMap[*item.ProfilePhotoOriginalFileID]
				if !ok {
					return nil, appErrors.Chainf(appErrors.ErrInternal.Extend("file not fetched"), "%s.%s", uc.pkg, op)
				}

				resItem.ProfilePhotoOriginalFile = file
			}
		}

		out = append(out, resItem)
	}

	return out, nil
}

type StartEmailData struct {
	Cfg            config.Config
	AccountDTO     *usecase.AccountDTO
	IsSendPassword bool
	Password       string
	NeedActivate   bool
	ActivateUrl    string
	Tenant         *entity.Tenant
	TenantURL      string
}

func (uc *UsecaseImpl) sendStartEmailToUser(
	ctx context.Context,
	accountDTO *usecase.AccountDTO,
	activationCode *entity.AccountCode,
	isSendPassword bool,
	passwordForSend string,
) error {
	const op = "sendStartEmailToUser"

	data := StartEmailData{
		Cfg:            uc.cfg,
		AccountDTO:     accountDTO,
		IsSendPassword: isSendPassword,
		Password:       passwordForSend,
		Tenant:         accountDTO.Tenant,
		TenantURL:      accountDTO.Tenant.GetUrl(uc.cfg.Global.TenantMainDomain, uc.cfg.Global.TenantUrlScheme),
	}

	if activationCode != nil {
		data.NeedActivate = true

		data.ActivateUrl = fmt.Sprintf(
			"%s/api/verify-email?code_id=%s&code=%s",
			accountDTO.Tenant.GetUrl(uc.cfg.Global.TenantMainDomain, uc.cfg.Global.TenantUrlScheme),
			activationCode.ID,
			activationCode.Code,
		)
	}

	msg, err := emailing.NewMessageFromJetTpl("tenants/email_on_new_user_created.jet", data)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.emailing.SendAsyc(emailing.Message{
		To:            accountDTO.Account.Email,
		Subject:       fmt.Sprintf("%s: Активация аккаунта", accountDTO.Tenant.TextID),
		TextHtml:      msg,
		AutoTextPlain: true,
	}, func(sendErr error) {
		if sendErr != nil {
			uc.logger.ErrorContext(loghandler.WithSource(ctx), "email sending", slog.Any("error", sendErr))
		}
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

func (uc *UsecaseImpl) createAccountEmailVerificationCode(
	ctx context.Context,
	account *entity.Account,
	requestIP *net.IP,
) (*entity.AccountCode, error) {
	const op = "createAccountEmailVerificationCode"

	code, err := entity.NewAccountCode(account.ID, entity.AccountCodeTypeEmailVerification, requestIP)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.repoAccountCode.Create(ctx, code)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return code, nil
}

func (uc *UsecaseImpl) createAccountPasswordRecoveryCode(
	ctx context.Context,
	account *entity.Account,
	requestIP *net.IP,
) (*entity.AccountCode, error) {
	const op = "createAccountPasswordRecoveryCode"

	code, err := entity.NewAccountCode(account.ID, entity.AccountCodeTypePasswordRecovery, requestIP)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.repoAccountCode.Create(ctx, code)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return code, nil
}

type RecoveryCodeEmailData struct {
	Cfg       config.Config
	Account   *entity.Account
	Tenant    *entity.Tenant
	TenantURL string
	Code      string
	RequestIP string
}

func (uc *UsecaseImpl) sendRecoveryCodeEmailToUser(
	ctx context.Context,
	accountDTO *usecase.AccountDTO,
	recoveryCode *entity.AccountCode,
	requestIP *net.IP,
) error {
	const op = "sendRecoveryCodeEmailToUser"

	data := RecoveryCodeEmailData{
		Cfg:       uc.cfg,
		Account:   accountDTO.Account,
		Tenant:    accountDTO.Tenant,
		TenantURL: accountDTO.Tenant.GetUrl(uc.cfg.Global.TenantMainDomain, uc.cfg.Global.TenantUrlScheme),
		Code:      recoveryCode.Code,
		RequestIP: requestIP.String(),
	}

	msg, err := emailing.NewMessageFromJetTpl("tenants/email_password_recovery_code.jet", data)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.emailing.SendAsyc(emailing.Message{
		To:            accountDTO.Account.Email,
		Subject:       fmt.Sprintf("%s: Восстановление пароля", accountDTO.Tenant.TextID),
		TextHtml:      msg,
		AutoTextPlain: true,
	}, func(sendErr error) {
		if sendErr != nil {
			uc.logger.ErrorContext(loghandler.WithSource(ctx), "email sending", slog.Any("error", sendErr))
		}
	})
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}

// Максимальное количество попыток проверки кода
const maxCodeAttempts = 3

// Максимальное время жизни кода
const maxCodeLifetime = time.Hour * 24 * 7

func (uc *UsecaseImpl) checkCodeByID(ctx context.Context, codeID uuid.UUID, codeValue string) (*entity.AccountCode, error) {
	const op = "checkCodeByID"

	var codeItem *entity.AccountCode
	var resErr error

	err := uc.dbMasterClient.Do(ctx, func(ctx context.Context) error {
		var err error

		codeItem, err = uc.repoAccountCode.FindOneByID(ctx, codeID, &uctypes.QueryGetOneParams{
			ForUpdate: true,
		})
		if err != nil {
			return err
		}

		if !codeItem.IsActive {
			resErr = appErrors.ErrBadRequest

			return nil
		}

		if !codeItem.VerifyCode(codeValue) {
			codeItem.AddAttempt()

			if codeItem.Attempts >= maxCodeAttempts {
				codeItem.Deactivate()
			}

			maxAttempts := maxCodeAttempts - codeItem.Attempts
			if maxAttempts < 0 {
				maxAttempts = 0
			}

			err := uc.repoAccountCode.Update(ctx, codeItem)
			if err != nil {
				return err
			}

			resErr = usecase.ErrCodeInvalid.WithDetail("max_attempts", false, maxAttempts)

			return nil
		}

		if !codeItem.IsAlive(maxCodeLifetime) {
			codeItem.Deactivate()

			err := uc.repoAccountCode.Update(ctx, codeItem)
			if err != nil {
				return err
			}

			resErr = usecase.ErrCodeExpired

			return nil
		}

		return nil
	})
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return codeItem, resErr
}
