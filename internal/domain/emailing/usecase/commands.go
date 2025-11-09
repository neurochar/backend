package usecase

import (
	"context"
	"net"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/infra/emailing"

	"github.com/neurochar/backend/internal/domain/emailing/entity"
)

func (uc *UsecaseImpl) Create(ctx context.Context, data emailing.Message, requestIP net.IP) (*entity.Item, error) {
	const op = "Create"

	item, err := entity.New(data, requestIP)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	err = uc.repo.Create(ctx, item)
	if err != nil {
		return nil, appErrors.Chainf(err, "%s.%s", uc.pkg, op)
	}

	return item, nil
}
