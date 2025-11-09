package pg

import (
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/domain/tenant/entity"
	"github.com/neurochar/backend/pkg/dbhelper"
)

const (
	TenantTable = "tenant"
)

var TenantTableFields = []string{}

func init() {
	TenantTableFields = dbhelper.ExtractDBFields(&TenantDBModel{})
}

type TenantDBModel struct {
	ID       uuid.UUID `db:"id"`
	TextID   string    `db:"text_id"`
	IsDemo   bool      `db:"is_demo"`
	IsActive bool      `db:"is_active"`
	Name     string    `db:"name"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (db *TenantDBModel) ToEntity() *entity.Tenant {
	return &entity.Tenant{
		ID:       db.ID,
		TextID:   db.TextID,
		IsDemo:   db.IsDemo,
		IsActive: db.IsActive,
		Name:     db.Name,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func MapTenantEntityToDBModel(entity *entity.Tenant) *TenantDBModel {
	return &TenantDBModel{
		ID:       entity.ID,
		TextID:   entity.TextID,
		IsDemo:   entity.IsDemo,
		IsActive: entity.IsActive,
		Name:     entity.Name,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
