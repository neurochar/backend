package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
	"github.com/neurochar/backend/internal/domain/user/usecase"
	"github.com/neurochar/backend/internal/infra/emailing"
	"github.com/neurochar/backend/internal/infra/loghandler"
)

type StartEmailData struct {
	Cfg            config.Config
	IsAdminMode    bool
	Account        *userEntity.Account
	ProfileDTO     *usecase.FullProfileDTO
	IsSendPassword bool
	Password       string
	NeedActivate   bool
	ActivateUrl    string
}

func (uc *UsecaseImpl) sendStartEmailToUser(
	ctx context.Context,
	isAdminMode bool,
	account *userEntity.Account,
	activationCode *userEntity.AccountCode,
	profileDTO *usecase.FullProfileDTO,
	isSendPassword bool,
	passwordForSend string,
) error {
	const op = "sendStartEmailToUser"

	data := StartEmailData{
		Cfg:            uc.cfg,
		IsAdminMode:    isAdminMode,
		Account:        account,
		ProfileDTO:     profileDTO,
		IsSendPassword: isSendPassword,
		Password:       passwordForSend,
	}

	if activationCode != nil {
		data.NeedActivate = true
		if data.IsAdminMode {
			data.ActivateUrl = fmt.Sprintf(
				"%s%s/v1/users/accounts/verify-email?code_id=%s&code=%s",
				uc.cfg.Global.CpanelApiUrl,
				uc.cfg.CPanelApp.HTTP.Prefix,
				activationCode.ID,
				activationCode.Code,
			)
		}
	}

	msg, err := emailing.NewMessageFromJetTpl("users/email_on_new_user_created.jet", data)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.emailing.SendAsyc(emailing.Message{
		To:            account.Email,
		Subject:       fmt.Sprintf("Регистрация на сайте %s", uc.cfg.Global.ProjectName),
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
