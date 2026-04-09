package backend

import (
	"github.com/neurochar/backend/internal/delivery/http/backend/controller"
	"github.com/neurochar/backend/internal/delivery/http/backend/controller/auth"
	"github.com/neurochar/backend/internal/delivery/http/backend/controller/crm_tenant"
	"github.com/neurochar/backend/internal/delivery/http/backend/controller/registration"
	"github.com/neurochar/backend/internal/delivery/http/backend/controller/rooms"
	"github.com/neurochar/backend/internal/delivery/http/backend/controller/tenants"
	"github.com/neurochar/backend/internal/delivery/http/backend/controller/testing"
	"github.com/neurochar/backend/internal/delivery/http/backend/controller/tests"
	"github.com/neurochar/backend/internal/delivery/http/backend/controller/users"
	"github.com/neurochar/backend/internal/delivery/http/backend/gateway"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	"github.com/neurochar/backend/pkg/validation"
	"go.uber.org/fx"
)

// FxModule - fx module
var FxModule = fx.Options(
	fx.Provide(validation.New),
	fx.Options(
		fx.Provide(gateway.NewGrpcClient),
		fx.Provide(middleware.New),
		fx.Provide(controller.ProvideGroups),
		tests.FxModule,
		auth.FxModule,
		tenants.FxModule,
		registration.FxModule,
		users.FxModule,
		crm_tenant.FxModule,
		testing.FxModule,
		rooms.FxModule,
	),
	fx.Provide(
		fx.Annotate(gateway.InitGatewayHandler, fx.ResultTags(`group:"InvokeInit"`)),
		fx.Annotate(gateway.InitGrpcClient, fx.ResultTags(`group:"InvokeInit"`)),
	),
)
