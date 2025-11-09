package tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/domain/alert/entity"
)

func (r *Repository) SendMessage(ctx context.Context, alert *entity.Alert) error {
	const op = "SendMessage"

	text := alert.Message

	message := tgbotapi.NewMessage(r.cfg.Alerts.TargetChannel, text)

	_, err := r.bot.Send(message)
	if err != nil {
		return appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	return nil
}
