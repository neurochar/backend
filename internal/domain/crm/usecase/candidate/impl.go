package candidate

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/domain/crm/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/emailing"
)

type UsecaseImpl struct {
	pkg             string
	logger          *slog.Logger
	cfg             config.Config
	dbMasterClient  db.MasterClient
	emailing        emailing.Emailing
	repo            usecase.CandidateRepository
	tenantAccountUC tenantUC.AccountUsecase
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	emailing emailing.Emailing,
	repo usecase.CandidateRepository,
	tenantAccountUC tenantUC.AccountUsecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:             "CRM.Usecase.Candidate",
		logger:          logger,
		cfg:             cfg,
		emailing:        emailing,
		dbMasterClient:  dbMasterClient,
		repo:            repo,
		tenantAccountUC: tenantAccountUC,
	}
	return uc
}
