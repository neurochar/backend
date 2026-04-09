package mapper

import (
	auth_tenantv1 "github.com/neurochar/backend/pkg/proto_pb/public/auth_tenant/v1"
)

func AuthTenantTokensToPb(refreshJwt string, refreshLifeSec int32, accessJwt string) *auth_tenantv1.Tokens {
	return &auth_tenantv1.Tokens{
		RefreshJwt:     refreshJwt,
		RefreshLifeSec: refreshLifeSec,
		AccessJwt:      accessJwt,
	}
}
