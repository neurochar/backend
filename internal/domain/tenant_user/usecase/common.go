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
		author *AccountDTO,
		requestIP net.IP,
	) (resAccountDTO *AccountDTO, resErr error)

	UpdatePasswordByRecoveryCode(
		ctx context.Context,
		codeID uuid.UUID,
		codeValue string,
		newPassword string,
		removeSessions bool,
	) (resErr error)
}
