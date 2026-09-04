package trace

import (
	"context"
	"fmt"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skpr/api/pb"
	"github.com/skpr/compass/pkg/app"
	"github.com/skpr/compass/pkg/app/events"
	applogger "github.com/skpr/compass/pkg/app/logger"
	compasstrace "github.com/skpr/compass/pkg/trace"
	"golang.org/x/sync/errgroup"

	"github.com/skpr/cli/internal/client"
)

// Command to trace environments.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	p := tea.NewProgram(app.NewModel("", app.DefaultMaxTraces, app.DefaultMaxLogs), tea.WithAltScreen())

	logger, err := applogger.New(p)
	if err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)

	eg := errgroup.Group{}

	// Start the collector.
	eg.Go(func() error {
		logger.Info("Connecting to Skpr API...")

		ctx, client, err := client.New(ctx)
		if err != nil {
			return err
		}

		stream, err := client.Trace().StreamTraces(ctx, &pb.StreamTracesRequest{
			Environment: cmd.Environment,
		})
		if err != nil {
			return err
		}

		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				resp, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					return fmt.Errorf("streaming trace failed: %w", err)
				}

				for _, t := range resp.Traces {
					converted := traceFromProto(t)
					p.Send(events.Trace{
						IngestionTime: converted.Metadata.StartTime,
						Trace:         converted,
					})
				}
			}
		}
	})

	// Start the application.
	eg.Go(func() error {
		_, err := p.Run()
		if err != nil {
			return fmt.Errorf("failed to run program: %w", err)
		}

		cancel()

		return nil
	})

	return eg.Wait()
}

func traceFromProto(src *pb.Trace) compasstrace.Trace {
	var dst compasstrace.Trace
	if src == nil {
		return dst
	}

	if metadata := src.GetMetadata(); metadata != nil {
		dst.Metadata.ID = metadata.GetRequestId()
		dst.Metadata.Source = traceSourceFromProto(metadata.GetSource())
		dst.Metadata.Runtime = traceRuntimeFromProto(metadata.GetRuntime())

		if startTime := metadata.GetStartTime(); startTime != nil {
			dst.Metadata.StartTime = startTime.AsTime()
		}
		if endTime := metadata.GetEndTime(); endTime != nil {
			dst.Metadata.EndTime = endTime.AsTime()
		}
		if httpMetadata := metadata.GetHttp(); httpMetadata != nil {
			dst.Metadata.HTTP = compasstrace.MetadataHTTP{
				Method: httpMetadata.GetMethod(),
				URI:    httpMetadata.GetUri(),
			}
		}
		if cliMetadata := metadata.GetCli(); cliMetadata != nil {
			dst.Metadata.CLI.Command = cliMetadata.GetCommand()
		}
	}

	if resources := src.GetResourceUtilisation(); resources != nil {
		dst.ResourceUtilisation.MaxMemory = resources.GetMaxMemory()
	}

	dst.FunctionCalls = make([]compasstrace.FunctionCall, 0, len(src.GetFunctionCalls()))
	for _, functionCall := range src.GetFunctionCalls() {
		if functionCall == nil {
			continue
		}

		converted := compasstrace.FunctionCall{
			Name:   functionCall.GetName(),
			Memory: functionCall.GetMemory(),
		}
		if offset := functionCall.GetOffset(); offset != nil {
			converted.Offset = offset.AsDuration()
		}
		if elapsed := functionCall.GetElapsed(); elapsed != nil {
			converted.Elapsed = elapsed.AsDuration()
		}

		dst.FunctionCalls = append(dst.FunctionCalls, converted)
	}
	dst.FunctionCallsDropped = int(src.GetFunctionCallsDropped())

	if drupal := src.GetDrupal(); drupal != nil {
		dst.Drupal = &compasstrace.Drupal{
			CacheEvents:        make([]compasstrace.CacheEvent, 0, len(drupal.GetCacheEvents())),
			CacheEventsDropped: int(drupal.GetCacheEventsDropped()),
		}

		for _, cacheEvent := range drupal.GetCacheEvents() {
			if cacheEvent == nil {
				continue
			}

			converted := compasstrace.CacheEvent{
				Origin:     traceDrupalCacheOriginFromProto(cacheEvent.GetOrigin()),
				Caller:     cacheEvent.GetCaller(),
				ObjectType: cacheEvent.GetObjectType(),
				MaxAge:     cacheEvent.GetMaxAge(),
				Tags:       append([]string(nil), cacheEvent.GetTags()...),
				Contexts:   append([]string(nil), cacheEvent.GetContexts()...),
				Calls:      cacheEvent.GetCalls(),
			}
			if offset := cacheEvent.GetOffset(); offset != nil {
				converted.Offset = offset.AsDuration()
			}

			dst.Drupal.CacheEvents = append(dst.Drupal.CacheEvents, converted)
		}
	}

	return dst
}

func traceSourceFromProto(src pb.TraceSource) compasstrace.Source {
	switch src {
	case pb.TraceSource_TRACE_SOURCE_HTTP:
		return compasstrace.SourceHTTP
	case pb.TraceSource_TRACE_SOURCE_CLI:
		return compasstrace.SourceCLI
	default:
		return ""
	}
}

func traceRuntimeFromProto(src pb.TraceRuntime) compasstrace.Runtime {
	switch src {
	case pb.TraceRuntime_TRACE_RUNTIME_PHP:
		return compasstrace.RuntimePHP
	case pb.TraceRuntime_TRACE_RUNTIME_NODE:
		return compasstrace.RuntimeNode
	default:
		return ""
	}
}

func traceDrupalCacheOriginFromProto(src pb.TraceDrupalCacheOrigin) compasstrace.CacheOrigin {
	switch src {
	case pb.TraceDrupalCacheOrigin_TRACE_DRUPAL_CACHE_ORIGIN_RENDER_ARRAY:
		return compasstrace.CacheOriginRenderArray
	case pb.TraceDrupalCacheOrigin_TRACE_DRUPAL_CACHE_ORIGIN_OBJECT:
		return compasstrace.CacheOriginObject
	default:
		return ""
	}
}
