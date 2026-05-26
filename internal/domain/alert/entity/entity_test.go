package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAlert(t *testing.T) {
	t.Parallel()

	alert, err := New("something went wrong")
	require.NoError(t, err)
	require.NotNil(t, alert)

	assert.Equal(t, "something went wrong", alert.Message)
	assert.False(t, alert.CreatedAt.IsZero())
	assert.True(t, alert.CreatedAt.Before(time.Now()))
	assert.True(t, alert.CreatedAt.Equal(alert.CreatedAt.Truncate(time.Microsecond)))
}
