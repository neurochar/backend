package profile

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	"github.com/neurochar/backend/internal/domain/user/usecase"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/imageproc"
)

type UsecaseImpl struct {
	pkg            string
	logger         *slog.Logger
	cfg            config.Config
	imageProc      imageproc.ImageProcessor
	repo           usecase.ProfileRepository
	dbMasterClient db.MasterClient
	accountUC      usecase.AccountUsecase
	fileUC         fileUC.Usecase
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	imageProc imageproc.ImageProcessor,
	dbMasterClient db.MasterClient,
	repo usecase.ProfileRepository,
	accountUC usecase.AccountUsecase,
	fileUC fileUC.Usecase,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:            "User.usercase.Profile",
		logger:         logger,
		cfg:            cfg,
		imageProc:      imageProc,
		dbMasterClient: dbMasterClient,
		repo:           repo,
		accountUC:      accountUC,
		fileUC:         fileUC,
	}
	return uc
}
