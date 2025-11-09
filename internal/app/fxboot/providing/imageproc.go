package providing

import "github.com/neurochar/backend/internal/infra/imageproc"

func NewImageProc() imageproc.ImageProcessor {
	return imageproc.New()
}
