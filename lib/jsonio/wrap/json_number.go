package wrap

import "encoding/json"

type JsonNumber json.Number

func (n JsonNumber) Int64() (int64, bool) {
	i, err := json.Number(n).Int64()
	return i, err == nil
}
func (n JsonNumber) Float64() (float64, bool) {
	f, err := json.Number(n).Float64()
	return f, err == nil
}
