package gateway

import (
	"context"
	"net/http"

	"github.com/neurochar/backend/internal/delivery/http/middleware"
	"google.golang.org/grpc/metadata"
)

func MetadataAnnotator(ctx context.Context, r *http.Request) metadata.MD {
	md := metadata.MD{}

	if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
		md.Set("x-request-id", requestID)
	}

	authToken, ok := ctx.Value(middleware.AuthCtxKeyToken).(string)
	if ok && authToken != "" {
		md.Set("authorization", authToken)
	}

	return md
}
