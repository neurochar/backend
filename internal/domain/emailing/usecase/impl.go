package usecase

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/emailing"
	"github.com/neurochar/backend/internal/infra/storage"
)

type UsecaseImpl struct {
	pkg            string
	logger         *slog.Logger
	cfg            config.Config
	storageClient  storage.Client
	repo           ItemRepository
	dbMasterClient db.MasterClient
	emailing       emailing.Emailing
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	storageClient storage.Client,
	dbMasterClient db.MasterClient,
	repo ItemRepository,
	emailing emailing.Emailing,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:            "Emailing.usecase",
		logger:         logger,
		cfg:            cfg,
		storageClient:  storageClient,
		dbMasterClient: dbMasterClient,
		repo:           repo,
		emailing:       emailing,
	}
	return uc
}
