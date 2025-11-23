package pg

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/domain/tenant/entity"
	"github.com/neurochar/backend/pkg/dbhelper"
)

const (
	RegistrationTable = "registration"
)

var RegistrationTableFields = []string{}

func init() {
	RegistrationTableFields = dbhelper.ExtractDBFields(&RegistrationDBModel{})
}

type RegistrationDBModel struct {
	ID         uuid.UUID  `db:"id"`
	Email      string     `db:"email"`
	Tariff     uint64     `db:"tariff"`
	IsFinished bool       `db:"is_finished"`
	TenantID   *uuid.UUID `db:"tenant_id"`
	RequestIP  *net.IP    `db:"request_ip"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (db *RegistrationDBModel) ToEntity() *entity.Registration {
	return &entity.Registration{
		ID:         db.ID,
		Email:      db.Email,
		Tariff:     db.Tariff,
		IsFinished: db.IsFinished,
		TenantID:   db.TenantID,
		RequestIP:  db.RequestIP,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
	}
}

func MapRegistrationEntityToDBModel(entity *entity.Registration) *RegistrationDBModel {
	return &RegistrationDBModel{
		ID:         entity.ID,
		Email:      entity.Email,
		Tariff:     entity.Tariff,
		IsFinished: entity.IsFinished,
		TenantID:   entity.TenantID,
		RequestIP:  entity.RequestIP,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
