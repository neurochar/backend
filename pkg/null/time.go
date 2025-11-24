package null

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type NullableTime struct {
	IsSet   bool
	IsValid bool
	Time    time.Time
}

var timeLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02",
}

func (nt *NullableTime) UnmarshalJSON(data []byte) error {
	nt.IsSet = true

	if bytes.Equal(data, []byte("null")) {
		nt.IsValid = false
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("NullableTime: expected string or null, got %s: %w", string(data), err)
	}

	s = strings.TrimSpace(s)
	if s == "" {
		nt.IsValid = false
		return nil
	}

	var lastErr error
	for _, layout := range timeLayouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			nt.IsValid = true
			nt.Time = t
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("NullableTime: cannot parse %q as time: %w", s, lastErr)
}

func (nt NullableTime) MarshalJSON() ([]byte, error) {
	if !nt.IsSet || !nt.IsValid {
		return []byte("null"), nil
	}

	return json.Marshal(nt.Time)
}
