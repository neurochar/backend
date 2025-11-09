package accountcode

import (
	"context"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
	"github.com/neurochar/backend/internal/infra/loghandler"
	"github.com/neurochar/backend/pkg/dbhelper"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (r *Repository) Create(ctx context.Context, item *userEntity.AccountCode) error {
	const op = "Create"

	dataMap, err := dbhelper.DBModelToMap(mapEntityToDBModel(item))
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "convert struct to db map", slog.Any("error", err))
		return appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	query, args, err := r.qb.Insert(AccountCodeTable).SetMap(dataMap).ToSql()
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

func (r *Repository) Update(ctx context.Context, item *userEntity.AccountCode) error {
	const op = "Update"

	currentUpdatedAt := item.UpdatedAt
	timeNow := time.Now()

	dataMap, err := dbhelper.DBModelToMap(mapEntityToDBModel(item))
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "convert struct to db map", slog.Any("error", err))
		return appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}
	delete(dataMap, "id")
	dataMap["updated_at"] = timeNow

	err = r.pgClient.Do(ctx, func(ctx context.Context) error {
		checkQuery, checkArgs, err := r.qb.Select("id").From(AccountCodeTable).Where(squirrel.Eq{"id": item.ID}).ToSql()
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

		updQuery, updArgs, err := r.qb.Update(AccountCodeTable).Where(
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
			return appErrors.Chainf(appErrors.ErrStoreOptimisticConflict, "%s.%s", r.pkg, op)
		}

		item.UpdatedAt = timeNow

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
