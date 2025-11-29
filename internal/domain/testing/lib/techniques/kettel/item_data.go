package kettel

import (
	appErrors "github.com/neurochar/backend/internal/app/errors"
	"github.com/neurochar/backend/internal/domain/testing/entity"
)

type KettelItemDataImpl struct {
	ID uint64 `json:"id"`
}

func (i *KettelItemDataImpl) GetItem() (entity.TechniqueItem, error) {
	item, ok := ItemsLib[i.ID]
	if !ok {
		return nil, appErrors.ErrNotFound
	}

	return item, nil
}
