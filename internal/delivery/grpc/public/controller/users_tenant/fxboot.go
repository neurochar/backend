package users_tenant

import (
	"go.uber.org/fx"
)

var FxModule = fx.Module(
	"users_tenant_grpc_public_controller",
	fx.Provide(New),
	fx.Invoke(func(ctrl *Controller) {
		ctrl.Register()
	}),
)
