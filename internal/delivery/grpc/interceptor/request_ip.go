package interceptor

import (
	"context"
	"net"
	"strings"

	"github.com/neurochar/backend/internal/delivery/common"
	"github.com/neurochar/backend/internal/infra/loghandler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func getClientIPFromGRPCRequest(privateIPs []string) func(ctx context.Context) string {
	return func(ctx context.Context) string {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
				ips := strings.Split(xff[0], ",")
				clientIPs := make([]string, 0, len(xff))
				for _, ip := range ips {
					v := strings.TrimSpace(ip)
					if v != "" {
						clientIPs = append(clientIPs, v)
					}
				}

				realIP, err := common.ParseRealIP(clientIPs, privateIPs)
				if err == nil {
					return realIP
				}
			}
		}

		if p, ok := peer.FromContext(ctx); ok {
			host, _, err := net.SplitHostPort(p.Addr.String())
			if err == nil {
				return host
			}
			return p.Addr.String()
		}

		return ""
	}
}

type contextReqIPKey string

const (
	contextReqIPKeyValue contextReqIPKey = "x-request-ip"
)

func InterceptorRequestIP(privateIPs []string) grpc.UnaryServerInterceptor {
	ipGetter := getClientIPFromGRPCRequest(privateIPs)

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		requestIP := ipGetter(ctx)

		ctx = context.WithValue(ctx, contextReqIPKeyValue, requestIP)
		ctx = loghandler.SetContextData(ctx, "request.ip", requestIP)

		return handler(ctx, req)
	}
}

func GetClientIP(ctx context.Context) string {
	if v := ctx.Value(contextReqIPKeyValue); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}

	return ""
}
