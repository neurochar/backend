package gateway

import (
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/httpgw/server"
)

func Recover() runtime.Middleware {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			defer func() {
				if errRec := recover(); errRec != nil {
					var err error
					switch errData := errRec.(type) {
					case error:
						err = errData
					case string:
						err = appErrors.ErrInternal.Extend(fmt.Sprintf("panic: %s", errData))
					default:
						err = appErrors.ErrInternal.Extend("panic: unknown error happend")
					}

					server.SetError(r.Context(), err)
				}
			}()

			next(w, r, pathParams)
		}
	}
}
