package providing

import "github.com/neurochar/backend/pkg/backoff"

// NewBackoff - create new backoff
func NewBackoff() *backoff.Controller {
	return backoff.NewController()
}
