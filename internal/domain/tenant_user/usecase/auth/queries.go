package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"
	"github.com/neurochar/backend/internal/domain/tenant_user/entity"
)

func (uc *UsecaseImpl) FindSessionByID(
	ctx context.Context,
	id uuid.UUID,
	queryParams *uctypes.QueryGetOneParams,
) (*entity.Session, error) {
	const op = "FindSessionByID"

	item, err := uc.repo.FindOneByID(ctx, id, queryParams)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return item, nil
}

func (uc *UsecaseImpl) IsSessionRevoked(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	const op = "IsSessionRevoked"

	_, err := uc.repo.FindOneByID(ctx, id, nil)
	if err != nil {
		if errors.Is(err, appErrors.ErrNotFound) {
			return true, nil
		}
		return false, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return false, nil
}
