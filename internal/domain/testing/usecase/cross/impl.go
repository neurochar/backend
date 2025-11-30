package cross

import (
	"log/slog"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/domain/testing/usecase"
	"github.com/neurochar/backend/internal/infra/db"
	"github.com/neurochar/backend/internal/infra/emailing"
)

type UsecaseImpl struct {
	pkg            string
	logger         *slog.Logger
	cfg            config.Config
	dbMasterClient db.MasterClient
	emailing       emailing.Emailing
	profileRepo    usecase.ProfileRepository
	roomRepo       usecase.RoomRepository
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	emailing emailing.Emailing,
	profileRepo usecase.ProfileRepository,
	roomRepo usecase.RoomRepository,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:            "Testing.Usecase.Cross",
		logger:         logger,
		cfg:            cfg,
		emailing:       emailing,
		dbMasterClient: dbMasterClient,
		profileRepo:    profileRepo,
		roomRepo:       roomRepo,
	}
	return uc
}
