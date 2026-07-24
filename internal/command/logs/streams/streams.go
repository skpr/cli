package streams

import (
	"context"
	"fmt"

	"github.com/skpr/api/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// List returns the streams available for an environment using the ListStreamsV2
// RPC. Older servers may not implement ListStreamsV2, in which case this falls
// back to the original ListStreams RPC and adapts the response to the V2 shape.
func List(ctx context.Context, client pb.LogsClient, environment string, types ...pb.LogStreamType) (*pb.LogListStreamsV2Response, error) {
	resp, err := client.ListStreamsV2(ctx, &pb.LogListStreamsV2Request{
		Environment: environment,
		Types:       types,
	})
	if err == nil {
		return resp, nil
	}

	// Any error other than the RPC being missing is a genuine failure.
	if status.Code(err) != codes.Unimplemented {
		return nil, err
	}

	// Fall back to the original ListStreams RPC for older servers. It has no
	// concept of stream capabilities, so every stream is returned without types
	// and the single default is treated as the default for tailing.
	legacy, err := client.ListStreams(ctx, &pb.LogListStreamsRequest{
		Environment: environment,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list streams: %w", err)
	}

	out := &pb.LogListStreamsV2Response{
		Defaults: &pb.LogStreamDefaults{},
	}

	for _, name := range legacy.Streams {
		out.Streams = append(out.Streams, &pb.LogStream{
			Name: name,
		})
	}

	if legacy.Default != "" {
		out.Defaults.Tail = []string{legacy.Default}
	}

	return out, nil
}
