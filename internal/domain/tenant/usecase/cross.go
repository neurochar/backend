package usecase

import (
	"context"

	"github.com/google/uuid"
)

type CrossUsecase interface {
	PatchTenantByDTO(
		ctx context.Context,
		id uuid.UUID,
		in PatchTenantDataInput,
		skipVersionCheck bool,
	) (resErr error)
}
