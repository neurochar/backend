package pg

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/domain/tenant_user/entity"
	"github.com/neurochar/backend/pkg/dbhelper"
)

const (
	AccountCodeTable = "tenant_account_code"
)

var AccountCodeTableFields = []string{}

func init() {
	AccountCodeTableFields = dbhelper.ExtractDBFields(&AccountCodeDBModel{})
}

type AccountCodeDBModel struct {
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

func (db *AccountCodeDBModel) ToEntity() *entity.AccountCode {
	return &entity.AccountCode{
		ID:        db.ID,
		AccountID: db.AccountID,
		Type:      entity.AccountCodeType(db.Type),
		IsActive:  db.IsActive,
		Code:      db.Code,
		RequestIP: db.RequestIP,
		Attempts:  db.Attempts,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
	}
}

func MapAccountCodeEntityToDBModel(entity *entity.AccountCode) *AccountCodeDBModel {
	return &AccountCodeDBModel{
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
