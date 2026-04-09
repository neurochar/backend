package gateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/neurochar/backend/internal/app/config"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/app/fxboot/invoking"
	"github.com/neurochar/backend/internal/delivery/http/gateway"
	authV1Pb "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
	crmV1Pb "github.com/neurochar/backend/pkg/proto_pb/public/crm/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitGatewayHandler(fiberApp *fiber.App, cfg config.Config) invoking.InvokeInit {
	const op = "Delivery.HTTP.NewGatewayHandler"

	return invoking.InvokeInit{
		StartBeforeOpen: func(ctx context.Context) error {
			gw, err := newGatewayHandler(ctx, fmt.Sprintf("127.0.0.1:%d", cfg.BackendApp.GRPC.Port))
			if err != nil {
				return appErrors.Chainf(err, "%s", op)
			}

			fiberApp.Use(adaptor.HTTPHandler(gw))

			return nil
		},
	}
}

func newGatewayHandler(ctx context.Context, grpcAddr string) (http.Handler, error) {
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(gateway.CustomHTTPErrorHandler),
		runtime.WithMetadata(gateway.MetadataAnnotator),
	)

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := crmV1Pb.RegisterCrmTenantPublicServiceHandlerFromEndpoint(ctx, mux, grpcAddr, dialOpts)
	if err != nil {
		return nil, err
	}

	err = authV1Pb.RegisterAuthTenantPublicServiceHandlerFromEndpoint(ctx, mux, grpcAddr, dialOpts)
	if err != nil {
		return nil, err
	}

	return stripAuthorization(mux), nil
}

func stripAuthorization(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r2 := r.Clone(r.Context())
		r2.Header = r.Header.Clone()
		r2.Header.Del("Authorization")
		r2.Header.Del("Grpc-Metadata-Authorization")
		r2.Header.Del("Grpc-Metadata-S2s")
		h.ServeHTTP(w, r2)
	})
}
