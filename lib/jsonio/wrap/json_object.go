package wrap

import "fmt"

type JsonObject map[string]*JsonValue

func (o JsonObject) Empty() bool {
	return len(o) == 0
}

func (o JsonObject) Has(k string) bool {
	_, ok := o[k]
	return ok
}

func (o JsonObject) Get(k string) *JsonValue {
	if !o.Has(k) {
		panic(fmt.Errorf("key %s not found", k))
	}
	return o[k]
}

func (o JsonObject) Set(k string, v *JsonValue) {
	if v == nil {
		v = Null()
	}
	o[k] = v
}

func (o JsonObject) Len() int {
	return len(o)
}

func (o JsonObject) Keys() []string {
	keys := []string{}
	for k := range o {
		keys = append(keys, k)
	}
	return keys
}
