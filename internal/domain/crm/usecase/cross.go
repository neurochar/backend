package usecase

import (
	"context"

	"github.com/google/uuid"
)

type CrossUsecase interface {
	Delete(ctx context.Context, id uuid.UUID) (resErr error)
}
