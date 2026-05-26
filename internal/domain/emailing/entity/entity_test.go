package entity

import (
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/infra/emailing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEmailingItem(t *testing.T) {
	t.Parallel()

	msg := emailing.Message{
		To:        "user@example.com",
		Subject:   "Test",
		TextPlain: "Hello",
	}
	ip := net.ParseIP("10.0.0.1")

	item, err := New(msg, ip)
	require.NoError(t, err)
	require.NotNil(t, item)

	assert.NotEqual(t, uuid.Nil, item.ID)
	assert.Equal(t, msg, item.MessageData)
	assert.Equal(t, ip, item.RequestIP)
	assert.Nil(t, item.SentAt)
	assert.False(t, item.CreatedAt.IsZero())
	assert.False(t, item.UpdatedAt.IsZero())
}

func TestEmailingItem_Version(t *testing.T) {
	t.Parallel()

	now := time.Now().Truncate(time.Microsecond)
	item := &Item{UpdatedAt: now}
	assert.Equal(t, now.UnixMicro(), item.Version())
}
