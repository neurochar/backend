package usecase

import (
	"time"

	"github.com/neurochar/backend/internal/common/uctypes"
)

type ListOptions struct {
	FilterSentAtCompare *uctypes.CompareOption[*time.Time]
}
