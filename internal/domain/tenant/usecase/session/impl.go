package auth

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/infra/db"

	"github.com/neurochar/backend/internal/domain/tenant/usecase"
)

type UsecaseImpl struct {
	pkg            string
	logger         *slog.Logger
	cfg            config.Config
	dbMasterClient db.MasterClient
	repo           usecase.SessionRepository
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	repo usecase.SessionRepository,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:            "TenantUser.usecase.Session",
		logger:         logger,
		cfg:            cfg,
		dbMasterClient: dbMasterClient,
		repo:           repo,
	}

	return uc
}
