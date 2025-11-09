package user

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	"github.com/neurochar/backend/internal/domain/user/usecase"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/emailing"
)

type UsecaseImpl struct {
	pkg                string
	logger             *slog.Logger
	cfg                config.Config
	dbMasterClient     db.MasterClient
	emailing           emailing.Emailing
	repoProfileAccount usecase.ProfileAccountRepository
	fileUC             fileUC.Usecase
	accountUC          usecase.AccountUsecase
	profileUC          usecase.ProfileUsecase
	roleUC             usecase.RoleUsecase
	adminAuthUC        usecase.AdminAuthUsecase
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	emailing emailing.Emailing,
	repoProfileAccount usecase.ProfileAccountRepository,
	fileUC fileUC.Usecase,
	accountUC usecase.AccountUsecase,
	profileUC usecase.ProfileUsecase,
	roleUC usecase.RoleUsecase,
	adminAuthUC usecase.AdminAuthUsecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:                "User.usecase.Common",
		logger:             logger,
		cfg:                cfg,
		dbMasterClient:     dbMasterClient,
		emailing:           emailing,
		repoProfileAccount: repoProfileAccount,
		fileUC:             fileUC,
		accountUC:          accountUC,
		profileUC:          profileUC,
		roleUC:             roleUC,
		adminAuthUC:        adminAuthUC,
	}
	return uc
}
