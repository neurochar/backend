package common

import (
	"github.com/neurochar/backend/internal/delivery/common/auth"
	"github.com/neurochar/backend/internal/delivery/common/limiter"
	"go.uber.org/fx"
)

var FxModule = fx.Module(
	"delivery_common_module",
	fx.Provide(limiter.NewController),
	fx.Provide(auth.New),
)
