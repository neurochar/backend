package interceptor

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc/metadata"
)

func getAllForwardedIPs(ctx context.Context) []string {
	var result []string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		vals := md.Get("x-forwarded-for")

		for _, v := range vals {
			parts := strings.Split(v, ",")
			for _, ip := range parts {
				ip = strings.TrimSpace(ip)
				if ip != "" {
					result = append(result, ip)
				}
			}
		}
	}

	return result
}

func getRequestID(ctx context.Context) *uuid.UUID {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get("x-request-id"); len(vals) > 0 {
			val, err := uuid.Parse(vals[0])
			if err == nil {
				return &val
			}
		}
	}

	return lo.ToPtr(uuid.New())
}

func getAuthorization(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get("authorization"); len(vals) > 0 {
			return vals[0]
		}
	}

	return ""
}
