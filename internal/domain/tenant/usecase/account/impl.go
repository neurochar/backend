package account

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	"github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/emailing"
	"github.com/neurochar/backend/internal/infra/imageproc"
)

type UsecaseImpl struct {
	pkg             string
	logger          *slog.Logger
	cfg             config.Config
	dbMasterClient  db.MasterClient
	emailing        emailing.Emailing
	imageProc       imageproc.ImageProcessor
	repoAccount     usecase.AccountRepository
	repoAccountCode usecase.AccountCodeRepository
	tenantUC        usecase.TenantUsecase
	sessionUC       usecase.SessionUsecase
	fileUC          fileUC.Usecase
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	emailing emailing.Emailing,
	imageProc imageproc.ImageProcessor,
	repoAccount usecase.AccountRepository,
	repoAccountCode usecase.AccountCodeRepository,
	tenantUC usecase.TenantUsecase,
	sessionUC usecase.SessionUsecase,
	fileUC fileUC.Usecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:             "TenantUser.Usecase.Account",
		logger:          logger,
		cfg:             cfg,
		emailing:        emailing,
		imageProc:       imageProc,
		dbMasterClient:  dbMasterClient,
		repoAccount:     repoAccount,
		repoAccountCode: repoAccountCode,
		tenantUC:        tenantUC,
		sessionUC:       sessionUC,
		fileUC:          fileUC,
	}
	return uc
}
