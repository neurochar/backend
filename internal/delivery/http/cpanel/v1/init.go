// Package v1 contains v1 http handlers
package v1

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/http/cpanel/middleware"
	"github.com/neurochar/backend/internal/delivery/http/cpanel/middleware/limiter"
)

type Groups struct {
	Prefix      string
	Default     fiber.Router
	RateLimiter *limiter.Controller
}

// ProvideGroups - provide v1 group
func ProvideGroups(cfg config.Config, fiberApp *fiber.App, cpanelMdwr *middleware.Controller) *Groups {
	prefix := fmt.Sprintf("%s/v1", cfg.CPanelApp.HTTP.Prefix)

	defaultGroup := fiberApp.Group(prefix, cpanelMdwr.MiddlewareAuth)

	return &Groups{
		Prefix:      prefix,
		Default:     defaultGroup,
		RateLimiter: limiter.NewController(),
	}
}
