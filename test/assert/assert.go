package assert

import "testing"

func Equal(t *testing.T, actual any, expect any) {
	t.Helper()
	if actual != expect {
		t.Errorf("ASSERT EQUAL\n  expect: %v:%T\n  actual: %v:%T", expect, expect, actual, actual)
	}
}

func NotEqual(t *testing.T, actual any, expect any) {
	t.Helper()
	if actual == expect {
		t.Errorf("ASSERT NOT EQUAL\n  expect: %v:%T\n  actual: %v:%T", expect, expect, actual, actual)
	}
}
