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
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	emailing emailing.Emailing,
	repo usecase.ProfileRepository,
	tenantAccountUC tenantUC.AccountUsecase,
	personalityTraitUC usecase.PersonalityTraitUsecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:                "Tesing.Usecase.Profile",
		logger:             logger,
		cfg:                cfg,
		emailing:           emailing,
		dbMasterClient:     dbMasterClient,
		repo:               repo,
		tenantAccountUC:    tenantAccountUC,
		personalityTraitUC: personalityTraitUC,
	}
	return uc
}
