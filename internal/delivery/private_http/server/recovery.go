package server

import (
	"fmt"
	"net/http"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if errRec := recover(); errRec != nil {
				var err error
				switch errData := errRec.(type) {
				case error:
					err = errData
				case string:
					err = fmt.Errorf("panic: %s", errData)
				default:
					err = fmt.Errorf("panic: unknown error")
				}

				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
