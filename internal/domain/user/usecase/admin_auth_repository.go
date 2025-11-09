package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

type SessionRepository interface {
	Create(
		ctx context.Context,
		item *userEntity.AdminSession,
	) (err error)

	Update(
		ctx context.Context,
		item *userEntity.AdminSession,
	) (err error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (session *userEntity.AdminSession, err error)

	FindList(
		ctx context.Context,
		listOptions *AdminAuthListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*userEntity.AdminSession, err error)
}
