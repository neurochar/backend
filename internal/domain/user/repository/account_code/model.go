package accountcode

import (
	"net"
	"time"

	"github.com/google/uuid"

	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

// DBModel - database model
type DBModel struct {
	ID        uuid.UUID `db:"id"`
	AccountID uuid.UUID `db:"account_id"`
	Type      uint8     `db:"code_type"`
	IsActive  bool      `db:"is_active"`
	Code      string    `db:"code"`
	RequestIP net.IP    `db:"request_ip"`
	Attempts  int       `db:"attempts"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (db *DBModel) ToEntity() *userEntity.AccountCode {
	return &userEntity.AccountCode{
		ID:        db.ID,
		AccountID: db.AccountID,
		Type:      userEntity.AccountCodeType(db.Type),
		IsActive:  db.IsActive,
		Code:      db.Code,
		RequestIP: db.RequestIP,
		Attempts:  db.Attempts,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
	}
}

func mapEntityToDBModel(entity *userEntity.AccountCode) *DBModel {
	return &DBModel{
		ID:        entity.ID,
		AccountID: entity.AccountID,
		Type:      uint8(entity.Type),
		IsActive:  entity.IsActive,
		Code:      entity.Code,
		RequestIP: entity.RequestIP,
		Attempts:  entity.Attempts,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
