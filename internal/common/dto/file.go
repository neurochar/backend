package dto

import (
	"github.com/google/uuid"
	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
)

type FileDTO struct {
	ID       uuid.UUID `json:"id,omitzero"`
	URL      string    `json:"url"`
	Filename string    `json:"filename,omitempty"`
}

func NewFileDTO(file *fileEntity.File, fileUC fileUC.Usecase, isFull bool) *FileDTO {
	if file == nil {
		return nil
	}

	if !isFull {
		return &FileDTO{
			URL: fileUC.FileUrl(file),
		}
	}

	return &FileDTO{
		ID:       file.ID,
		URL:      fileUC.FileUrl(file),
		Filename: file.OriginalFileName,
	}
}

type UploadedFilePackDTO struct {
	Files map[string]*FileDTO `json:"files"`
}
