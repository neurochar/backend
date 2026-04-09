package controller

import (
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/auth_tenant"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/crm"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/registration"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/rooms"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/tenant"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/testing"
	"github.com/neurochar/backend/internal/delivery/grpc/public/controller/users_tenant"
	"go.uber.org/fx"
)

// FxModule - fx module
var FxModule = fx.Options(
	tenant.FxModule,
	auth_tenant.FxModule,
	users_tenant.FxModule,
	crm.FxModule,
	registration.FxModule,
	rooms.FxModule,
	testing.FxModule,
)
