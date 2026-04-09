package crm_tenant

import (
	"github.com/neurochar/backend/internal/app/config"
	crmUC "github.com/neurochar/backend/internal/domain/crm/usecase"
	desc "github.com/neurochar/backend/pkg/proto_pb/public/crm/v1"
	"google.golang.org/grpc"
)

type controller struct {
	desc.UnimplementedCrmTenantPublicServiceServer
	pkg       string
	cfg       config.Config
	crmFacade *crmUC.Facade
}

func Register(
	gRPCServer *grpc.Server,
	cfg config.Config,
	crmFacade *crmUC.Facade,
) {
	ctrl := &controller{
		pkg:       "grpc.Controller.CrmTenant",
		cfg:       cfg,
		crmFacade: crmFacade,
	}

	desc.RegisterCrmTenantPublicServiceServer(gRPCServer, ctrl)
}
