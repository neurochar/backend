package gateway

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/middleware"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func CustomHTTPErrorHandler(
	ctx context.Context,
	mux *runtime.ServeMux,
	marshaler runtime.Marshaler,
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	appErr := appErrors.FromGRPCError(err)
	errJson := middleware.ErrorToHTTP(appErr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errJson.Code)
	// nolint
	_ = json.NewEncoder(w).Encode(errJson)
}
