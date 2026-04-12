package candidate

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/domain/crm/usecase"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/emailing"
	"github.com/neurochar/backend/internal/infra/storage"
)

type UsecaseImpl struct {
	pkg               string
	logger            *slog.Logger
	cfg               config.Config
	dbMasterClient    db.MasterClient
	storageClient     storage.Client
	emailing          emailing.Emailing
	fileUC            fileUC.Usecase
	repo              usecase.CandidateRepository
	candidateResumeUC usecase.CandidateResumeUsecase
	tenantAccountUC   tenantUC.AccountUsecase
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	storageClient storage.Client,
	emailing emailing.Emailing,
	fileUC fileUC.Usecase,
	repo usecase.CandidateRepository,
	candidateResumeUC usecase.CandidateResumeUsecase,
	tenantAccountUC tenantUC.AccountUsecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:               "CRM.Usecase.Candidate",
		logger:            logger,
		cfg:               cfg,
		emailing:          emailing,
		storageClient:     storageClient,
		fileUC:            fileUC,
		dbMasterClient:    dbMasterClient,
		repo:              repo,
		candidateResumeUC: candidateResumeUC,
		tenantAccountUC:   tenantAccountUC,
	}
	return uc
}
