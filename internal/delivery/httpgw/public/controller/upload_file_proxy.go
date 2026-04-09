package controller

import (
	"context"
	"io"
	"net/http"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/httpgw/server"
)

func (ctrl *Controller) UploadFileProxy(
	exec func(ctx context.Context, file []byte, filename string) ([]byte, string, error),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			server.SetError(r.Context(), appErrors.ErrMethodNotAllowed)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			server.SetError(r.Context(), appErrors.ErrRequestEntityTooLarge.WithHints("file too large"))
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			server.SetError(r.Context(), appErrors.ErrBadRequest.WithHints("file not found"))
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			server.SetError(r.Context(), appErrors.ErrInternal.Extend("cannot read file"))
			return
		}

		resp, contentType, err := exec(r.Context(), data, header.Filename)
		if err != nil {
			server.SetError(r.Context(), err)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
	}
}
