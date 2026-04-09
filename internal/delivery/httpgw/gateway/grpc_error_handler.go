package gateway

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/httpgw/server"
)

func GRPCErrorHandler(
	ctx context.Context,
	_ *runtime.ServeMux,
	_ runtime.Marshaler,
	_ http.ResponseWriter,
	_ *http.Request,
	err error,
) {
	appErr := appErrors.FromGRPCError(err)
	if appErr != nil {
		server.SetError(ctx, appErr)
	}
}
