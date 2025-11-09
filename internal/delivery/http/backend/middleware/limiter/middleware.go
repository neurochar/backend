package limiter

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/http/backend/middleware"
	httpMiddleware "github.com/neurochar/backend/internal/delivery/http/middleware"
	"github.com/ulule/limiter/v3"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

type Middleware struct {
	limiter *limiter.Limiter
}

func newMiddleware(rateFormatted string) *Middleware {
	rate, err := limiter.NewRateFromFormatted(rateFormatted)
	if err != nil {
		panic(err)
	}

	store := memory.NewStore()

	mw := &Middleware{
		limiter: limiter.New(store, rate),
	}

	return mw
}

var ErrRateLimiter = appErrors.ErrInternal.Extend("rate limiter error")

func (m *Middleware) Create(useIP bool, useAccountID bool, postfix string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var key strings.Builder

		if useIP {
			key.WriteString(httpMiddleware.GetRealIP(c))
		}

		if useAccountID {
			authData := middleware.GetAuthData(c)
			if authData != nil {
				key.WriteString(":")
				key.WriteString(authData.AccountID.String())
			}
		}

		if postfix != "" {
			key.WriteString(":")
			key.WriteString(postfix)
		}

		lctx, err := m.limiter.Get(c.Context(), key.String())
		if err != nil {
			return ErrRateLimiter.WithWrap(err)
		}

		if lctx.Reached {
			resetTime := time.Unix(lctx.Reset, 0)
			retryAfter := time.Until(resetTime).Seconds()
			if retryAfter < 1 {
				retryAfter = 1
			}

			c.Set("Retry-After", fmt.Sprintf("%d", int64(math.Ceil(retryAfter))))
			return appErrors.ErrTooManyRequests.WithDetail("retry_after_sec", false, int64(math.Ceil(retryAfter)))
		}

		return c.Next()
	}
}
