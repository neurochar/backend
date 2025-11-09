package usecase

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
)

type UsecaseImpl struct {
	pkg    string
	logger *slog.Logger
	cfg    config.Config
	repo   TelegramRepository
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	repo TelegramRepository,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:    "Alerts.usecase",
		logger: logger,
		cfg:    cfg,
		repo:   repo,
	}
	return uc
}
