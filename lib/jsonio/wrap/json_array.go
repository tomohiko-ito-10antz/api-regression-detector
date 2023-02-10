package wrap

import "fmt"

type JsonArray []*JsonValue

func (a JsonArray) Empty() bool {
	return len(a) == 0
}
func (a JsonArray) Get(i int) *JsonValue {
	if i < 0 || a.Len() <= i {
		panic(fmt.Errorf("index %d > len %d", i, a.Len()))
	}
	return a[i]
}
func (a JsonArray) Set(i int, v *JsonValue) {
	if i < 0 || a.Len() <= i {
		panic(fmt.Errorf("index %d > len %d", i, a.Len()))
	}
	if v == nil {
		v = Null()
	}
	a[i] = v
}
func (a JsonArray) Append(v *JsonValue) JsonArray {
	if v == nil {
		v = Null()
	}
	return append(a, v)
}
func (a JsonArray) Len() int {
	return len(a)
}
