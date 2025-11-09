package auth

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/infra/db"

	"github.com/neurochar/backend/internal/domain/tenant_user/usecase"
)

type UsecaseImpl struct {
	pkg            string
	logger         *slog.Logger
	cfg            config.Config
	repo           usecase.SessionRepository
	dbMasterClient db.MasterClient
	accountUC      usecase.AccountUsecase
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	repo usecase.SessionRepository,
	accountUC usecase.AccountUsecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:            "TenantUser.usecase.Auth",
		logger:         logger,
		cfg:            cfg,
		dbMasterClient: dbMasterClient,
		repo:           repo,
		accountUC:      accountUC,
	}

	return uc
}
