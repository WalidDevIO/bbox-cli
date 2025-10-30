package client

import (
	"encoding/json"
	"fmt"
)

type StringOrInt string

func (s *StringOrInt) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch val := v.(type) {
	case string:
		*s = StringOrInt(val)
	case float64:
		*s = StringOrInt(fmt.Sprintf("%d", int(val)))
	default:
		return fmt.Errorf("cannot unmarshal %v into StringOrInt", v)
	}
	return nil
}

func (s StringOrInt) String() string {
	return string(s)
}
