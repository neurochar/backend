package candidate_resume

import (
	"context"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/crm/entity"
	"github.com/neurochar/backend/internal/domain/crm/repository/pg"
	"github.com/neurochar/backend/internal/domain/crm/usecase"
	"github.com/neurochar/backend/internal/infra/loghandler"
	"github.com/samber/lo"
)

func (r *Repository) buildWhereForList(listOptions *usecase.CandidateResumeListOptions, withDeleted bool) (where squirrel.And) {
	defer func() {
		if !withDeleted {
			where = append(where, squirrel.Expr("deleted_at IS NULL"))
		}
	}()

	if listOptions == nil {
		return where
	}

	if listOptions.FilterIDs != nil {
		where = append(where, squirrel.Eq{"id": lo.Uniq(*listOptions.FilterIDs)})
	}

	if listOptions.FilterCandidatesIDs != nil {
		where = append(where, squirrel.Eq{"candidate_id": lo.Uniq(*listOptions.FilterCandidatesIDs)})
	}

	if listOptions.FilterStatus != nil {
		where = append(where, squirrel.Eq{"status": *listOptions.FilterStatus})
	}

	return where
}

func (r *Repository) buildSortForList(listOptions *usecase.CandidateResumeListOptions) []string {
	if listOptions == nil || len(listOptions.Sort) == 0 {
		return []string{"created_at DESC"}
	}

	sort := make([]string, 0, len(listOptions.Sort))

	for _, sortOption := range listOptions.Sort {
		switch sortOption.Field {
		case usecase.CandidateResumeListOptionsSortFieldCreatedAt:
			if sortOption.IsDesc {
				sort = append(sort, "created_at DESC")
			} else {
				sort = append(sort, "created_at ASC")
			}
		}
	}

	return sort
}

func (r *Repository) FindList(
	ctx context.Context,
	listOptions *usecase.CandidateResumeListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*entity.CandidateResume, error) {
	const op = "FindList"

	withDeleted := queryParams != nil && queryParams.WithDeleted

	where := r.buildWhereForList(listOptions, withDeleted)

	fields := pg.CandidateResumeTableFields

	q := r.qb.Select(fields...).From(pg.CandidateResumeTable).Where(where)

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

	dbData := []*pg.CandidateResumeDBModel{}

	if err := pgxscan.ScanAll(&dbData, rows); err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "scan row error", slog.Any("error", err))
		}
		return nil, appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	result := make([]*entity.CandidateResume, 0, len(dbData))
	for _, dbItem := range dbData {
		result = append(result, dbItem.ToEntity())
	}

	return result, nil
}

func (r *Repository) FindPagedList(
	ctx context.Context,
	listOptions *usecase.CandidateResumeListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*entity.CandidateResume, uint64, error) {
	const op = "FindPagedList"

	withDeleted := queryParams != nil && queryParams.WithDeleted

	where := r.buildWhereForList(listOptions, withDeleted)

	fields := pg.CandidateResumeTableFields

	q := r.qb.Select(fields...).From(pg.CandidateResumeTable).Where(where)

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
		return nil, 0, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	totalQ := r.qb.Select("COUNT(*) as total").From(pg.CandidateResumeTable).Where(where)
	totalQuery, totalArgs, err := totalQ.ToSql()
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "building query for total", slog.Any("error", err))
		return nil, 0, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	var total uint64
	var result []*entity.CandidateResume

	err = r.pgClient.Do(ctx, func(ctx context.Context) error {
		rows, err := r.pgClient.GetConn(ctx).Query(ctx, query, args...)
		if err != nil {
			convErr, ok := appErrors.ConvertPgxToAppErr(err)
			if !ok {
				r.logger.ErrorContext(loghandler.WithSource(ctx), "query row error", slog.Any("error", err))
			}
			return appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
		}

		defer rows.Close()

		dbData := []*pg.CandidateResumeDBModel{}

		if err := pgxscan.ScanAll(&dbData, rows); err != nil {
			convErr, ok := appErrors.ConvertPgxToAppErr(err)
			if !ok {
				r.logger.ErrorContext(loghandler.WithSource(ctx), "scan row error", slog.Any("error", err))
			}
			return appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
		}

		result = make([]*entity.CandidateResume, 0, len(dbData))
		for _, dbItem := range dbData {
			result = append(result, dbItem.ToEntity())
		}

		row := r.pgClient.GetConn(ctx).QueryRow(ctx, totalQuery, totalArgs...)
		if err := row.Scan(&total); err != nil {
			convErr, ok := appErrors.ConvertPgxToAppErr(err)
			if !ok {
				r.logger.ErrorContext(loghandler.WithSource(ctx), "scan total error", slog.Any("error", err))
			}
			return appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
		}

		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *Repository) FindOneByID(
	ctx context.Context,
	id uuid.UUID,
	queryParams *uctypes.QueryGetOneParams,
) (*entity.CandidateResume, error) {
	const op = "FindOneByID"

	withDeleted := queryParams != nil && queryParams.WithDeleted

	where := squirrel.And{
		squirrel.Eq{"id": id},
	}

	if !withDeleted {
		where = append(where, squirrel.Expr("deleted_at IS NULL"))
	}

	q := r.qb.Select(pg.CandidateResumeTableFields...).From(pg.CandidateResumeTable).Where(where)

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

	dbData := &pg.CandidateResumeDBModel{}

	if err := pgxscan.ScanOne(dbData, rows); err != nil {
		convErr, ok := appErrors.ConvertPgxToAppErr(err)
		if !ok {
			r.logger.ErrorContext(loghandler.WithSource(ctx), "scan row error", slog.Any("error", err))
		}
		return nil, appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
	}

	return dbData.ToEntity(), nil
}
