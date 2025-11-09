package entity

import (
	"time"
)

type RoleToRight struct {
	RoleID  uint64
	RightID uint64
	Value   int

	CreatedAt time.Time
}

func NewRoleToRight(roleID uint64, rightID uint64, value int) *RoleToRight {
	timeNow := time.Now()

	item := &RoleToRight{
		RoleID:    roleID,
		RightID:   rightID,
		Value:     value,
		CreatedAt: timeNow,
	}

	return item
}
