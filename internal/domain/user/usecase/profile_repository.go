package usecase

import (
	"context"

	"github.com/neurochar/backend/internal/common/uctypes"

	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

type ProfileRepository interface {
	FindList(
		ctx context.Context,
		listOptions *ProfileListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*userEntity.Profile, err error)

	FindOneByID(
		ctx context.Context,
		id uint64,
		queryParams *uctypes.QueryGetOneParams,
	) (item *userEntity.Profile, err error)

	Create(ctx context.Context, item *userEntity.Profile) (err error)

	Update(ctx context.Context, item *userEntity.Profile) (err error)
}
