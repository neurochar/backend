package usecase

import (
	"github.com/google/uuid"

	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

type FullProfileDTO struct {
	Profile          *userEntity.Profile
	Photo100x100File *fileEntity.File
}

type ProfileListOptions struct {
	AccountID *uuid.UUID
}

type ProfileFileTarget string

const (
	FileTargetProfilePhoto100x100 ProfileFileTarget = "profile:photo:100x100"
)

type UploadFileOut []*fileEntity.File

type ProfileDataInput struct {
	Version int64

	Name               string
	Surname            string
	Photo100x100FileID *uuid.UUID
}
