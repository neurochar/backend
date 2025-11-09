package profileaccount

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Masterminds/squirrel"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	accountRepo "github.com/neurochar/backend/internal/domain/user/repository/account"
	profileRepo "github.com/neurochar/backend/internal/domain/user/repository/profile"
	"github.com/neurochar/backend/internal/domain/user/usecase"
	"github.com/neurochar/backend/internal/infra/loghandler"
	"github.com/neurochar/backend/pkg/escape"
)

func (r *Repository) buildWhereForList(listOptions *usecase.UserListOptions, withDeleted bool) (where squirrel.And) {
	defer func() {
		if !withDeleted {
			where = append(where, squirrel.Expr("p.deleted_at IS NULL"))
		}
	}()

	if listOptions == nil {
		return where
	}

	if listOptions.RoleID != nil {
		where = append(where, squirrel.Eq{"a.role_id": *listOptions.RoleID})
	}

	if listOptions.Query != nil && strings.TrimSpace(*listOptions.Query) != "" {
		q := strings.TrimSpace(*listOptions.Query)
		qEscaped := escape.EscapeLikePattern(q)

		emailPred := squirrel.ILike{"a.email": qEscaped + "%"}

		parts := strings.Fields(q)
		var termClauses []squirrel.Sqlizer
		for _, term := range parts {
			escapedTerm := escape.EscapeLikePattern(term)

			termClauses = append(termClauses, squirrel.Or{
				squirrel.ILike{"p.name": escapedTerm + "%"},
				squirrel.ILike{"p.surname": escapedTerm + "%"},
				squirrel.ILike{"p.surname": "%-" + escapedTerm + "%"},
			})
		}

		namePred := squirrel.And(termClauses)

		where = append(where, squirrel.Or{emailPred, namePred})
	}

	return where
}

func (r *Repository) buildSortForList(_ *usecase.UserListOptions) (sort []string) {
	sort = append(sort, "p.id ASC")

	return sort
}

func (r *Repository) FindPagedList(
	ctx context.Context,
	listOptions *usecase.UserListOptions,
	queryParams *uctypes.QueryGetListParams,
) ([]*usecase.User, uint64, error) {
	const op = "FindPagedList"

	withDeleted := queryParams != nil && queryParams.WithDeleted

	where := r.buildWhereForList(listOptions, withDeleted)

	q := r.qb.Select(JoinFields...).
		From(fmt.Sprintf("%s p", profileRepo.ProfileTable)).
		Join(fmt.Sprintf("%s a ON p.account_id = a.id", accountRepo.AccountTable)).
		Where(where).
		OrderBy(r.buildSortForList(listOptions)...)

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
		return nil, 0, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	totalQ := r.qb.Select("COUNT(p.id) as total").
		From(fmt.Sprintf("%s p", profileRepo.ProfileTable)).
		Join(fmt.Sprintf("%s a ON p.account_id = a.id", accountRepo.AccountTable)).
		Where(where)

	totalQuery, totalArgs, err := totalQ.ToSql()
	if err != nil {
		r.logger.ErrorContext(loghandler.WithSource(ctx), "building query for total", slog.Any("error", err))
		return nil, 0, appErrors.Chainf(appErrors.ErrInternal.WithWrap(err), "%s.%s", r.pkg, op)
	}

	var total uint64
	var result []*usecase.User

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

		dbData := []*DBModel{}

		if err := r.pgScanApi.ScanAll(&dbData, rows); err != nil {
			convErr, ok := appErrors.ConvertPgxToAppErr(err)
			if !ok {
				r.logger.ErrorContext(loghandler.WithSource(ctx), "scan row error", slog.Any("error", err))
			}
			return appErrors.Chainf(convErr, "%s.%s", r.pkg, op)
		}

		result = make([]*usecase.User, 0, len(dbData))
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
