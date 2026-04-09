package server

import (
	"encoding/json"
	"errors"
	"net/http"

	appErrors "github.com/neurochar/backend/internal/app/errors"
)

type ErrorJSON struct {
	Code     int            `json:"code"`
	TextCode string         `json:"textCode"`
	Hints    []string       `json:"hints"`
	Details  map[string]any `json:"details"`
}

type HTTPError struct {
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Message
}

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func httpTextCode(code int) string {
	switch {
	case code == http.StatusMethodNotAllowed:
		return "METHOD_NOT_ALLOWED"
	case code >= 400 && code < 500:
		return "BAD_REQUEST"
	default:
		return "INTERNAL_ERROR"
	}
}

func ErrorToHTTP(err error) ErrorJSON {
	code := http.StatusInternalServerError
	jsonRes := ErrorJSON{
		Code:     code,
		TextCode: "INTERNAL_ERROR",
		Hints:    []string{},
		Details:  map[string]any{},
	}

	if appError, ok := appErrors.ExtractError(err); ok {
		jsonRes.Code = int(appError.Meta().Code)
		jsonRes.TextCode = appError.Meta().TextCode
		jsonRes.Hints = appError.Hints()
		if jsonRes.Hints == nil {
			jsonRes.Hints = []string{}
		}
		jsonRes.Details = appError.Details(false)

		return jsonRes
	}

	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		jsonRes.Code = httpErr.Code
		jsonRes.TextCode = httpTextCode(httpErr.Code)

		if httpErr.Message != "" {
			jsonRes.Hints = []string{httpErr.Message}
		}

		return jsonRes
	}

	return jsonRes
}

func ErrorHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			err := GetError(r.Context())
			if err != nil {
				jsonRes := ErrorToHTTP(err)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(jsonRes.Code)
				_ = json.NewEncoder(w).Encode(jsonRes)
			}
		})
	}
}
