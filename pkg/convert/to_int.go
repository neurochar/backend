package convert

import (
	"encoding/json"
	"reflect"
)

func ToInt(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int8, int16, int32, int64:
		return int(reflect.ValueOf(n).Int()), true
	case uint, uint8, uint16, uint32, uint64:
		return int(reflect.ValueOf(n).Uint()), true
	case float32, float64:
		return int(reflect.ValueOf(n).Float()), true
	case json.Number:
		if i, err := n.Int64(); err == nil {
			return int(i), true
		}
		if f, err := n.Float64(); err == nil {
			return int(f), true
		}
	}

	return 0, false
}
