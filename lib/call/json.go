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

func ToReader(v *wrap.JsonValue) (*bytes.Buffer, error) {
	a := wrap.ToAny(v)

	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	if err := encoder.Encode(a); err != nil {
		return buffer, errors.Wrap(
			errors.Join(err, errors.IOFailure),
			"fail to write JSON")
	}

	return buffer, nil
}
