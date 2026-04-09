package mapper

import (
	fileEntity "github.com/neurochar/backend/internal/domain/file/entity"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	typesv1 "github.com/neurochar/backend/pkg/proto_pb/common/types"
	"github.com/samber/lo"
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
		Url:      fileUC.FileUrl(file),
		Id:       lo.ToPtr(file.ID.String()),
		Filename: lo.ToPtr(file.OriginalFileName),
	}
}

func FilesToMapPb(files []*fileEntity.File, fileUC fileUC.Usecase, isFull bool) map[string]*typesv1.File {
	result := make(map[string]*typesv1.File, len(files))

	for _, file := range files {
		result[file.Target] = FileToPb(file, fileUC, isFull)
	}

	return result
}
