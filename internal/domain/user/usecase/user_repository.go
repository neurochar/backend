package usecase

import (
	"context"

	"github.com/neurochar/backend/internal/common/uctypes"
)

type ProfileAccountRepository interface {
	FindPagedList(
		ctx context.Context,
		listOptions *UserListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (items []*User, total uint64, err error)
}
