package call

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

func FromReader(reader io.Reader) (*wrap.JsonValue, error) {
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()

	var a any
	if err := decoder.Decode(&a); err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to read JSON")
	}

	ret, err := wrap.FromAny(a)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to read JSON")
	}

	return ret, nil
}

func ToAny(v *wrap.JsonValue) (any, error) {
	switch v.Type {
	case wrap.JsonTypeNull:
		return nil, nil
	case wrap.JsonTypeBoolean:
		return v.Bool(), nil
	case wrap.JsonTypeNumber:
		return json.Number(v.NumberValue), nil
	case wrap.JsonTypeString:
		return v.String(), nil
	case wrap.JsonTypeArray:
		a := []any{}
		for _, vi := range v.Array() {
			ai, err := ToAny(vi)
			if err != nil {
				return nil, errors.Wrap(
					errors.Join(err, errors.BadConversion),
					"fail to convert JsonValue to any")
			}

			a = append(a, ai)
		}

		return a, nil
	case wrap.JsonTypeObject:
		m := map[string]any{}
		for k, vk := range v.Object() {
			mk, err := ToAny(vk)
			if err != nil {
				return nil, errors.Wrap(
					errors.Join(err, errors.BadConversion),
					"fail to convert JsonValue to any")
			}

			m[k] = mk
		}

		return m, nil
	default:
		return nil, errors.Wrap(errors.BadState, "unexpected case %v", v.Type)
	}
}

func ToReader(v *wrap.JsonValue) (*bytes.Buffer, error) {
	a, err := ToAny(v)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to convert JsonValue to any")
	}

	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(a); err != nil {
		return buffer, errors.Wrap(
			errors.Join(err, errors.IOFailure),
			"fail to write JSON")
	}

	return buffer, nil
}
