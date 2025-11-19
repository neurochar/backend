package usecase

import (
	"context"
	"net"

	"github.com/google/uuid"
)

type CommonUsecase interface {
	CreateUser(
		ctx context.Context,
		tenantID uuid.UUID,
		in CreateAccountDataInput,
		authorID uuid.UUID,
		requestIP net.IP,
	) (resAccountDTO *AccountDTO, resErr error)

	PatchAccountByDTO(
		ctx context.Context,
		id uuid.UUID,
		in PatchAccountDataInput,
		authorID uuid.UUID,
		skipVersionCheck bool,
	) (resErr error)

	UpdatePasswordByRecoveryCode(
		ctx context.Context,
		codeID uuid.UUID,
		codeValue string,
		newPassword string,
		removeSessions bool,
	) (resErr error)
}
