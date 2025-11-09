package entity

import (
	"time"
)

type Alert struct {
	Message string

	CreatedAt time.Time
}

func New(
	message string,
) (*Alert, error) {
	tNow := time.Now().Truncate(time.Microsecond)

	item := &Alert{
		Message:   message,
		CreatedAt: tNow,
	}

	return item, nil
}
