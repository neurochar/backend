// Package session contains session repository
package adminsession

import (
	"net"
	"time"

	"github.com/google/uuid"

	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

// DBModel - database model
type DBModel struct {
	ID            uuid.UUID `db:"id"`
	AccountID     uuid.UUID `db:"account_id"`
	LastRequestAt time.Time `db:"last_request_at"`
	LastRequestIP net.IP    `db:"last_request_ip"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (db *DBModel) ToEntity() *userEntity.AdminSession {
	return &userEntity.AdminSession{
		ID:            db.ID,
		AccountID:     db.AccountID,
		LastRequestAt: db.LastRequestAt,
		LastRequestIP: db.LastRequestIP,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func mapEntityToDBModel(entity *userEntity.AdminSession) *DBModel {
	return &DBModel{
		ID:            entity.ID,
		AccountID:     entity.AccountID,
		LastRequestAt: entity.LastRequestAt,
		LastRequestIP: entity.LastRequestIP,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
