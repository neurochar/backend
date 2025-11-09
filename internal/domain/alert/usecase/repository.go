package usecase

import (
	"context"

	"github.com/neurochar/backend/internal/domain/alert/entity"
)

type TelegramRepository interface {
	SendMessage(ctx context.Context, alert *entity.Alert) error
}
