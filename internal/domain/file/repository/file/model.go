package file

import (
	"time"

	"github.com/google/uuid"
	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
)

// DBModel - database model
type DBModel struct {
	ID                  uuid.UUID `db:"id"`
	GroupID             uuid.UUID `db:"group_id"`
	Target              string    `db:"file_target"`
	AssignedToTarget    bool      `db:"assigned_to_target"`
	StorageFileKey      *string   `db:"storage_file_key"`
	OriginalFileName    string    `db:"original_file_name"`
	UploadedToStorage   bool      `db:"uploaded_to_storage"`
	ToDeleteFromStorage bool      `db:"to_delete_from_storage"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (db *DBModel) ToEntity() *fileEntity.File {
	return &fileEntity.File{
		ID:                  db.ID,
		GroupID:             db.GroupID,
		Target:              db.Target,
		AssignedToTarget:    db.AssignedToTarget,
		StorageFileKey:      db.StorageFileKey,
		OriginalFileName:    db.OriginalFileName,
		UploadedToStorage:   db.UploadedToStorage,
		ToDeleteFromStorage: db.ToDeleteFromStorage,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: db.DeletedAt,
	}
}

func mapEntityToDBModel(entity *fileEntity.File) *DBModel {
	return &DBModel{
		ID:                  entity.ID,
		GroupID:             entity.GroupID,
		Target:              entity.Target,
		AssignedToTarget:    entity.AssignedToTarget,
		StorageFileKey:      entity.StorageFileKey,
		OriginalFileName:    entity.OriginalFileName,
		UploadedToStorage:   entity.UploadedToStorage,
		ToDeleteFromStorage: entity.ToDeleteFromStorage,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
