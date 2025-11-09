package user

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/infra/db"
)

type UsecaseImpl struct {
	pkg            string
	logger         *slog.Logger
	cfg            config.Config
	dbMasterClient db.MasterClient
	repo           usecase.TenantRepository
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	repo usecase.TenantRepository,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:            "Tenant.usecase.Tenant",
		logger:         logger,
		cfg:            cfg,
		dbMasterClient: dbMasterClient,
		repo:           repo,
	}
	return uc
}
