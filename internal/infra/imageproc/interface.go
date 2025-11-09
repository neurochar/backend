package imageproc

import (
	appErrors "github.com/neurochar/backend/internal/app/errors"
)

type ImageProcessor interface {
	IsOpenable(fileData []byte) bool

	ScaleAndCrop(fileData []byte, width int, height int, options ...option) ([]byte, *appErrors.AppError)

	DownscaleIfLarger(fileData []byte, maxWidth int, maxHeight int, options ...option) ([]byte, *appErrors.AppError)
}
