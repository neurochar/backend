package usecase

import (
	"context"

	appErrors "github.com/neurochar/backend/internal/app/errors"

	"github.com/neurochar/backend/internal/domain/alert/entity"
)

func (uc *UsecaseImpl) SendAlert(ctx context.Context, message string) error {
	const op = "SendAlert"

	alert, err := entity.New(message)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.repo.SendMessage(ctx, alert)
	if err != nil {
		return appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return nil
}
