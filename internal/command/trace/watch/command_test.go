package watch

import (
	"testing"
	"time"

	"github.com/skpr/api/pb"
	compasstrace "github.com/skpr/compass/pkg/trace"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTraceFromProto(t *testing.T) {
	startTime := time.Date(2026, time.September, 4, 12, 0, 0, 0, time.UTC)
	endTime := startTime.Add(1500 * time.Millisecond)

	actual := traceFromProto(&pb.Trace{
		Metadata: &pb.TraceMetadata{
			RequestId: "request-123",
			StartTime: timestamppb.New(startTime),
			EndTime:   timestamppb.New(endTime),
			Source:    pb.TraceSource_TRACE_SOURCE_HTTP,
			Runtime:   pb.TraceRuntime_TRACE_RUNTIME_PHP,
			Http: &pb.TraceMetadataHTTP{
				Method: "GET",
				Uri:    "/articles/123",
			},
		},
		ResourceUtilisation: &pb.TraceResourceUtilisation{MaxMemory: 4096},
		Spans: []*pb.TraceSpan{
			{
				Name:    "Drupal\\Core\\Kernel::handle",
				Offset:  durationpb.New(25 * time.Millisecond),
				Elapsed: durationpb.New(750 * time.Millisecond),
				Total:   durationpb.New(900 * time.Millisecond),
				Calls:   5,
				Memory:  1024,
			},
			nil,
		},
		Calls:        7,
		CallsDropped: 2,
		Drupal: &pb.TraceDrupal{
			CacheEvents: []*pb.TraceDrupalCacheEvent{
				{
					Origin:     pb.TraceDrupalCacheOrigin_TRACE_DRUPAL_CACHE_ORIGIN_OBJECT,
					Caller:     "Drupal\\node\\Entity\\Node::getCacheTags",
					ObjectType: "node",
					MaxAge:     3600,
					Tags:       []string{"node:123"},
					Contexts:   []string{"url.path"},
					Offset:     durationpb.New(100 * time.Millisecond),
					Calls:      3,
				},
				nil,
			},
			CacheEventsDropped: 1,
		},
	})

	expected := compasstrace.Trace{
		Metadata: compasstrace.Metadata{
			Source:    compasstrace.SourceHTTP,
			Runtime:   compasstrace.RuntimePHP,
			ID:        "request-123",
			HTTP:      compasstrace.MetadataHTTP{Method: "GET", URI: "/articles/123"},
			StartTime: startTime,
			EndTime:   endTime,
		},
		ResourceUtilisation: compasstrace.ResourceUtilisation{MaxMemory: 4096},
		Spans: []compasstrace.Span{
			{
				Name:    "Drupal\\Core\\Kernel::handle",
				Offset:  25 * time.Millisecond,
				Elapsed: 750 * time.Millisecond,
				Total:   900 * time.Millisecond,
				Calls:   5,
				Memory:  1024,
			},
		},
		Calls:        7,
		CallsDropped: 2,
		Drupal: &compasstrace.Drupal{
			CacheEvents: []compasstrace.CacheEvent{
				{
					Origin:     compasstrace.CacheOriginObject,
					Caller:     "Drupal\\node\\Entity\\Node::getCacheTags",
					ObjectType: "node",
					MaxAge:     3600,
					Tags:       []string{"node:123"},
					Contexts:   []string{"url.path"},
					Offset:     100 * time.Millisecond,
					Calls:      3,
				},
			},
			CacheEventsDropped: 1,
		},
	}

	assert.Equal(t, expected, actual)
}

func TestTraceFromProtoCLI(t *testing.T) {
	actual := traceFromProto(&pb.Trace{
		Metadata: &pb.TraceMetadata{
			Source:  pb.TraceSource_TRACE_SOURCE_CLI,
			Runtime: pb.TraceRuntime_TRACE_RUNTIME_NODE,
			Cli:     &pb.TraceMetadataCLI{Command: "drush cr"},
		},
	})

	assert.Equal(t, compasstrace.SourceCLI, actual.Metadata.Source)
	assert.Equal(t, compasstrace.RuntimeNode, actual.Metadata.Runtime)
	assert.Equal(t, "drush cr", actual.Metadata.CLI.Command)
}

func TestTraceFromProtoNil(t *testing.T) {
	assert.Equal(t, compasstrace.Trace{}, traceFromProto(nil))
}
