package usecase

import (
	"context"
	"net"
	"time"

	"github.com/neurochar/backend/internal/domain/emailing/entity"
	"github.com/neurochar/backend/internal/infra/emailing"
)

type Usecase interface {
	Create(ctx context.Context, data emailing.Message, requestIP net.IP) (item *entity.Item, err error)

	JobProcessItemsToSend(ctx context.Context) (anyJobDone bool, err error)

	JobProcessItemsToDelete(ctx context.Context, ttl time.Duration) (anyJobDone bool, err error)
}
