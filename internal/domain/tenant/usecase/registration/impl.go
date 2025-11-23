package registration

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/emailing"
)

type UsecaseImpl struct {
	pkg                 string
	logger              *slog.Logger
	cfg                 config.Config
	dbMasterClient      db.MasterClient
	emailing            emailing.Emailing
	repo                usecase.RegistrationRepository
	tenantUC            usecase.TenantUsecase
	tenantUserAccountUC usecase.AccountUsecase
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	emailing emailing.Emailing,
	repo usecase.RegistrationRepository,
	tenantUC usecase.TenantUsecase,
	tenantUserAccountUC usecase.AccountUsecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:                 "Tenant.usecase.Registration",
		logger:              logger,
		cfg:                 cfg,
		dbMasterClient:      dbMasterClient,
		emailing:            emailing,
		repo:                repo,
		tenantUC:            tenantUC,
		tenantUserAccountUC: tenantUserAccountUC,
	}
	return uc
}
