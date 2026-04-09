package backend

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/neurochar/backend/internal/delivery/http/middleware"
	alertUC "github.com/neurochar/backend/internal/domain/alert/usecase"
)

// HTTPConfig - config for http server
type HTTPConfig struct {
	AppTitle         string
	UseLogger        bool
	BodyLimit        int
	CorsAllowOrigins []string
	ServerIPs        []string
}

// NewHTTPFiber provides fiber app
func NewHTTPFiber(httpCfg HTTPConfig, logger *slog.Logger, alertUsecase alertUC.Usecase) *fiber.App {
	fiberCfg := fiber.Config{
		ErrorHandler: middleware.ErrorHandler(httpCfg.AppTitle, logger, alertUsecase),
		BodyLimit:    httpCfg.BodyLimit,
	}

	fiberCfg.ProxyHeader = fiber.HeaderXForwardedFor

	app := fiber.New(fiberCfg)

	app.Use(middleware.Recovery(logger))
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestIP(httpCfg.ServerIPs))

	if len(httpCfg.CorsAllowOrigins) > 0 {
		app.Use(middleware.Cors(httpCfg.CorsAllowOrigins))
	}

	app.Use(middleware.Authorization())
	app.Use(middleware.XCheck())

	if httpCfg.UseLogger {
		app.Use(middleware.Logger(logger))
	}

	return app
}
