package usecase

import (
	"context"

	"github.com/neurochar/backend/internal/common/uctypes"

	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

type RoleRepository interface {
	FindList(
		ctx context.Context,
		listOptions *RoleListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*userEntity.Role, err error)

	FindOneByID(
		ctx context.Context,
		id uint64,
		queryParams *uctypes.QueryGetOneParams,
	) (item *userEntity.Role, err error)

	Create(ctx context.Context, item *userEntity.Role) (err error)

	Update(ctx context.Context, item *userEntity.Role) (err error)
}

type RoleToRightRepository interface {
	FindList(
		ctx context.Context,
		listOptions *RoleToRightListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*userEntity.RoleToRight, err error)

	Create(ctx context.Context, item *userEntity.RoleToRight) (err error)

	DeleteByRoleID(ctx context.Context, roleID uint64) (err error)
}
