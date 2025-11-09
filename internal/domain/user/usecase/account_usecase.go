package usecase

import (
	"context"
	"net"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/common/uctypes"

	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

type AccountUsecase interface {
	FindOneByEmail(
		ctx context.Context,
		email string,
		queryParams *uctypes.QueryGetOneParams,
	) (resAccount *userEntity.Account, resErr error)

	FindOneByID(
		ctx context.Context,
		id uuid.UUID,
		queryParams *uctypes.QueryGetOneParams,
	) (resAccount *userEntity.Account, resErr error)

	FindList(
		ctx context.Context,
		listOptions *AccountListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems []*userEntity.Account, resErr error)

	FindListInMap(
		ctx context.Context,
		listOptions *AccountListOptions,
		queryParams *uctypes.QueryGetListParams,
	) (resItems map[uuid.UUID]*userEntity.Account, resErr error)

	CreateAccountByDTO(
		ctx context.Context,
		in AccountDataInput,
		requestIP net.IP,
	) (resAccount *userEntity.Account, activationCode *userEntity.AccountCode, resErr error)

	PatchAccountByDTO(
		ctx context.Context,
		id uuid.UUID,
		in PatchAccountDataInput,
		skipVersionCheck bool,
	) (resErr error)

	UpdateAccount(ctx context.Context, item *userEntity.Account) (resErr error)

	VerifyAccountEmailByCode(
		ctx context.Context,
		codeID uuid.UUID,
		codeValue string,
	) (resErr error)

	RequestPasswordRecoveryByEmail(
		ctx context.Context,
		email string,
		requestIP net.IP,
	) (resCode *userEntity.AccountCode, resErr error)

	UpdatePasswordByRecoveryCode(
		ctx context.Context,
		codeID uuid.UUID,
		codeValue string,
		newPassword string,
	) (resCode *userEntity.AccountCode, resErr error)

	CheckCode(
		ctx context.Context,
		codeID uuid.UUID,
		codeValue string,
	) (resCode *userEntity.AccountCode, resErr error)

	UpdateCode(
		ctx context.Context,
		code *userEntity.AccountCode,
	) (resErr error)
}
