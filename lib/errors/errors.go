package errors

import (
	stderrors "errors"

	"github.com/pkg/errors"
)

var (
	BadIndexAccess error = errors.New("BadIndexAccess")
	BadKeyAccess   error = errors.New("BadKeyAccess")
	BadConversion  error = errors.New("BadConversion")
	BadArgs        error = errors.New("BadArgs")
	BadState       error = errors.New("BadState")
	IOFailure      error = errors.New("IOFailure")
	DBFailure      error = errors.New("DBFailure")
	BadJSON        error = errors.New("BadJSON")
)

var (
	Wrap = errors.Wrapf
	Join = stderrors.Join
)
