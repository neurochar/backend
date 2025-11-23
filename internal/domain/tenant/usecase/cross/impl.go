package tenant

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
	accoutUC       usecase.AccountUsecase
	tenantUC       usecase.TenantUsecase
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	accoutUC usecase.AccountUsecase,
	tenantUC usecase.TenantUsecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:            "Tenant.usecase.Tenant",
		logger:         logger,
		cfg:            cfg,
		dbMasterClient: dbMasterClient,
		accoutUC:       accoutUC,
		tenantUC:       tenantUC,
	}
	return uc
}
