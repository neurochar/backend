package profile

import (
	"time"

	"github.com/google/uuid"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

// DBModel - database model
type DBModel struct {
	ID                 uint64     `db:"id"`
	AccountID          uuid.UUID  `db:"account_id"`
	Name               string     `db:"name"`
	Surname            string     `db:"surname"`
	Photo100x100FileID *uuid.UUID `db:"photo_100x100_file_id"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (db *DBModel) ToEntity() *userEntity.Profile {
	return &userEntity.Profile{
		ID:                 db.ID,
		AccountID:          db.AccountID,
		Name:               db.Name,
		Surname:            db.Surname,
		Photo100x100FileID: db.Photo100x100FileID,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func mapEntityToDBModel(entity *userEntity.Profile) *DBModel {
	return &DBModel{
		ID:                 entity.ID,
		AccountID:          entity.AccountID,
		Name:               entity.Name,
		Surname:            entity.Surname,
		Photo100x100FileID: entity.Photo100x100FileID,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
