package tools

import (
	"context"
	"net/netip"

	"github.com/google/uuid"
)

type enrichHolderKey struct{}

type enrichHolder struct {
	RequestIP          *netip.Addr
	RequestIPChain     []string
	RequestID          *uuid.UUID
	AuthorizationToken string
	S2SToken           string
}

func WithRequestEnrich(ctx context.Context) (context.Context, *enrichHolder) {
	holder := &enrichHolder{}
	return context.WithValue(ctx, enrichHolderKey{}, holder), holder
}

func GetEnrich(ctx context.Context) *enrichHolder {
	holder, ok := ctx.Value(enrichHolderKey{}).(*enrichHolder)
	if !ok {
		return nil
	}
	return holder
}

func GetRealIP(ctx context.Context) *netip.Addr {
	enrich := GetEnrich(ctx)
	if enrich == nil || enrich.RequestIP == nil {
		return nil
	}

	return enrich.RequestIP
}

func GetRealIPString(ctx context.Context) string {
	ip := GetRealIP(ctx)
	if ip == nil {
		return ""
	}

	return ip.String()
}
