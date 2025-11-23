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
	accountUC      usecase.AccountUsecase
	sessionUC      usecase.SessionUsecase
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	accountUC usecase.AccountUsecase,
	sessionUC usecase.SessionUsecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:            "TenantUser.usecase.Auth",
		logger:         logger,
		cfg:            cfg,
		dbMasterClient: dbMasterClient,
		accountUC:      accountUC,
		sessionUC:      sessionUC,
	}

	return uc
}
