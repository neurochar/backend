package tenant

import (
	"go.uber.org/fx"
)

var FxModule = fx.Module(
	"tenant_grpc_public_controller",
	fx.Provide(New),
	fx.Invoke(func(ctrl *Controller) {
		ctrl.Register()
	}),
)
