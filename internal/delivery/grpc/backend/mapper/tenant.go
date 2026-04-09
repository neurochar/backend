package mapper

import (
	tenantEntity "github.com/neurochar/backend/internal/domain/tenant/entity"
	tenantv1 "github.com/neurochar/backend/pkg/proto_pb/public/tenant/v1"
)

func TenantToPb(item *tenantEntity.Tenant) *tenantv1.Tenant {
	if item == nil {
		return nil
	}

	return &tenantv1.Tenant{
		Id:      item.ID.String(),
		Version: item.Version(),
		TextId:  item.TextID,
		Name:    item.Name,
	}
}
