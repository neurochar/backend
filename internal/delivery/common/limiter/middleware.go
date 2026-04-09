package limiter

import (
	"context"
	"math"
	"net/netip"
	"strings"
	"time"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/ulule/limiter/v3"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

type Limiter struct {
	limiter *limiter.Limiter
}

func newLimiter(rateFormatted string) *Limiter {
	rate, err := limiter.NewRateFromFormatted(rateFormatted)
	if err != nil {
		panic(err)
	}

	store := memory.NewStore()

	mw := &Limiter{
		limiter: limiter.New(store, rate),
	}

	return mw
}

var ErrRateLimiter = appErrors.ErrInternal.Extend("rate limiter error")

type RegisterKey struct {
	IP        *netip.Addr
	Key       string
	AccountID string
}

func (m *Limiter) Register(ctx context.Context, reqKey *RegisterKey) error {
	if reqKey == nil {
		return nil
	}

	var key strings.Builder

	if reqKey.IP != nil {
		key.WriteString(":ip:")
		key.WriteString(reqKey.IP.String())
	}

	if reqKey.AccountID != "" {
		key.WriteString(":accountID:")
		key.WriteString(reqKey.AccountID)
	}

	if reqKey.Key != "" {
		key.WriteString(":key:")
		key.WriteString(reqKey.Key)
	}

	keyValue := key.String()

	if keyValue == "" {
		return nil
	}

	lctx, err := m.limiter.Get(ctx, keyValue)
	if err != nil {
		return ErrRateLimiter.WithWrap(err)
	}

	if lctx.Reached {
		resetTime := time.Unix(lctx.Reset, 0)
		retryAfter := time.Until(resetTime).Seconds()
		if retryAfter < 1 {
			retryAfter = 1
		}

		return appErrors.ErrTooManyRequests.WithDetail("retry_after_sec", false, int64(math.Ceil(retryAfter)))
	}

	return nil
}
