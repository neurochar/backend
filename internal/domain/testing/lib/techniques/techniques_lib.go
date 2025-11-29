package techniques

import (
	"github.com/neurochar/backend/internal/domain/testing/entity"
	"github.com/neurochar/backend/internal/domain/testing/lib/techniques/kettel"
)

var TechniquesLib = map[uint64]entity.Technique{
	1: &kettel.Kettel,
}
