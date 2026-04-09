package server

import (
	"errors"
	"net/http"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

func ErrorCatcherMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rec := &ResponseWriter{ResponseWriter: w}

			next.ServeHTTP(rec, r)

			err := GetError(r.Context())
			if err == nil {
				switch rec.Code {
				case http.StatusNotFound:
					err = appErrors.ErrNotFound
				default:
					if rec.Code >= 400 && rec.Code < 500 {
						err = appErrors.ErrBadRequest
					} else if rec.Code >= 500 {
						err = appErrors.ErrInternal
					}
				}
			} else {
				_, ok := appErrors.ExtractError(err)
				if !ok {
					err = appErrors.ErrInternal.WithWrap(err)
				}
			}

			if err != nil && errors.Is(err, appErrors.ErrUnimplemented) {
				err = appErrors.ErrMethodNotAllowed.WithTextCode("METHOD_NOT_ALLOWED")
			}

			if err != nil {
				SetError(r.Context(), err)
			}
		})
	}
}
