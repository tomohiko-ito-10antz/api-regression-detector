package grpcurl

import (
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/fullstorydev/grpcurl"
)

type Args struct {
	Target  string
	Method  string
	Options Options
}
type Options struct {
	Plaintext    bool
	EmitDefaults bool

	ConnectTimeout     float64
	KeepaliveTime      float64
	MaxTime            float64
	MaxMsgSz           int
	Insecure           bool
	Cert               string
	Cacert             string
	Key                string
	Format             grpcurl.Format
	Authority          string
	UserAgent          string
	Verbose            bool
	AllowUnknownFields bool
}

func (options Options) Validate() error {
	// Do extra validation on arguments and figure out what user asked us to do.
	if options.ConnectTimeout < 0 {
		return errors.Wrap(errors.BadArgs, "The -connect-timeout argument must not be negative.")
	}
	if options.KeepaliveTime < 0 {
		return errors.Wrap(errors.BadArgs, "The -keepalive-time argument must not be negative.")
	}
	if options.MaxTime < 0 {
		return errors.Wrap(errors.BadArgs, "The -max-time argument must not be negative.")
	}
	if options.MaxMsgSz < 0 {
		return errors.Wrap(errors.BadArgs, "The -max-msg-sz argument must not be negative.")
	}
	if options.Plaintext && options.Insecure {
		return errors.Wrap(errors.BadArgs, "The -plaintext and -insecure arguments are mutually exclusive.")
	}
	if options.Plaintext && options.Cert != "" {
		return errors.Wrap(errors.BadArgs, "The -plaintext and -cert arguments are mutually exclusive.")
	}
	if options.Plaintext && options.Key != "" {
		return errors.Wrap(errors.BadArgs, "The -plaintext and -key arguments are mutually exclusive.")
	}
	if (options.Key == "") != (options.Cert == "") {
		return errors.Wrap(errors.BadArgs, "The -cert and -key arguments must be used together and both be present.")
	}
	if options.Format != "json" && options.Format != "text" {
		return errors.Wrap(errors.BadArgs, "The -format option must be 'json' or 'text'.")
	}
	if options.EmitDefaults && options.Format != "json" {
		return errors.Wrap(errors.BadArgs, "The -emit-defaults is only used when using json format.")
	}
	return nil
}
