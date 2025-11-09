package tg

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/neurochar/backend/internal/app/config"
)

type Repository struct {
	pkg    string
	cfg    config.Config
	logger *slog.Logger
	bot    *tgbotapi.BotAPI
}

func NewRepository(cfg config.Config, logger *slog.Logger) *Repository {
	r := &Repository{
		cfg:    cfg,
		pkg:    "Alerts.repository.Telegram",
		logger: logger,
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Alerts.BotToken)
	if err != nil {
		logger.Error("failed to create telegram bot", slog.Any("error", err))
	}

	r.bot = bot

	return r
}
