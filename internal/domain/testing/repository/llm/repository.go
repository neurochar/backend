package llm

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/openai/openai-go/v3"
)

type Repository struct {
	pkg          string
	cfg          config.Config
	logger       *slog.Logger
	openaiClient openai.Client
}

func NewRepository(logger *slog.Logger, cfg config.Config, openaiClient openai.Client) *Repository {
	return &Repository{
		pkg:          "Tesing.repository.LLM",
		cfg:          cfg,
		logger:       logger,
		openaiClient: openaiClient,
	}
}
