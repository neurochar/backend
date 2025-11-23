// Package v1 contains v1 http handlers
package v1

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/neurochar/backend/internal/app/config"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware/limiter"
	"github.com/neurochar/backend/pkg/backoff"
)

type Groups struct {
	Prefix      string
	Default     fiber.Router
	RateLimiter *limiter.Controller
}

const BackoffDefaultGroupID = "default"

// ProvideGroups - provide v1 group
func ProvideGroups(
	cfg config.Config,
	fiberApp *fiber.App,
	cpanelMdwr *middleware.Controller,
	backoffCtrl *backoff.Controller,
) *Groups {
	prefix := fmt.Sprintf("%s/v1", cfg.BackendApp.HTTP.Prefix)

	defaultGroup := fiberApp.Group(prefix, cpanelMdwr.Auth)

	backoffCtrl.SetConfigForGroup(
		BackoffDefaultGroupID,
		backoff.WithTtl(time.Minute*20),
		backoff.WithInitialInterval(time.Second*5),
		backoff.WithMultiplier(2),
		backoff.WithMaxInterval(time.Minute*5),
	)

	return &Groups{
		Prefix:      prefix,
		Default:     defaultGroup,
		RateLimiter: limiter.NewController(),
	}
}
