package create

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/skpr/api/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type Mock struct {
	pb.BackupClient
}

func (m *Mock) List(ctx context.Context, in *pb.BackupListRequest, opts ...grpc.CallOption) (*pb.BackupListResponse, error) {
	return &pb.BackupListResponse{
		List: []*pb.BackupStatus{
			{
				Name:  "completed",
				Phase: pb.BackupStatus_Completed,
			},
			{
				Name:  "failed",
				Phase: pb.BackupStatus_Failed,
			},
			{
				Name:  "unknown",
				Phase: pb.BackupStatus_Unknown,
			},
		},
	}, nil
}

func TestGetBackup(t *testing.T) {
	restore := getBackup("foo", []*pb.BackupStatus{
		{
			Name: "bar",
		},
	})
	assert.Nil(t, restore)

	restore = getBackup("foo", []*pb.BackupStatus{
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
	assert.ErrorContains(t, err, "the backup failed")
	assert.True(t, exit)
}

func TestWaitUnknown(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	mock := &Mock{}

	f := wait(context.TODO(), logger, mock, "dev", "unknown")

	exit, err := f()
	assert.ErrorContains(t, err, "the backup failed for an unknown reason")
	assert.True(t, exit)
}
