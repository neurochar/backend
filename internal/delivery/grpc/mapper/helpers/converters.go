package helpers

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
)

func PbAnyToAny(v *anypb.Any) (any, error) {
	b, err := protojson.Marshal(v)
	if err != nil {
		return nil, err
	}

	var res any
	if err := json.Unmarshal(b, &res); err != nil {
		return nil, err
	}

	return res, nil
}
