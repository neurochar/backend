package cpanel

import (
	"github.com/neurochar/backend/internal/delivery/http/cpanel/middleware"
	v1 "github.com/neurochar/backend/internal/delivery/http/cpanel/v1"
	"github.com/neurochar/backend/internal/delivery/http/cpanel/v1/tests"
	"github.com/neurochar/backend/internal/delivery/http/cpanel/v1/users"
	"github.com/neurochar/backend/pkg/validation"
	"go.uber.org/fx"
)

// FxModule - fx module
var FxModule = fx.Options(
	fx.Provide(validation.New),
	fx.Options(
		fx.Provide(middleware.New),
		fx.Provide(v1.ProvideGroups),
		users.FxModule,
		tests.FxModule,
	),
)
