package entity

import (
	"time"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

var ErrRoleNameEmpty = appErrors.ErrBadRequest.Extend("role name is empty")

// Role - role entity
type Role struct {
	ID       uint64
	Name     string
	IsSystem bool
	IsSuper  bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Version - get version
func (item *Role) Version() int64 {
	return item.UpdatedAt.UnixMicro()
}

// NewRole - constructor for new role
func NewRole(name string) (*Role, error) {
	if name == "" {
		return nil, ErrRoleNameEmpty
	}

	timeNow := time.Now().Truncate(time.Microsecond)

	role := &Role{
		Name:      name,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}

	return role, nil
}
