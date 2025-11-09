package usecase

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/storage"
)

type UsecaseImpl struct {
	pkg            string
	logger         *slog.Logger
	cfg            config.Config
	storageClient  storage.Client
	repo           FileRepository
	dbMasterClient db.MasterClient
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	storageClient storage.Client,
	dbMasterClient db.MasterClient,
	repo FileRepository,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:            "File.usecase",
		logger:         logger,
		cfg:            cfg,
		storageClient:  storageClient,
		dbMasterClient: dbMasterClient,
		repo:           repo,
	}
	return uc
}
