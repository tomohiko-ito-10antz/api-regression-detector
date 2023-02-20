package assert

import "testing"

func anyEqual(actual any, expect any) bool {
	return actual == expect
}
func Equal[T any](t *testing.T, actual T, expect T) {
	t.Helper()

	if !anyEqual(actual, expect) {
		t.Errorf("ASSERT EQUAL\n  expect: %v:%T\n  actual: %v:%T", expect, expect, actual, actual)
	}
}

func NotEqual(t *testing.T, actual any, expect any) {
	t.Helper()

	if anyEqual(actual, expect) {
		t.Errorf("ASSERT NOT EQUAL\n  expect: %v:%T\n  actual: %v:%T", expect, expect, actual, actual)
	}
}
