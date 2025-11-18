package account

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	tenantEntity "github.com/neurochar/backend/internal/domain/tenant/entity"
	"github.com/neurochar/backend/internal/domain/tenant_user/entity"
	"github.com/neurochar/backend/internal/domain/tenant_user/usecase"
	"github.com/neurochar/backend/internal/infra/emailing"
	"github.com/neurochar/backend/internal/infra/loghandler"
)

type StartEmailData struct {
	Cfg            config.Config
	AccountDTO     *usecase.AccountDTO
	IsSendPassword bool
	Password       string
	NeedActivate   bool
	ActivateUrl    string
	Tenant         *tenantEntity.Tenant
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

	msg, err := emailing.NewMessageFromJetTpl("tenant_users/email_on_new_user_created.jet", data)
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
