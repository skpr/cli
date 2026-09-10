package client

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	skprcredentials "github.com/skpr/cli/internal/client/credentials"
)

// Returned by the platform when our request did not include credentials which
// it could authenticate with.
const messageCredentials = "failed to extract credentials"

// Helper function which maps authentication failures returned by the API into
// an error which tells the developer what to do next.
//
// This covers developers who have never logged in, because we deliberately
// allow commands to run without credentials, as well as sessions which the
// platform has rejected.
func mapAuthError(err error) error {
	if err == nil {
		return nil
	}

	s, ok := status.FromError(err)
	if !ok {
		return err
	}

	// The platform does not always return a status code we can rely on for
	// this eg. Failing to extract credentials is returned as an unknown error.
	if s.Code() == codes.Unauthenticated || strings.Contains(s.Message(), messageCredentials) {
		return fmt.Errorf("%w: %w", skprcredentials.ErrLoginRequired, err)
	}

	return err
}

// Interceptor which maps authentication failures for unary calls.
func authUnaryInterceptor(ctx context.Context, method string, req, reply any, conn *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return mapAuthError(invoker(ctx, method, req, reply, conn, opts...))
}

// Interceptor which maps authentication failures for streaming calls.
func authStreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, conn *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	stream, err := streamer(ctx, desc, conn, method, opts...)
	if err != nil {
		return nil, mapAuthError(err)
	}

	return authClientStream{ClientStream: stream}, nil
}

// Wraps a stream so that authentication failures returned while the stream is
// being consumed are also mapped.
type authClientStream struct {
	grpc.ClientStream
}

// RecvMsg from the stream.
func (s authClientStream) RecvMsg(m any) error {
	return mapAuthError(s.ClientStream.RecvMsg(m))
}

// SendMsg to the stream.
func (s authClientStream) SendMsg(m any) error {
	return mapAuthError(s.ClientStream.SendMsg(m))
}
