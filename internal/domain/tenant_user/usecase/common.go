package usecase

import (
	"context"

	"github.com/google/uuid"
)

type CommonUsecase interface {
	UpdatePasswordByRecoveryCode(
		ctx context.Context,
		codeID uuid.UUID,
		codeValue string,
		newPassword string,
		removeSessions bool,
	) (resErr error)
}
