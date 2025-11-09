package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

type AccountRepository interface {
	FindOneByEmail(
		ctx context.Context,
		email string,
		queryParams *uctypes.QueryGetOneParams,
	) (account *userEntity.Account, err error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (account *userEntity.Account, err error)

	FindList(
		ctx context.Context,
		listOptions *AccountListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*userEntity.Account, err error)

	Create(ctx context.Context, item *userEntity.Account) (err error)

	Update(ctx context.Context, item *userEntity.Account) (err error)
}

type AccountCodeRepository interface {
	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (item *userEntity.AccountCode, err error)

	FindList(
		ctx context.Context,
		listOptions *AccountCodeListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*userEntity.AccountCode, err error)

	Create(ctx context.Context, item *userEntity.AccountCode) (err error)

	Update(ctx context.Context, item *userEntity.AccountCode) (err error)
}
