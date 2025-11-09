package adminauth

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/infra/db"

	"github.com/neurochar/backend/internal/domain/user/usecase"
)

type UsecaseImpl struct {
	pkg            string
	logger         *slog.Logger
	cfg            config.Config
	repo           usecase.SessionRepository
	dbMasterClient db.MasterClient
	accountUC      usecase.AccountUsecase
	roleUC         usecase.RoleUsecase
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	repo usecase.SessionRepository,
	accountUC usecase.AccountUsecase,
	roleUC usecase.RoleUsecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:            "User.usecase.AdminAuth",
		logger:         logger,
		cfg:            cfg,
		dbMasterClient: dbMasterClient,
		repo:           repo,
		accountUC:      accountUC,
		roleUC:         roleUC,
	}

	return uc
}
