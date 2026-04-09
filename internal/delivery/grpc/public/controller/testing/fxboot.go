package testing

import (
	"go.uber.org/fx"
)

var FxModule = fx.Module(
	"testing_grpc_public_controller",
	fx.Provide(New),
	fx.Invoke(func(ctrl *Controller) {
		ctrl.Register()
	}),
)
