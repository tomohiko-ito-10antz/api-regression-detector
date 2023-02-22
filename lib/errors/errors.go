package errors

import (
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

var (
	Wrap = errors.Wrapf
	Join = multierr.Combine
	Is   = errors.Is
	New  = errors.Errorf
)

func Assert(condition bool, format string, args ...any) {
	if !condition {
		panic(Unexpected.New(format, args...))
	}
}

func Unreachable[T any](format string, args ...any) T {
	panic(Unexpected.New(format, args...))
}

func Unreachable2[T1, T2 any](format string, args ...any) (T1, T2) {
	panic(Unexpected.New(format, args...))
}

func Unreachable3[T1, T2, T3 any](format string, args ...any) (T1, T2, T3) {
	panic(Unexpected.New(format, args...))
}

func Unreachable4[T1, T2, T3, T4 any](format string, args ...any) (T1, T2, T3, T4) {
	panic(Unexpected.New(format, args...))
}
