package auth

import "context"

const (
	ContextKeyAuthCheckTenantAccess contextKey = "auth_context_auth_check_right"
)

func WithoutCheckTenantAccess(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextKeyAuthCheckTenantAccess, false)
}

func WithCheckTenantAccess(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextKeyAuthCheckTenantAccess, true)
}

func IsNeedToCheckTenantAccess(ctx context.Context) bool {
	is, ok := ctx.Value(ContextKeyAuthCheckTenantAccess).(bool)
	return ok && is
}
