package usecase

import (
	"context"

	"github.com/google/uuid"
)

type CrossUsecase interface {
	DeleteProfile(ctx context.Context, id uuid.UUID) (resErr error)

	DeleteRoom(ctx context.Context, id uuid.UUID) (resErr error)
}
