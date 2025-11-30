package personalitytrait

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
)

type UsecaseImpl struct {
	pkg    string
	logger *slog.Logger
	cfg    config.Config
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:    "Testing.Usecase.PersonalityTrait",
		logger: logger,
		cfg:    cfg,
	}
	return uc
}
