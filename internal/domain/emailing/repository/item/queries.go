package item

import (
	"context"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/emailing/entity"
	"github.com/neurochar/backend/internal/domain/emailing/usecase"
	"github.com/neurochar/backend/internal/infra/loghandler"
)

func (r *Repository) buildWhereForList(listOptions *usecase.ListOptions) (where squirrel.And) {
	if listOptions == nil {
		return where
	}

	if listOptions.FilterSentAtCompare != nil {
		col := "sent_at"
		val := listOptions.FilterSentAtCompare.Value

		switch listOptions.FilterSentAtCompare.Type {
		case uctypes.CompareEqual:
			if val == nil {
				where = append(where, squirrel.Eq{col: nil})
			} else {
				where = append(where, squirrel.Eq{col: val})
			}

		case uctypes.CompareNotEqual:
			if val == nil {
				where = append(where, squirrel.NotEq{col: nil})
			} else {
				where = append(where, squirrel.Or{
					squirrel.NotEq{col: val},
					squirrel.Eq{col: nil},
				})
			}

		case uctypes.CompareMore:
			where = append(where, squirrel.Gt{col: val})
		case uctypes.CompareLess:
			where = append(where, squirrel.Lt{col: val})
		case uctypes.CompareMoreOrEqual:
			where = append(where, squirrel.GtOrEq{col: val})
		case uctypes.CompareLessOrEqual:
			where = append(where, squirrel.LtOrEq{col: val})
		}
	}

	return where
}

func (r *Repository) buildSortForList(listOptions *usecase.ListOptions) []string {
	if listOptions == nil {
		return []string{}
	}

	return []string{}
}

func (r *Repository) FindList(
	ctx context.Context,
	listOptions *usecase.ListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*entity.Item, error) {
	const op = "FindList"

	where := r.buildWhereForList(listOptions)

	fields := TableFields

	q := r.qb.Select(fields...).From(Table).Where(where)

	sort := r.buildSortForList(listOptions)
	if len(sort) > 0 {
		q = q.OrderBy(sort...)
	}

	if queryParams != nil {
		if queryParams.ForUpdateSkipLocked {
			q = q.Suffix("FOR UPDATE SKIP LOCKED")
		} else if queryParams.ForUpdateNoWait {
			q = q.Suffix("FOR UPDATE NOWAIT")
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

	dbData := []*DBModel{}

	if err := pgxscan.ScanAll(&dbData, rows); err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "scan row error", slog.Any("error", err))
		}
		return nil, appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	result := make([]*entity.Item, 0, len(dbData))
	for _, dbItem := range dbData {
		result = append(result, dbItem.ToEntity())
	}

	return result, nil
}

func (r *Repository) FindOneByID(
	ctx context.Context,
	id uuid.UUID,
	queryParams *uctypes.QueryGetOneParams,
) (*entity.Item, error) {
	const op = "FindOneByID"

	where := squirrel.And{
		squirrel.Eq{"id": id},
	}

	q := r.qb.Select(TableFields...).From(Table).Where(where)

	if queryParams != nil {
		if queryParams.ForUpdateSkipLocked {
			q = q.Suffix("FOR UPDATE SKIP LOCKED")
		} else if queryParams.ForUpdateNoWait {
			q = q.Suffix("FOR UPDATE NOWAIT")
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

	dbData := &DBModel{}

	if err := pgxscan.ScanOne(dbData, rows); err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "scan row error", slog.Any("error", err))
		}
		return nil, appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	return dbData.ToEntity(), nil
}
