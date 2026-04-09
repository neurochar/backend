package public

import (
	"net/http"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	authDelivery "github.com/neurochar/backend/internal/delivery/common/auth"
	"github.com/neurochar/backend/internal/delivery/httpgw/server"
)

func PublicMiddleware(
	authDeliveryCtrl *authDelivery.Controller,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authCtx, err := authDeliveryCtrl.EnrichAuth()(r.Context())
			if err != nil {
				server.SetError(r.Context(), appErrors.ErrInternal.WithWrap(err))
				return
			} else {
				r = r.WithContext(authCtx)
			}

			next.ServeHTTP(w, r)
		})
	}
}
