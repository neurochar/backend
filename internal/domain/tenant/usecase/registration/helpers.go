package registration

import (
	"context"
	"fmt"
	"log/slog"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/domain/tenant/entity"
	"github.com/neurochar/backend/internal/infra/emailing"
	"github.com/neurochar/backend/internal/infra/loghandler"
)

type RegistrationEmailData struct {
	ActivateUrl string
}

func (uc *UsecaseImpl) sendRegistrationEmailToUser(
	ctx context.Context,
	registration *entity.Registration,
) error {
	const op = "sendRegistrationEmailToUser"

	data := RegistrationEmailData{
		ActivateUrl: fmt.Sprintf(
			"%s/registration/%s",
			uc.cfg.Global.ProjectFrontendUrl,
			registration.ID.String(),
		),
	}

	msg, err := emailing.NewMessageFromJetTpl("tenants/email_on_new_created.jet", data)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.emailing.SendAsyc(emailing.Message{
		To:            registration.Email,
		Subject:       "Регистрация в Neurochar",
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
