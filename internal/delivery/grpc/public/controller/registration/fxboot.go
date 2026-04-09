package registration

import (
	"go.uber.org/fx"
)

var FxModule = fx.Module(
	"registration_grpc_public_controller",
	fx.Provide(New),
	fx.Invoke(func(ctrl *Controller) {
		ctrl.Register()
	}),
)
