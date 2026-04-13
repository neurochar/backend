package room

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	candidateUC "github.com/neurochar/backend/internal/domain/crm/usecase"
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
	repo               usecase.RoomRepository
	repoLLM            usecase.LLMRepository
	candidateUC        candidateUC.CandidateUsecase
	tenantAccountUC    tenantUC.AccountUsecase
	personalityTraitUC usecase.PersonalityTraitUsecase
	profileUC          usecase.ProfileUsecase
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	emailing emailing.Emailing,
	repo usecase.RoomRepository,
	repoLLM usecase.LLMRepository,
	candidateUC candidateUC.CandidateUsecase,
	tenantAccountUC tenantUC.AccountUsecase,
	personalityTraitUC usecase.PersonalityTraitUsecase,
	profileUC usecase.ProfileUsecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:                "Testing.Usecase.Room",
		logger:             logger,
		cfg:                cfg,
		emailing:           emailing,
		dbMasterClient:     dbMasterClient,
		repo:               repo,
		repoLLM:            repoLLM,
		candidateUC:        candidateUC,
		tenantAccountUC:    tenantAccountUC,
		personalityTraitUC: personalityTraitUC,
		profileUC:          profileUC,
	}
	return uc
}
