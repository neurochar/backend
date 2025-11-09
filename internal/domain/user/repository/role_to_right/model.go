package roletoright

import (
	"time"

	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

// DBModel - database model
type DBModel struct {
	RoleID      uint64 `db:"role_id"`
	RoleRightID uint64 `db:"role_right_id"`
	Value       int    `db:"value"`

	CreatedAt time.Time `db:"created_at"`
}

func (db *DBModel) ToEntity() *userEntity.RoleToRight {
	return &userEntity.RoleToRight{
		RoleID:  db.RoleID,
		RightID: db.RoleRightID,
		Value:   db.Value,

		CreatedAt: db.CreatedAt,
	}
}

func mapEntityToDBModel(entity *userEntity.RoleToRight) *DBModel {
	return &DBModel{
		RoleID:      entity.RoleID,
		RoleRightID: entity.RightID,
		Value:       entity.Value,

		CreatedAt: entity.CreatedAt,
	}
}
