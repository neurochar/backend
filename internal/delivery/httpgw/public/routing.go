package public

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	authDelivery "github.com/neurochar/backend/internal/delivery/common/auth"
	"github.com/neurochar/backend/internal/delivery/httpgw/gateway"
	"github.com/neurochar/backend/internal/delivery/httpgw/server"
	authV1Pb "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
	crmV1Pb "github.com/neurochar/backend/pkg/proto_pb/public/crm/v1"
	registrationV1Pb "github.com/neurochar/backend/pkg/proto_pb/public/registration/v1"
	roomsV1Pb "github.com/neurochar/backend/pkg/proto_pb/public/rooms/v1"
	tenantV1Pb "github.com/neurochar/backend/pkg/proto_pb/public/tenant/v1"
	testingV1Pb "github.com/neurochar/backend/pkg/proto_pb/public/testing/v1"
	usersV1Pb "github.com/neurochar/backend/pkg/proto_pb/public/users_tenant/v1"
)

func (gw *Gateway) RegisterHandlers(
	ctx context.Context,
	authDeliveryCtrl *authDelivery.Controller,
) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/swagger.json", swaggerJsonHandler)
	mux.HandleFunc("/swagger", swaggerHandler)

	rmux := runtime.NewServeMux(
		runtime.WithErrorHandler(gateway.GRPCErrorHandler),
		runtime.WithMiddlewares(
			gateway.Recover(),
		),
	)

	err := crmV1Pb.RegisterCrmPublicServiceHandlerServer(ctx, rmux, gw.ctrl.Controls().CRM)
	if err != nil {
		return err
	}

	err = authV1Pb.RegisterAuthTenantPublicServiceHandlerServer(ctx, rmux, gw.ctrl.Controls().AuthTenant)
	if err != nil {
		return err
	}

	err = tenantV1Pb.RegisterTenantPublicServiceHandlerServer(ctx, rmux, gw.ctrl.Controls().Tenant)
	if err != nil {
		return err
	}

	err = registrationV1Pb.RegisterRegistrationPublicServiceHandlerServer(ctx, rmux, gw.ctrl.Controls().Registration)
	if err != nil {
		return err
	}

	err = roomsV1Pb.RegisterRoomsPublicServiceHandlerServer(ctx, rmux, gw.ctrl.Controls().Rooms)
	if err != nil {
		return err
	}

	err = usersV1Pb.RegisterUsersTenantPublicServiceHandlerServer(ctx, rmux, gw.ctrl.Controls().UsersTenant)
	if err != nil {
		return err
	}

	err = testingV1Pb.RegisterTestingPublicServiceHandlerServer(ctx, rmux, gw.ctrl.Controls().Testing)
	if err != nil {
		return err
	}

	publicMdwlrs := []func(http.Handler) http.Handler{
		PublicMiddleware(authDeliveryCtrl),
	}

	mux.Handle(
		"/v1/crm/candidates-resume",
		server.ChainMiddleware(
			http.HandlerFunc(gw.ctrl.UploadCandidateResumeFile),
			publicMdwlrs...,
		),
	)

	mux.Handle(
		"/v1/tenant/users/profile-photo",
		server.ChainMiddleware(
			http.HandlerFunc(gw.ctrl.UploadProfilePhotoFile),
			publicMdwlrs...,
		),
	)

	mux.Handle(
		"/",
		server.ChainMiddleware(
			rmux,
			publicMdwlrs...,
		),
	)

	gw.server.RegisterHandlers(mux)

	return nil
}
