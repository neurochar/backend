package null

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Nullable[T any] struct {
	IsSet   bool
	IsValid bool
	Value   T
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	n.IsSet = true

	if bytes.Equal(data, []byte("null")) {
		n.IsValid = false
		var zero T
		n.Value = zero
		return nil
	}

	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("Nullable: failed to unmarshal into target type: %w", err)
	}

	n.IsValid = true
	n.Value = v
	return nil
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.IsSet || !n.IsValid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Value)
}
