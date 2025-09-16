package create

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/skpr/api/pb"
)

type Mock struct {
	pb.RestoreClient
}

func (m *Mock) List(ctx context.Context, in *pb.RestoreListRequest, opts ...grpc.CallOption) (*pb.RestoreListResponse, error) {
	return &pb.RestoreListResponse{
		List: []*pb.RestoreStatus{
			{
				Name:  "completed",
				Phase: pb.RestoreStatus_Completed,
			},
			{
				Name:  "failed",
				Phase: pb.RestoreStatus_Failed,
			},
			{
				Name:  "unknown",
				Phase: pb.RestoreStatus_Unknown,
			},
		},
	}, nil
}

func TestGetRestore(t *testing.T) {
	restore := getRestore("foo", []*pb.RestoreStatus{
		{
			Name: "bar",
		},
	})
	assert.Nil(t, restore)

	restore = getRestore("foo", []*pb.RestoreStatus{
		{
			Name: "foo",
		},
	})
	assert.NotNil(t, restore)
}

func TestWaitNotExist(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	mock := &Mock{}

	f := wait(context.TODO(), logger, mock, "dev", "not-found")

	exit, err := f()
	assert.ErrorContains(t, err, "backup does not exist")
	assert.True(t, exit)
}

func TestWaitCompleted(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	mock := &Mock{}

	f := wait(context.TODO(), logger, mock, "dev", "completed")

	exit, err := f()
	assert.NoError(t, err)
	assert.True(t, exit)
}

func TestWaitFailed(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	mock := &Mock{}

	f := wait(context.TODO(), logger, mock, "dev", "failed")

	exit, err := f()
	assert.ErrorContains(t, err, "the restore failed")
	assert.True(t, exit)
}

func TestWaitUnknown(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	mock := &Mock{}

	f := wait(context.TODO(), logger, mock, "dev", "unknown")

	exit, err := f()
	assert.ErrorContains(t, err, "the restore failed for an unknown reason")
	assert.True(t, exit)
}
