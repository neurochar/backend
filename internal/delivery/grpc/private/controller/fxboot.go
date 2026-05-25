package controller

import (
	"github.com/neurochar/backend/internal/delivery/grpc/private/controller/crm"
	"go.uber.org/fx"
)

// FxModule - fx module
var FxModule = fx.Options(
	crm.FxModule,
)
