package candidate_resume

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/domain/crm/usecase"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/storage"
	temporalClient "github.com/neurochar/backend/internal/infra/temporal/client"
)

type UsecaseImpl struct {
	pkg            string
	logger         *slog.Logger
	cfg            config.Config
	dbMasterClient db.MasterClient
	storageClient  storage.Client
	fileUC         fileUC.Usecase
	repo           usecase.CandidateResumeRepository
	temporalClient temporalClient.Client
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	storageClient storage.Client,
	fileUC fileUC.Usecase,
	repo usecase.CandidateResumeRepository,
	temporalClient temporalClient.Client,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:            "CRM.Usecase.CandidateResume",
		logger:         logger,
		cfg:            cfg,
		storageClient:  storageClient,
		fileUC:         fileUC,
		dbMasterClient: dbMasterClient,
		repo:           repo,
		temporalClient: temporalClient,
	}
	return uc
}
