package backend

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/neurochar/backend/internal/delivery/http/middleware"
	alertUC "github.com/neurochar/backend/internal/domain/alert/usecase"
)

const defaultBodyLimit = 10 * 1024 * 1024

// HTTPConfig - config for http server
type HTTPConfig struct {
	AppTitle         string
	UnderProxy       bool
	UseLogger        bool
	BodyLimit        int
	CorsAllowOrigins []string
}

// NewHTTPFiber provides fiber app
func NewHTTPFiber(httpCfg HTTPConfig, logger *slog.Logger, alertUsecase alertUC.Usecase) *fiber.App {
	if httpCfg.BodyLimit == -1 {
		httpCfg.BodyLimit = defaultBodyLimit
	}

	fiberCfg := fiber.Config{
		ErrorHandler: middleware.ErrorHandler(httpCfg.AppTitle, logger, alertUsecase),
		BodyLimit:    httpCfg.BodyLimit,
	}

	if httpCfg.UnderProxy {
		fiberCfg.ProxyHeader = fiber.HeaderXForwardedFor
	}

	app := fiber.New(fiberCfg)

	app.Use(middleware.Recovery(logger))
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestIP())

	if len(httpCfg.CorsAllowOrigins) > 0 {
		app.Use(middleware.Cors(httpCfg.CorsAllowOrigins))
	}

	app.Use(middleware.XCheck())

	if httpCfg.UseLogger {
		app.Use(middleware.Logger(logger))
	}

	return app
}
