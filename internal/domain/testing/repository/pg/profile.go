package pg

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/domain/testing/entity"
	"github.com/neurochar/backend/pkg/dbhelper"
)

const (
	ProfileTable = "testing_profile"
)

var ProfileTableFields = []string{}

func init() {
	ProfileTableFields = dbhelper.ExtractDBFields(&ProfileDBModel{})
}

type ProfileDBModel struct {
	ID                   uuid.UUID       `db:"id"`
	TenantID             uuid.UUID       `db:"tenant_id"`
	Name                 string          `db:"name"`
	Description          string          `db:"description"`
	PersonalityTraitsMap json.RawMessage `db:"personality_traits_map"`
	CreatedBy            *uuid.UUID      `db:"created_by"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (db *ProfileDBModel) ToEntity() *entity.Profile {
	traitsMap := make(entity.ProfilePersonalityTraitsMap)
	err := json.Unmarshal(db.PersonalityTraitsMap, &traitsMap)
	if err != nil {
		panic(err)
	}

	return &entity.Profile{
		ID:                   db.ID,
		TenantID:             db.TenantID,
		Name:                 db.Name,
		Description:          db.Description,
		PersonalityTraitsMap: traitsMap,
		CreatedBy:            db.CreatedBy,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func MapProfileEntityToDBModel(entity *entity.Profile) *ProfileDBModel {
	traitsMap, err := json.Marshal(entity.PersonalityTraitsMap)
	if err != nil {
		panic(err)
	}

	return &ProfileDBModel{
		ID:                   entity.ID,
		TenantID:             entity.TenantID,
		Name:                 entity.Name,
		Description:          entity.Description,
		PersonalityTraitsMap: traitsMap,
		CreatedBy:            entity.CreatedBy,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
