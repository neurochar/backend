package pg

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/domain/tenant/entity"
	"github.com/neurochar/backend/pkg/dbhelper"
)

const (
	AccountTable = "tenant_account"
)

var AccountTableFields = []string{}

func init() {
	AccountTableFields = dbhelper.ExtractDBFields(&AccountDBModel{})
}

type AccountDBModel struct {
	ID                         uuid.UUID  `db:"id"`
	TenantID                   uuid.UUID  `db:"tenant_id"`
	RoleID                     uint64     `db:"role_id"`
	Email                      string     `db:"email"`
	PasswordHash               string     `db:"password_hash"`
	IsConfirmed                bool       `db:"is_confirmed"`
	IsEmailVerified            bool       `db:"is_email_verified"`
	IsBlocked                  bool       `db:"is_blocked"`
	LastLoginAt                *time.Time `db:"last_login_at"`
	LastRequestAt              *time.Time `db:"last_request_at"`
	LastRequestIP              *net.IP    `db:"last_request_ip"`
	ProfileName                string     `db:"profile_name"`
	ProfileSurname             string     `db:"profile_surname"`
	ProfilePhoto100x100FileID  *uuid.UUID `db:"profile_photo_100x100_file_id"`
	ProfilePhotoOriginalFileID *uuid.UUID `db:"profile_photo_original_file_id"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (db *AccountDBModel) ToEntity() *entity.Account {
	return &entity.Account{
		ID:                         db.ID,
		TenantID:                   db.TenantID,
		RoleID:                     db.RoleID,
		Email:                      db.Email,
		PasswordHash:               db.PasswordHash,
		IsConfirmed:                db.IsConfirmed,
		IsEmailVerified:            db.IsEmailVerified,
		IsBlocked:                  db.IsBlocked,
		LastLoginAt:                db.LastLoginAt,
		LastRequestAt:              db.LastRequestAt,
		LastRequestIP:              db.LastRequestIP,
		ProfileName:                db.ProfileName,
		ProfileSurname:             db.ProfileSurname,
		ProfilePhoto100x100FileID:  db.ProfilePhoto100x100FileID,
		ProfilePhotoOriginalFileID: db.ProfilePhotoOriginalFileID,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func MapAccountEntityToDBModel(entity *entity.Account) *AccountDBModel {
	return &AccountDBModel{
		ID:                         entity.ID,
		TenantID:                   entity.TenantID,
		RoleID:                     entity.RoleID,
		Email:                      entity.Email,
		PasswordHash:               entity.PasswordHash,
		IsConfirmed:                entity.IsConfirmed,
		IsEmailVerified:            entity.IsEmailVerified,
		IsBlocked:                  entity.IsBlocked,
		LastLoginAt:                entity.LastLoginAt,
		LastRequestAt:              entity.LastRequestAt,
		LastRequestIP:              entity.LastRequestIP,
		ProfileName:                entity.ProfileName,
		ProfileSurname:             entity.ProfileSurname,
		ProfilePhoto100x100FileID:  entity.ProfilePhoto100x100FileID,
		ProfilePhotoOriginalFileID: entity.ProfilePhotoOriginalFileID,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
