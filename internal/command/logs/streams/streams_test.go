package streams

import (
	"context"
	"testing"

	"github.com/skpr/api/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mockLogsClient is a stub pb.LogsClient that only implements the two list RPCs.
type mockLogsClient struct {
	pb.LogsClient

	v2Resp *pb.LogListStreamsV2Response
	v2Err  error

	v1Resp            *pb.LogListStreamsResponse
	listStreamsCalled bool
}

func (m *mockLogsClient) ListStreamsV2(_ context.Context, _ *pb.LogListStreamsV2Request, _ ...grpc.CallOption) (*pb.LogListStreamsV2Response, error) {
	return m.v2Resp, m.v2Err
}

func (m *mockLogsClient) ListStreams(_ context.Context, _ *pb.LogListStreamsRequest, _ ...grpc.CallOption) (*pb.LogListStreamsResponse, error) {
	m.listStreamsCalled = true
	return m.v1Resp, nil
}

func TestListUsesV2WhenAvailable(t *testing.T) {
	mock := &mockLogsClient{
		v2Resp: &pb.LogListStreamsV2Response{
			Streams: []*pb.LogStream{
				{Name: "nginx", Types: []pb.LogStreamType{pb.LogStreamType_Tail}},
			},
			Defaults: &pb.LogStreamDefaults{Tail: []string{"nginx"}},
		},
	}

	resp, err := List(context.Background(), mock, "dev")
	require.NoError(t, err)
	assert.False(t, mock.listStreamsCalled, "should not fall back when V2 succeeds")
	assert.Equal(t, "nginx", resp.Streams[0].Name)
	assert.Equal(t, []string{"nginx"}, resp.Defaults.Tail)
}

func TestListFallsBackWhenV2Unimplemented(t *testing.T) {
	mock := &mockLogsClient{
		v2Err: status.Error(codes.Unimplemented, "method ListStreamsV2 not implemented"),
		v1Resp: &pb.LogListStreamsResponse{
			Streams: []string{"nginx", "fpm"},
			Default: "nginx",
		},
	}

	resp, err := List(context.Background(), mock, "dev")
	require.NoError(t, err)
	assert.True(t, mock.listStreamsCalled, "should fall back when V2 is unimplemented")

	require.Len(t, resp.Streams, 2)
	assert.Equal(t, "nginx", resp.Streams[0].Name)
	assert.Equal(t, "fpm", resp.Streams[1].Name)
	assert.Empty(t, resp.Streams[0].Types, "legacy streams have no capability metadata")
	assert.Equal(t, []string{"nginx"}, resp.Defaults.Tail)
}

func TestListPropagatesOtherErrors(t *testing.T) {
	mock := &mockLogsClient{
		v2Err: status.Error(codes.PermissionDenied, "nope"),
	}

	_, err := List(context.Background(), mock, "dev")
	require.Error(t, err)
	assert.False(t, mock.listStreamsCalled, "should not fall back on non-Unimplemented errors")
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
}
