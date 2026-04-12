package providing

import (
	"github.com/neurochar/backend/internal/app/config"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func NewOpenAIClient(cfg config.Config) openai.Client {
	openaiClient := openai.NewClient(
		option.WithAPIKey(cfg.Openai.Token),
		option.WithBaseURL(cfg.Openai.BaseURL),
		option.WithMaxRetries(3),
	)

	return openaiClient
}
