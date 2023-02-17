package errors

import (
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

var (
	BadIndexAccess error = errors.New("BadIndexAccess")
	BadKeyAccess   error = errors.New("BadKeyAccess")
	Unsupported    error = errors.New("Unsupported")
	Unexpected     error = errors.New("Unsupported")
	BadConversion  error = errors.New("BadConversion")
	BadArgs        error = errors.New("BadArgs")
	BadState       error = errors.New("BadState")
	IOFailure      error = errors.New("IOFailure")
	DBFailure      error = errors.New("DBFailure")
	BadJSON        error = errors.New("BadJSON")
	HTTPFailure    error = errors.New("HTTPFailure")
	GRPCFailure    error = errors.New("GRPCFailure")
)

var (
	Wrap = errors.Wrapf
	Join = multierr.Combine // stderrors.Join
)
