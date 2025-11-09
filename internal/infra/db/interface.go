// Package db contains db interfaces and helpers
package db

import "github.com/neurochar/backend/pkg/pgclient"

// MasterClient - interface for master (read + write)
type MasterClient interface {
	pgclient.Client
}
