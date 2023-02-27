package errors

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

const (
	BadIndexAccess Tag = "BadIndexAccess"
	BadKeyAccess   Tag = "BadKeyAccess"
	Unsupported    Tag = "Unsupported"
	Unexpected     Tag = "Unexpected"
	BadConversion  Tag = "BadConversion"
	BadArgs        Tag = "BadArgs"
	BadState       Tag = "BadState"
	IOFailure      Tag = "IOFailure"
	DBFailure      Tag = "DBFailure"
	BadProtoBuf    Tag = "BadProtoBuf"
	BadJSON        Tag = "BadJSON"
	HTTPFailure    Tag = "HTTPFailure"
	GRPCFailure    Tag = "GRPCFailure"
)

type withTag struct {
	tag   Tag
	cause error
}

type Tag string

func (tag Tag) New(format string, args ...any) error {
	return errors.WithStack(tag.Err(errors.Errorf(format, args...)))
}
func (tag Tag) Err(err error) error {
	return WithTag(err, tag)
}
func WithTag(err error, tags ...Tag) error {
	if err == nil {
		return nil
	}
	if len(tags) == 0 {
		return err
	}
	for _, tag := range tags {
		if !errors.Is(err, &withTag{tag: tag}) {
			err = &withTag{cause: err, tag: tag}
		}
	}
	return errors.WithStack(err)
}
func (e *withTag) Error() string {
	return fmt.Sprintf(`%s: %s`, e.tag, e.cause)
}
func (e *withTag) Unwrap() error {
	return e.cause
}
func (e *withTag) Is(target error) bool {
	if t, ok := target.(*withTag); ok && e.tag == t.tag {
		return true
	}
	return errors.Is(e.Unwrap(), target)
}
func (e *withTag) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", e.Unwrap())
			io.WriteString(s, string(e.tag))
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, e.Error())
	}
}
