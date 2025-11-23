package auth

import (
	"context"
	"time"

	"github.com/neurochar/backend/internal/domain/tenant/entity"
)

func (uc *UsecaseImpl) delete(ctx context.Context, session *entity.Session) error {
	nowTime := time.Now()
	session.DeletedAt = &nowTime
	return uc.repo.Update(ctx, session)
}
