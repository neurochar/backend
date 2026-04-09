package controller

import (
	"github.com/neurochar/backend/internal/delivery/grpc/backend/controller/auth_tenant"
	"github.com/neurochar/backend/internal/delivery/grpc/backend/controller/crm_tenant"
	"go.uber.org/fx"
)

// FxModule - fx module
var FxModule = fx.Options(
	fx.Invoke(auth_tenant.Register),
	fx.Invoke(crm_tenant.Register),
)
