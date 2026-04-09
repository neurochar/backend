package server

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"regexp"
	"strings"
	"time"

	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/delivery/common/tools"
)

func RootMiddleware(trustedProxies []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, _ := WithErrorHolder(r.Context())

			ctx, enrich := tools.WithRequestEnrich(ctx)

			ipChain := getAllForwardedIPs(r)
			remoteIP := tools.NormizeIP(r.RemoteAddr)
			if len(ipChain) == 0 || ipChain[len(ipChain)-1] != remoteIP {
				ipChain = append(ipChain, remoteIP)
			}
			enrich.RequestIPChain = ipChain
			realIP, err := tools.ParseRealIP(ipChain, trustedProxies)
			if err == nil {
				ip, err := netip.ParseAddr(realIP)
				if err == nil {
					enrich.RequestIP = &ip
				}
			}

			enrich.RequestID = getRequestID(r)

			enrich.AuthorizationToken = getAuthorization(r)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func LoggerMiddleware(isActive bool, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isActive {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()

			reqData := &tools.LogRequestData{
				Processor: "http",
				Method:    fmt.Sprintf("%s %s", r.Method, r.URL.Path),
				URI:       r.URL.Query().Encode(),
				Referer:   r.Header.Get("X-Referer"),
			}

			enrich := tools.GetEnrich(r.Context())
			if enrich != nil {
				if enrich.RequestIP != nil {
					reqData.IP = enrich.RequestIP.String()
				}

				reqData.IPChain = enrich.RequestIPChain

				reqData.RequestID = enrich.RequestID
			}

			ctx := tools.LogSetRequest(r.Context(), reqData)
			r = r.WithContext(ctx)

			rec := &ResponseWriter{ResponseWriter: w}

			next.ServeHTTP(rec, r)

			err := GetError(r.Context())

			duration := time.Since(start)

			if rec.Code == 0 {
				rec.Code = 200
			}

			respData := &tools.LogResponseData{
				Processor:  "http",
				DurationMS: duration.Milliseconds(),
				Code:       rec.Code,
			}

			if err != nil {
				respData.Error = err.Error()
				errStr := appErrors.ToJSONStruct(err, true, false)
				respData.AppError = &errStr
			}

			ctx = tools.LogSetResponse(r.Context(), respData)
			tools.LogHTTPContext(ctx, logger)
		})
	}
}

func RecoveryMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if errRec := recover(); errRec != nil {
					var err error
					switch errData := errRec.(type) {
					case error:
						err = errData
					case string:
						err = appErrors.ErrInternal.Extend(fmt.Sprintf("panic: %s", errData))
					default:
						err = appErrors.ErrInternal.Extend("panic: unknown error happend")
					}

					SetError(r.Context(), err)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Cors(corsAllowOrigins []string) func(next http.Handler) http.Handler {
	regexps := make([]*regexp.Regexp, 0, len(corsAllowOrigins))

	for _, raw := range corsAllowOrigins {
		raw = strings.TrimSpace(raw)

		if ip := net.ParseIP(raw); ip != nil {
			pattern := `^https?://` + regexp.QuoteMeta(raw) + `$`
			regexps = append(regexps, regexp.MustCompile(pattern))
			continue
		}

		pattern := `^https?://([a-zA-Z0-9-]+\.)?` + regexp.QuoteMeta(raw) + `$`
		regexps = append(regexps, regexp.MustCompile(pattern))
	}

	isAllowedOrigin := func(origin string) bool {
		for _, r := range regexps {
			if r.MatchString(origin) {
				return true
			}
		}
		return false
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			reqMethod := r.Header.Get("Access-Control-Request-Method")
			reqHeaders := r.Header.Get("Access-Control-Request-Headers")

			allowed := origin != "" && isAllowedOrigin(origin)

			// Аналогично Fiber: для не-wildcard origin добавляем ACAO только если origin разрешен.
			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Add("Vary", "Origin")
			}

			// Fiber считает preflight только OPTIONS + Access-Control-Request-Method.
			if r.Method == http.MethodOptions && reqMethod != "" {
				// Для preflight Fiber добавляет Vary по методу и заголовкам запроса.
				w.Header().Add("Vary", "Access-Control-Request-Method")
				w.Header().Add("Vary", "Access-Control-Request-Headers")

				if !allowed {
					w.WriteHeader(http.StatusNoContent)
					return
				}

				// У тебя в Fiber AllowMethods явно не задавался, но дефолт там есть.
				// Можно оставить такой набор.
				w.Header().Set("Access-Control-Allow-Methods", "GET,POST,HEAD,PUT,DELETE,PATCH")

				// Ключевой момент: как в Fiber, если AllowHeaders не задан,
				// возвращаем то, что браузер попросил.
				if reqHeaders != "" {
					w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
				}

				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
