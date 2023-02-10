// Command grpcurl makes gRPC requests (a la cURL, but HTTP/2). It can use a supplied descriptor
// file, protobuf sources, or service reflection to translate JSON or text request data into the
// appropriate protobuf messages and vice versa for presenting the response contents.
package grpcurl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	// Register gzip compressor so compressed responses will work
	_ "google.golang.org/grpc/encoding/gzip"
	// Register xds so xds and xds-experimental resolver schemes work
	_ "google.golang.org/grpc/xds"

	"github.com/fullstorydev/grpcurl"
)

type Request struct {
	Headers map[string][]string
	Body    string
}

type Response struct {
	Headers map[string][]string
	Body    string
}

func Grpcurl(
	ctx context.Context,
	req Request,
	target string,
	fullMethod string,
	options Options,
) (*Response, error) {
	if err := options.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid option")
	}

	if target == "" {
		return nil, errors.Wrap(errors.BadArgs, "no host:port specified")
	}

	verbosityLevel := 0
	if options.Verbose {
		verbosityLevel = 1
	}

	if options.MaxTime > 0 {
		timeout := time.Duration(options.MaxTime * float64(time.Second))
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	headers := []string{}
	for key, vals := range req.Headers {
		for _, val := range vals {
			headers = append(headers, fmt.Sprintf(`%s: %s`, key, val))
		}
	}

	var cc *grpc.ClientConn
	md := grpcurl.MetadataFromHeaders(headers)
	refCtx := metadata.NewOutgoingContext(ctx, md)

	cc, err := dialClientConn(ctx, target, options)
	if err != nil {
		return nil, errors.Wrap(err, "fail to dial client connection")
	}

	refClient := grpcreflect.NewClientV1Alpha(refCtx, reflectpb.NewServerReflectionClient(cc))
	descSource := grpcurl.DescriptorSourceFromServer(ctx, refClient)

	// arrange for the RPCs to be cleanly shutdown
	reset := func() {
		if refClient != nil {
			refClient.Reset()
			refClient = nil
		}
		if cc != nil {
			cc.Close()
			cc = nil
		}
	}
	defer reset()

	// Invoke an RPC
	cc, err = dialClientConn(ctx, target, options)
	if err != nil {
		return nil, errors.Wrap(err, "fail to dial client connection")
	}

	// if not verbose output, then also include record delimiters
	// between each message, so output could potentially be piped
	// to another grpcurl process
	includeSeparators := verbosityLevel == 0

	rf, formatter, err := grpcurl.RequestParserAndFormatter(
		options.Format,
		descSource,
		strings.NewReader(req.Body),
		grpcurl.FormatOptions{
			EmitJSONDefaultFields: options.EmitDefaults,
			IncludeTextSeparator:  includeSeparators,
			AllowUnknownFields:    options.AllowUnknownFields,
		})
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCURLFailure),
			"failed to construct request parser and formatter for %s", options.Format)
	}

	handler := &Handler{
		Formatter: formatter,
	}

	err = grpcurl.InvokeRPC(ctx, descSource, cc, fullMethod, headers, handler, rf.Next)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCURLFailure),
			"failed to invoke method %s", fullMethod)
	}

	if handler.Error != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCURLFailure),
			"failed to parse response as JSON", fullMethod)
	}

	return &handler.Response, nil
}

func dialClientConn(ctx context.Context, target string, options Options) (*grpc.ClientConn, error) {
	dialTime := 10 * time.Second
	if options.ConnectTimeout > 0 {
		dialTime = time.Duration(options.ConnectTimeout * float64(time.Second))
	}
	ctx, cancel := context.WithTimeout(ctx, dialTime)
	defer cancel()
	var dialOptions []grpc.DialOption
	if options.KeepaliveTime > 0 {
		timeout := time.Duration(options.KeepaliveTime * float64(time.Second))
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    timeout,
			Timeout: timeout,
		}))
	}
	if options.MaxMsgSz > 0 {
		dialOptions = append(dialOptions, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(options.MaxMsgSz)))
	}
	var creds credentials.TransportCredentials
	if !options.Plaintext {
		tlsConf, err := grpcurl.ClientTLSConfig(options.Insecure, options.Cacert, options.Cert, options.Key)
		if err != nil {
			return nil, errors.Wrap(
				errors.Join(err, errors.GRPCURLFailure),
				"fail to create TLS config")
		}

		creds = credentials.NewTLS(tlsConf)

		if options.Authority != "" {
			dialOptions = append(dialOptions, grpc.WithAuthority(options.Authority))
		}
	} else if options.Authority != "" {
		dialOptions = append(dialOptions, grpc.WithAuthority(options.Authority))
	}

	grpcurlUA := "grpcurl"
	if options.UserAgent != "" {
		grpcurlUA = options.UserAgent + " " + grpcurlUA
	}
	dialOptions = append(dialOptions, grpc.WithUserAgent(grpcurlUA))

	network := "tcp"
	cc, err := grpcurl.BlockingDial(ctx, network, target, creds, dialOptions...)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCURLFailure),
			"Failed to dial target host %s", target)
	}
	return cc, nil
}
