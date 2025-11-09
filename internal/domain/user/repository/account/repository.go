package account

import (
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/neurochar/backend/internal/infra/db"
)

// Repository - account repository
type Repository struct {
	pkg      string
	logger   *slog.Logger
	pgClient db.MasterClient
	qb       squirrel.StatementBuilderType
}

// NewRepository - constructor for account repository
func NewRepository(logger *slog.Logger, pgClient db.MasterClient) *Repository {
	return &Repository{
		pkg:      "User.repository.Account",
		logger:   logger,
		pgClient: pgClient,
		qb:       squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
