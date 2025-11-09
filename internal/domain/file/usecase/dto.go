package usecase

import (
	"time"

	"github.com/google/uuid"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/common/uctypes"

	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
)

type ListOptions struct {
	IDs                 *[]uuid.UUID
	AssignedToTarget    *bool
	StorageFileKey      **string
	ToDeleteFromStorage *bool
	CreatedAt           *ListOptionsCreatedAt
}

type ListOptionsCreatedAt struct {
	Value   *time.Time
	Compare uctypes.CompareType
}

type ProcessFileToTargetIn struct {
	CurrentFileID *uuid.UUID
	NewFileID     *uuid.UUID
	Target        string
	Group         string
}

type UploadFilesIn []struct {
	Target   string
	FileName string
	FileData []byte
	Process  func([]byte) ([]byte, *appErrors.AppError)
}

type UploadFilesOut map[uuid.UUID]*fileEntity.File
