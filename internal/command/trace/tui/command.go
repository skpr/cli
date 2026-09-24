package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/skpr/api/pb"
	"github.com/skpr/compass/pkg/app"
	applogger "github.com/skpr/compass/pkg/app/logger"
	compasstrace "github.com/skpr/compass/pkg/trace"
	"golang.org/x/sync/errgroup"
)

// Command which streams traces for an environment into the Compass app.
type Command struct {
	Environment string
}

// Run the command.
func (cmd *Command) Run(ctx context.Context) error {
	ctx, api, err := cmd.preflight(ctx, connectAPI)
	if err != nil {
		return err
	}

	model := app.NewModel("", app.Options{
		MaxTraces: app.DefaultMaxTraces,
		MaxLogs:   app.DefaultMaxLogs,
		MaxBytes:  app.DefaultMaxBytes,
	})

	p := tea.NewProgram(model, tea.WithAltScreen())

	logger, err := applogger.New(p)
	if err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eg := errgroup.Group{}

	// Start the collector.
	eg.Go(func() error {
		return collectTraces(ctx, api, cmd.Environment, p, logger)
	})

	// Start the application.
	eg.Go(func() error {
		_, err := p.Run()
		cancel()
		if err != nil {
			return fmt.Errorf("failed to run program: %w", err)
		}

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

	dst.Spans = make([]compasstrace.Span, 0, len(src.GetSpans()))
	for _, span := range src.GetSpans() {
		if span == nil {
			continue
		}

		converted := compasstrace.Span{
			Name:   span.GetName(),
			Calls:  span.GetCalls(),
			Memory: span.GetMemory(),
		}
		if offset := span.GetOffset(); offset != nil {
			converted.Offset = offset.AsDuration()
		}
		if elapsed := span.GetElapsed(); elapsed != nil {
			converted.Elapsed = elapsed.AsDuration()
		}
		if total := span.GetTotal(); total != nil {
			converted.Total = total.AsDuration()
		}

		dst.Spans = append(dst.Spans, converted)
	}
	dst.Calls = src.GetCalls()
	dst.CallsDropped = src.GetCallsDropped()

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
