package user

import (
	"context"

	"github.com/neurochar/backend/internal/domain/tenant/entity"
	"github.com/neurochar/backend/internal/domain/tenant/usecase"
)

// CreateByDTO TODO: сделать
func (uc *UsecaseImpl) CreateByDTO(
	ctx context.Context,
	in usecase.CreateTenantIn,
) (*entity.Tenant, error) {
	// const op = "CreateByDTO"

	return nil, nil
}
