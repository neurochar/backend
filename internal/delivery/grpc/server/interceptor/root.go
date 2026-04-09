package interceptor

import (
	"context"
	"net/netip"

	"github.com/neurochar/backend/internal/delivery/common/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func InterceptorRoot(trustedProxies []string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		ctx, enrich := tools.WithRequestEnrich(ctx)

		ipChain := getAllForwardedIPs(ctx)
		if p, ok := peer.FromContext(ctx); ok {
			remoteIP := tools.NormizeIP(p.Addr.String())
			if len(ipChain) == 0 || ipChain[len(ipChain)-1] != remoteIP {
				ipChain = append(ipChain, remoteIP)
			}
		}
		enrich.RequestIPChain = ipChain
		realIP, err := tools.ParseRealIP(ipChain, trustedProxies)
		if err == nil {
			ip, err := netip.ParseAddr(realIP)
			if err == nil {
				enrich.RequestIP = &ip
			}
		}

		enrich.RequestID = getRequestID(ctx)

		enrich.AuthorizationToken = getAuthorization(ctx)

		return handler(ctx, req)
	}
}
