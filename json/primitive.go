package json

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func ToNullableString(value any) (p *string, err error) {
	var s string
	switch value := value.(type) {
	case nil:
		return nil, nil
	case json.Number:
		s = value.String()
	case int:
		s = strconv.FormatInt(int64(value), 10)
	case int8:
		s = strconv.FormatInt(int64(value), 10)
	case int16:
		s = strconv.FormatInt(int64(value), 10)
	case int32:
		s = strconv.FormatInt(int64(value), 10)
	case int64:
		s = strconv.FormatInt(int64(value), 10)
	case uint:
		s = strconv.FormatUint(uint64(value), 10)
	case uint8:
		s = strconv.FormatUint(uint64(value), 10)
	case uint16:
		s = strconv.FormatUint(uint64(value), 10)
	case uint32:
		s = strconv.FormatUint(uint64(value), 10)
	case uint64:
		s = strconv.FormatUint(uint64(value), 10)
	case string:
		s = value
	case bool:
		s = strconv.FormatBool(value)
	default:
		return nil, fmt.Errorf("unexpected value %v", value)
	}
	return &s, nil
}

func ToNullableInteger(value any) (p *int64, err error) {
	var i int64
	switch value := value.(type) {
	case nil:
		return nil, nil
	case json.Number:
		i, err = value.Int64()
		if err != nil {
			return nil, fmt.Errorf("unexpected value %v", value)
		}
	case int:
		i = int64(value)
	case int8:
		i = int64(value)
	case int16:
		i = int64(value)
	case int32:
		i = int64(value)
	case int64:
		i = int64(value)
	case uint:
		i = int64(value)
	case uint8:
		i = int64(value)
	case uint16:
		i = int64(value)
	case uint32:
		i = int64(value)
	case uint64:
		i = int64(value)
	case string:
		i, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("unexpected value %v", value)
		}
	default:
		return nil, fmt.Errorf("unexpected value %v", value)
	}
	return &i, nil
}

func ToNullableFloat(value any) (p *float64, err error) {
	var f float64
	switch value := value.(type) {
	case nil:
		return nil, nil
	case json.Number:
		f, err = value.Float64()
		if err != nil {
			return nil, fmt.Errorf("unexpected value %v", value)
		}
	case int:
		f = float64(value)
	case int8:
		f = float64(value)
	case int16:
		f = float64(value)
	case int32:
		f = float64(value)
	case int64:
		f = float64(value)
	case uint:
		f = float64(value)
	case uint8:
		f = float64(value)
	case uint16:
		f = float64(value)
	case uint32:
		f = float64(value)
	case uint64:
		f = float64(value)
	case string:
		f, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("unexpected value %v", value)
		}
	default:
		return nil, fmt.Errorf("unexpected value %v", value)
	}
	return &f, nil
}

func ToNullableBoolean(value any) (p *bool, err error) {
	var b bool
	switch value := value.(type) {
	case nil:
		return nil, nil
	case json.Number:
		i, err := value.Int64()
		if err != nil {
			return nil, fmt.Errorf("unexpected value %v", value)
		}
		b = i != 0
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		b = value != 0
	case string:
		if s := strings.ToLower(value); s == "true" || s == "1" {
			b = true
		} else if s == "false" || s == "0" {
			b = false
		} else {
			return nil, fmt.Errorf("unexpected value %v", value)
		}
	case bool:
		b = value
	default:
		return nil, fmt.Errorf("unexpected value %v", value)
	}
	return &b, nil
}
