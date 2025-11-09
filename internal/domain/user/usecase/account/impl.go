package account

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/domain/user/usecase"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/emailing"
)

type UsecaseImpl struct {
	pkg             string
	logger          *slog.Logger
	cfg             config.Config
	emailing        emailing.Emailing
	repoAccount     usecase.AccountRepository
	repoAccountCode usecase.AccountCodeRepository
	dbMasterClient  db.MasterClient
	roleUC          usecase.RoleUsecase
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	emailing emailing.Emailing,
	dbMasterClient db.MasterClient,
	repoAccount usecase.AccountRepository,
	repoAccountCode usecase.AccountCodeRepository,
	roleUC usecase.RoleUsecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:             "User.Usercase.Account",
		logger:          logger,
		cfg:             cfg,
		emailing:        emailing,
		dbMasterClient:  dbMasterClient,
		repoAccount:     repoAccount,
		repoAccountCode: repoAccountCode,
		roleUC:          roleUC,
	}
	return uc
}
