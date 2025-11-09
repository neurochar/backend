package roletoright

import (
	"context"
	"log/slog"

	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
	"github.com/neurochar/backend/internal/infra/loghandler"
	"github.com/neurochar/backend/pkg/dbhelper"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func (r *Repository) Create(ctx context.Context, item *userEntity.RoleToRight) error {
	const op = "Create"

	dataMap, err := dbhelper.DBModelToMap(mapEntityToDBModel(item))
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "convert struct to db map", slog.Any("error", err))
		return appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	query, args, err := r.qb.Insert(RoleToRightTable).SetMap(dataMap).ToSql()
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

func (r *Repository) DeleteByRoleID(ctx context.Context, roleID uint64) error {
	const op = "DeleteByRoleID"

	query, args, err := r.qb.Delete(RoleToRightTable).Where("role_id = ?", roleID).ToSql()
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
