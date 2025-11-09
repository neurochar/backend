package usecase

import (
	"context"
	"net"

	"github.com/google/uuid"

	"github.com/neurochar/backend/internal/common/uctypes"
)

type UserUsecase interface {
	FindPagedList(
		ctx context.Context,
		listOptions *UserListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*UserDTO, total uint64, resErr error)

	FindOneByProfileID(
		ctx context.Context,
		profileID uint64,
	) (resItem *UserDTO, resErr error)

	FindOneByAccountID(
		ctx context.Context,
		accountID uuid.UUID,
	) (resItem *UserDTO, resErr error)

	PatchAccountByDTO(
		ctx context.Context,
		id uuid.UUID,
		in PatchAccountDataInput,
		removeSessions bool,
		skipVersionCheck bool,
	) (resErr error)

	DeleteAccountActiveSessions(
		ctx context.Context,
		accountID uuid.UUID,
	) (resErr error)

	CreateUser(
		ctx context.Context,
		in CreateUserInput,
		requestIP net.IP,
		isAdminMode bool,
	) (resItem *UserDTO, resErr error)

	DeleteRole(ctx context.Context, roleID uint64) error

	UpdatePasswordByRecoveryCode(
		ctx context.Context,
		codeID uuid.UUID,
		codeValue string,
		newPassword string,
		removeSessions bool,
	) (resErr error)
}
