package profileaccount

import (
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/neurochar/backend/internal/infra/db"
)

type Repository struct {
	pkg       string
	logger    *slog.Logger
	pgClient  db.MasterClient
	pgScanApi *pgxscan.API
	qb        squirrel.StatementBuilderType
}

func NewRepository(logger *slog.Logger, pgClient db.MasterClient) *Repository {
	dbAPI, err := pgxscan.NewDBScanAPI(
		dbscan.WithColumnSeparator("___"),
	)
	if err != nil {
		panic(err)
	}

	pgScanApi, err := pgxscan.NewAPI(dbAPI)
	if err != nil {
		panic(err)
	}

	return &Repository{
		pkg:       "User.repository.Profile__Account",
		logger:    logger,
		pgClient:  pgClient,
		pgScanApi: pgScanApi,
		qb:        squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
