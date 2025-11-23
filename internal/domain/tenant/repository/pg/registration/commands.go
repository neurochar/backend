package registration

import (
	"context"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/domain/tenant/entity"
	"github.com/neurochar/backend/internal/domain/tenant/repository/pg"
	"github.com/neurochar/backend/internal/domain/tenant/usecase"
	"github.com/neurochar/backend/internal/infra/loghandler"
	"github.com/neurochar/backend/pkg/dbhelper"
)

func (r *Repository) Create(ctx context.Context, item *entity.Registration) error {
	const op = "Create"

	dataMap, err := dbhelper.DBModelToMap(pg.MapRegistrationEntityToDBModel(item))
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "convert struct to db map", slog.Any("error", err))
		return appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	query, args, err := r.qb.Insert(pg.RegistrationTable).SetMap(dataMap).ToSql()
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "building query", slog.Any("error", err))
		return appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	_, err = r.pgClient.GetConn(ctx).Exec(ctx, query, args...)
	if err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "executing query", slog.Any("error", err))
		}
		return appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, item *entity.Registration) error {
	const op = "Update"

	currentUpdatedAt := item.UpdatedAt
	timeNow := time.Now().Truncate(time.Microsecond)

	dataMap, err := dbhelper.DBModelToMap(pg.MapRegistrationEntityToDBModel(item))
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "convert struct to db map", slog.Any("error", err))
		return appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}
	delete(dataMap, "id")
	dataMap["updated_at"] = timeNow

	err = r.pgClient.Do(ctx, func(ctx context.Context) error {
		checkQuery, checkArgs, err := r.qb.Select("id").From(pg.RegistrationTable).Where(squirrel.Eq{"id": item.ID}).ToSql()
		if err != nil {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "building check query", slog.Any("error", err))
			return appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
		}

		var checkID uuid.UUID
		err = r.pgClient.GetConn(ctx).QueryRow(ctx, checkQuery, checkArgs...).Scan(&checkID)
		if err != nil {
			convErr, ok := appErrors.ConvertPgxToAppErr(err)
			if !ok {
				r.logger.ErrorContext(loghandler.WithSource(ctx), "executing check query", slog.Any("error", err))
			}
			return appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
		}

		updQuery, updArgs, err := r.qb.Update(pg.RegistrationTable).Where(
			squirrel.Eq{"id": item.ID, "updated_at": currentUpdatedAt}).SetMap(dataMap).ToSql()
		if err != nil {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "building query", slog.Any("error", err))
			return appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
		}

		cmdTag, err := r.pgClient.GetConn(ctx).Exec(ctx, updQuery, updArgs...)
		if err != nil {
			convErr, ok := appErrors.ConvertPgxToAppErr(err)
			if !ok {
				r.logger.ErrorContext(loghandler.WithSource(ctx), "executing query", slog.Any("error", err))
			}
			return appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
		}

		if cmdTag.RowsAffected() == 0 {
			return appErrors.Chainf(appErrors.ErrConflict, "%s.%s", r.pkg, op)
		}

		item.UpdatedAt = timeNow

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	const op = "DeleteByID"

	query, args, err := r.qb.Delete(pg.RegistrationTable).Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "building query", slog.Any("error", err))
		return appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	ct, err := r.pgClient.GetConn(ctx).Exec(ctx, query, args...)
	if err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "executing query", slog.Any("error", err))
		}
		return appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	if ct.RowsAffected() == 0 {
		return appErrors.Chainf(appErrors.ErrBadRequest, "%s.%s", r.pkg, op)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, listOptions *usecase.RegistrationListOptions) (uint64, error) {
	const op = "Delete"

	where := r.buildWhereForList(listOptions)

	query, args, err := r.qb.Delete(pg.RegistrationTable).Where(where).ToSql()
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "building query", slog.Any("error", err))
		return 0, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	ct, err := r.pgClient.GetConn(ctx).Exec(ctx, query, args...)
	if err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "executing query", slog.Any("error", err))
		}
		return 0, appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	return uint64(ct.RowsAffected()), nil
}
