package mapper

import (
	tenantEntity "github.com/neurochar/backend/internal/domain/tenant/entity"
	typesv1 "github.com/neurochar/backend/pkg/proto_pb/common/types"
)

func TenantToPb(item *tenantEntity.Tenant) *typesv1.Tenant {
	if item == nil {
		return nil
	}

	return &typesv1.Tenant{
		Id:      item.ID.String(),
		Version: item.Version(),
		TextId:  item.TextID,
		Name:    item.Name,
	}
}
