package accountcode

import (
	"context"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/tenant_user/entity"
	"github.com/neurochar/backend/internal/domain/tenant_user/repository/pg"
	"github.com/neurochar/backend/internal/domain/tenant_user/usecase"
	"github.com/neurochar/backend/internal/infra/loghandler"
)

func (r *Repository) buildWhereForList(listOptions *usecase.AccountCodeListOptions, withDeleted bool) (where squirrel.And) {
	if listOptions == nil {
		return where
	}

	return where
}

func (r *Repository) buildSortForList(_ *usecase.AccountCodeListOptions) []string {
	return nil
}

func (r *Repository) FindList(
	ctx context.Context,
	listOptions *usecase.AccountCodeListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*entity.AccountCode, error) {
	const op = "FindList"

	withDeleted := queryParams != nil && queryParams.WithDeleted

	where := r.buildWhereForList(listOptions, withDeleted)

	fields := pg.AccountCodeTableFields

	q := r.qb.Select(fields...).From(pg.AccountCodeTable).Where(where)

	sort := r.buildSortForList(listOptions)
	if len(sort) > 0 {
		q = q.OrderBy(sort...)
	}

	if queryParams != nil {
		if queryParams.ForUpdateSkipLocked {
			q = q.Suffix("FOR UPDATE SKIP LOCKED")
		} else if queryParams.ForUpdate {
			q = q.Suffix("FOR UPDATE")
		} else if queryParams.ForShare {
			q = q.Suffix("FOR SHARE")
		}

		if queryParams.Limit > 0 {
			q = q.Limit(queryParams.Limit)
		}

		if queryParams.Offset > 0 {
			q = q.Offset(queryParams.Offset)
		}
	}

	query, args, err := q.ToSql()
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "building query", slog.Any("error", err))
		return nil, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	rows, err := r.pgClient.GetConn(ctx).Query(ctx, query, args...)
	if err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "query row error", slog.Any("error", err))
		}
		return nil, appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	defer rows.Close()

	dbData := []*pg.AccountCodeDBModel{}

	if err := pgxscan.ScanAll(&dbData, rows); err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "scan row error", slog.Any("error", err))
		}
		return nil, appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	result := make([]*entity.AccountCode, 0, len(dbData))
	for _, dbItem := range dbData {
		result = append(result, dbItem.ToEntity())
	}

	return result, nil
}

func (r *Repository) FindOneByID(
	ctx context.Context,
	id uuid.UUID,
	queryParams *uctypes.QueryGetOneParams,
) (*entity.AccountCode, error) {
	const op = "FindOneByID"

	where := squirrel.And{
		squirrel.Eq{"id": id},
	}

	q := r.qb.Select(pg.AccountCodeTableFields...).From(pg.AccountCodeTable).Where(where)

	if queryParams != nil {
		if queryParams.ForUpdateSkipLocked {
			q = q.Suffix("FOR UPDATE SKIP LOCKED")
		} else if queryParams.ForUpdate {
			q = q.Suffix("FOR UPDATE")
		} else if queryParams.ForShare {
			q = q.Suffix("FOR SHARE")
		}
	}

	query, args, err := q.ToSql()
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "building query", slog.Any("error", err))
		return nil, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	rows, err := r.pgClient.GetConn(ctx).Query(ctx, query, args...)
	if err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "query row error", slog.Any("error", err))
		}
		return nil, appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	defer rows.Close()

	dbData := &pg.AccountCodeDBModel{}

	if err := pgxscan.ScanOne(dbData, rows); err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "scan row error", slog.Any("error", err))
		}
		return nil, appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	return dbData.ToEntity(), nil
}
