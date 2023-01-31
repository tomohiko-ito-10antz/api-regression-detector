package assert

import "testing"

func Equal(t *testing.T, actual any, expect any) {
	if actual != expect {
		t.Errorf(`expect: %v, actual: %v`, expect, actual)
	}
}
func NotEqual(t *testing.T, actual any, expect any) {
	if actual == expect {
		t.Errorf(`expect: not %v, actual: %v`, expect, actual)
	}
}
