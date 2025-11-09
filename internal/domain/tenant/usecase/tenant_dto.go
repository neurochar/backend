package usecase

import "github.com/google/uuid"

type TenantListOptions struct {
	FilterIDs *[]uuid.UUID
}

type TenantIn struct{}

type CreateTenantIn struct {
	Data TenantIn
}

type UpdateTenantIn struct {
	Version int64
	Data    TenantIn
}
