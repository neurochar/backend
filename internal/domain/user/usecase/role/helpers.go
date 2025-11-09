package role

import (
	"context"
	"time"

	userEntity "github.com/neurochar/backend/internal/domain/user/entity"
)

func (uc *UsecaseImpl) deleteRole(ctx context.Context, item *userEntity.Role) error {
	nowTime := time.Now()
	item.DeletedAt = &nowTime

	return uc.repoRole.Update(ctx, item)
}
