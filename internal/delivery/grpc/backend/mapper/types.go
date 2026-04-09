package mapper

import (
	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	typesv1 "github.com/neurochar/backend/pkg/proto_pb/common/types"
)

func FileToPb(file *fileEntity.File, fileUC fileUC.Usecase, isFull bool) *typesv1.File {
	if file == nil {
		return nil
	}

	if !isFull {
		return &typesv1.File{
			Url: fileUC.FileUrl(file),
		}
	}

	return &typesv1.File{
		Id:       file.ID.String(),
		Url:      fileUC.FileUrl(file),
		Filename: file.OriginalFileName,
	}
}
