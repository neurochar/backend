package entity

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/infra/emailing"
)

type Item struct {
	ID          uuid.UUID
	MessageData emailing.Message
	RequestIP   net.IP
	SentAt      *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (item *Item) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

func New(
	messageData emailing.Message,
	requestIP net.IP,
) (*Item, error) {
	tNow := time.Now().Truncate(time.Microsecond)

	item := &Item{
		ID:          uuid.New(),
		MessageData: messageData,
		RequestIP:   requestIP,
		CreatedAt:   tNow,
		UpdatedAt:   tNow,
	}

	return item, nil
}
