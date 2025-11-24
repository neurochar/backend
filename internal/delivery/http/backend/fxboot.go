package backend

import (
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	v1 "github.com/neurochar/backend/internal/delivery/http/backend/v1"
	"github.com/neurochar/backend/internal/delivery/http/backend/v1/auth"
	"github.com/neurochar/backend/internal/delivery/http/backend/v1/crm"
	"github.com/neurochar/backend/internal/delivery/http/backend/v1/registration"
	"github.com/neurochar/backend/internal/delivery/http/backend/v1/tenants"
	"github.com/neurochar/backend/internal/delivery/http/backend/v1/testing"
	"github.com/neurochar/backend/internal/delivery/http/backend/v1/tests"
	"github.com/neurochar/backend/internal/delivery/http/backend/v1/users"
	"github.com/neurochar/backend/pkg/validation"
	"go.uber.org/fx"
)

// FxModule - fx module
var FxModule = fx.Options(
	fx.Provide(validation.New),
	fx.Options(
		fx.Provide(middleware.New),
		fx.Provide(v1.ProvideGroups),
		tests.FxModule,
		auth.FxModule,
		tenants.FxModule,
		registration.FxModule,
		users.FxModule,
		crm.FxModule,
		testing.FxModule,
	),
)
