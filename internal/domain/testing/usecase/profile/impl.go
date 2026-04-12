package profile

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/emailing"
)

type UsecaseImpl struct {
	pkg                string
	logger             *slog.Logger
	cfg                config.Config
	dbMasterClient     db.MasterClient
	emailing           emailing.Emailing
	repo               usecase.ProfileRepository
	tenantAccountUC    tenantUC.AccountUsecase
	personalityTraitUC usecase.PersonalityTraitUsecase
	llmRepo            usecase.LLMRepository
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	emailing emailing.Emailing,
	repo usecase.ProfileRepository,
	tenantAccountUC tenantUC.AccountUsecase,
	personalityTraitUC usecase.PersonalityTraitUsecase,
	llmRepo usecase.LLMRepository,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:                "Testing.Usecase.Profile",
		logger:             logger,
		cfg:                cfg,
		emailing:           emailing,
		dbMasterClient:     dbMasterClient,
		repo:               repo,
		llmRepo:            llmRepo,
		tenantAccountUC:    tenantAccountUC,
		personalityTraitUC: personalityTraitUC,
	}
	return uc
}
