package usecase

import (
	"context"
)

type Usecase interface {
	SendAlert(ctx context.Context, message string) error
}
