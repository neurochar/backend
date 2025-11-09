package account

import (
	"net"
	"time"

	"github.com/google/uuid"

	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

// DBModel - database model
type DBModel struct {
	ID              uuid.UUID  `db:"id"`
	RoleID          uint64     `db:"role_id"`
	Email           string     `db:"email"`
	PasswordHash    string     `db:"password_hash"`
	IsEmailVerified bool       `db:"is_email_verified"`
	IsBlocked       bool       `db:"is_blocked"`
	LastLoginAt     *time.Time `db:"last_login_at"`
	LastRequestAt   *time.Time `db:"last_request_at"`
	LastRequestIP   *net.IP    `db:"last_request_ip"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (db *DBModel) ToEntity() *userEntity.Account {
	return &userEntity.Account{
		ID:              db.ID,
		RoleID:          db.RoleID,
		Email:           db.Email,
		PasswordHash:    db.PasswordHash,
		IsEmailVerified: db.IsEmailVerified,
		IsBlocked:       db.IsBlocked,
		LastLoginAt:     db.LastLoginAt,
		LastRequestAt:   db.LastRequestAt,
		LastRequestIP:   db.LastRequestIP,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func mapEntityToDBModel(entity *userEntity.Account) *DBModel {
	return &DBModel{
		ID:              entity.ID,
		RoleID:          entity.RoleID,
		Email:           entity.Email,
		PasswordHash:    entity.PasswordHash,
		IsEmailVerified: entity.IsEmailVerified,
		IsBlocked:       entity.IsBlocked,
		LastLoginAt:     entity.LastLoginAt,
		LastRequestAt:   entity.LastRequestAt,
		LastRequestIP:   entity.LastRequestIP,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
