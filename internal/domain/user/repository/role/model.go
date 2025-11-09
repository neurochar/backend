package role

import (
	"time"

	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

// DBModel - database model
type DBModel struct {
	ID       uint64 `db:"id"`
	Name     string `db:"name"`
	IsSystem bool   `db:"is_system"`
	IsSuper  bool   `db:"is_super"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (db *DBModel) ToEntity() *userEntity.Role {
	return &userEntity.Role{
		ID:       db.ID,
		Name:     db.Name,
		IsSystem: db.IsSystem,
		IsSuper:  db.IsSuper,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func mapEntityToDBModel(entity *userEntity.Role) *DBModel {
	return &DBModel{
		ID:       entity.ID,
		Name:     entity.Name,
		IsSystem: entity.IsSystem,
		IsSuper:  entity.IsSuper,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
