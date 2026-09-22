package get

import (
	"bytes"
	"context"
	"testing"

	"github.com/skpr/api/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
)

type fakeAPI struct {
	threshold *durationpb.Duration
}

func (api *fakeAPI) GetThreshold(context.Context, *pb.TraceGetThresholdRequest, ...grpc.CallOption) (*pb.TraceGetThresholdResponse, error) {
	return &pb.TraceGetThresholdResponse{Threshold: api.threshold}, nil
}

func connect(a api) connectFunc {
	return func(ctx context.Context) (context.Context, api, error) {
		return ctx, a, nil
	}
}

// A zero threshold traces every call, so it is a meaningful value rather than a
// stand in for one the environment never reported.
func TestGetThresholdZero(t *testing.T) {
	cmd := Command{Environment: "dev"}

	var b bytes.Buffer

	require.NoError(t, cmd.run(context.Background(), connect(&fakeAPI{threshold: durationpb.New(0)}), &b))
	assert.Equal(t, "0s\n", b.String())
}

func TestGetThresholdMissing(t *testing.T) {
	cmd := Command{Environment: "dev"}

	var b bytes.Buffer

	err := cmd.run(context.Background(), connect(&fakeAPI{}), &b)
	require.EqualError(t, err, `environment "dev" did not report a tracing threshold`)
	assert.Empty(t, b.String())
}
