package account

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/domain/tenant_user/usecase"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/emailing"
)

type UsecaseImpl struct {
	pkg             string
	logger          *slog.Logger
	cfg             config.Config
	dbMasterClient  db.MasterClient
	emailing        emailing.Emailing
	repoAccount     usecase.AccountRepository
	repoAccountCode usecase.AccountCodeRepository
	tenantUC        tenantUC.TenantUsecase
	accountUC       usecase.AccountUsecase
	authUC          usecase.AuthUsecase
	fileUC          fileUC.Usecase
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	emailing emailing.Emailing,
	repoAccount usecase.AccountRepository,
	repoAccountCode usecase.AccountCodeRepository,
	tenantUC tenantUC.TenantUsecase,
	accountUC usecase.AccountUsecase,
	authUC usecase.AuthUsecase,
	fileUC fileUC.Usecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:             "TenantUser.Usecase.Common",
		logger:          logger,
		cfg:             cfg,
		emailing:        emailing,
		dbMasterClient:  dbMasterClient,
		repoAccount:     repoAccount,
		repoAccountCode: repoAccountCode,
		tenantUC:        tenantUC,
		accountUC:       accountUC,
		authUC:          authUC,
		fileUC:          fileUC,
	}
	return uc
}
