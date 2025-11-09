package entity

import (
	"net"
	"time"

	"github.com/google/uuid"
)

// AdminSession - session entity
type AdminSession struct {
	ID            uuid.UUID
	AccountID     uuid.UUID
	LastRequestAt time.Time
	LastRequestIP net.IP

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Version - get version
func (item *AdminSession) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

// IsAlive - check if session is alive
func (item *AdminSession) IsAlive(duration time.Duration) bool {
	return time.Since(item.LastRequestAt) < duration
}

func (item *AdminSession) SetLastRequestAt(value time.Time) {
	item.LastRequestAt = value.Truncate(time.Microsecond)
}

func (item *AdminSession) SetLastRequestIP(value net.IP) {
	item.LastRequestIP = value
}

// NewSession - constructor for new session
func NewSession(accountID uuid.UUID, ip net.IP) *AdminSession {
	timeNow := time.Now().Truncate(time.Microsecond)

	item := &AdminSession{
		ID:        uuid.New(),
		AccountID: accountID,

		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	item.SetLastRequestAt(timeNow)
	item.SetLastRequestIP(ip)

	return item
}
