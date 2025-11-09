package file

import (
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/neurochar/backend/internal/infra/db"
)

type Repository struct {
	pkg      string
	logger   *slog.Logger
	pgClient db.MasterClient
	qb       squirrel.StatementBuilderType
}

func NewRepository(logger *slog.Logger, pgClient db.MasterClient) *Repository {
	return &Repository{
		pkg:      "File.repository.File",
		logger:   logger,
		pgClient: pgClient,
		qb:       squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
