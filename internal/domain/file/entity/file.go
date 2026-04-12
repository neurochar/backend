package entity

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID                  uuid.UUID
	GroupID             uuid.UUID
	Target              string
	AssignedToTarget    bool
	StorageFileKey      *string
	OriginalFileName    string
	FileMimetype        *string
	FileHash            *string
	UploadedToStorage   bool
	ToDeleteFromStorage bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewFile(
	groupID uuid.UUID,
	target string,
	assignedToTarget bool,
	originalFileName string,
) *File {
	tNow := time.Now()

	return &File{
		ID:               uuid.New(),
		GroupID:          groupID,
		Target:           target,
		AssignedToTarget: assignedToTarget,
		OriginalFileName: originalFileName,
		CreatedAt:        tNow,
		UpdatedAt:        tNow,
	}
}

func (i *File) SetStorageFileKey(key string) {
	i.StorageFileKey = &key
}

func (i *File) SetUploadedToStorage(value bool) {
	i.UploadedToStorage = value
}

func (i *File) SetAssignedToTarget(value bool) {
	i.AssignedToTarget = value
}

func (i *File) SetFileMimetype(value *string) {
	i.FileMimetype = value
}

func (i *File) SetFileHash(value *string) {
	i.FileHash = value
}
