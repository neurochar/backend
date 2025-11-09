package role

import (
	"log/slog"
	"sync"

	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/infra/db"

	"github.com/neurochar/backend/internal/domain/user/usecase"
)

type UsecaseImpl struct {
	pkg             string
	logger          *slog.Logger
	cfg             config.Config
	repoRole        usecase.RoleRepository
	repoRoleToRight usecase.RoleToRightRepository
	dbMasterClient  db.MasterClient

	rolesMap map[uint64]*usecase.RoleDTO
	mu       sync.RWMutex
}

func NewUsecaseImpl(
	logger *slog.Logger,
	cfg config.Config,
	dbMasterClient db.MasterClient,
	repoRole usecase.RoleRepository,
	repoRoleToRight usecase.RoleToRightRepository,
) *UsecaseImpl {
	uc := &UsecaseImpl{
		pkg:             "User.usercase.Role",
		logger:          logger,
		cfg:             cfg,
		dbMasterClient:  dbMasterClient,
		repoRole:        repoRole,
		repoRoleToRight: repoRoleToRight,
	}
	return uc
}
