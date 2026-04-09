package auth_tenant

import (
	"github.com/neurochar/backend/internal/app/config"
	fileUC "github.com/neurochar/backend/internal/domain/file/usecase"
	tenantUC "github.com/neurochar/backend/internal/domain/tenant/usecase"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
	"google.golang.org/grpc"
)

type controller struct {
	desc.UnimplementedAuthTenantPublicServiceServer
	pkg          string
	cfg          config.Config
	tenantFacade *tenantUC.Facade
	fileUC       fileUC.Usecase
}

func Register(
	gRPCServer *grpc.Server,
	cfg config.Config,
	tenantFacade *tenantUC.Facade,
	fileUC fileUC.Usecase,
) {
	ctrl := &controller{
		pkg:          "grpc.Controller.AuthTenant",
		cfg:          cfg,
		tenantFacade: tenantFacade,
		fileUC:       fileUC,
	}

	desc.RegisterAuthTenantPublicServiceServer(gRPCServer, ctrl)
}
