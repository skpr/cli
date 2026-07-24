package filter

import (
	"fmt"
	"time"

	"github.com/skpr/api/pb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	skprtime "github.com/skpr/cli/internal/time"
)

// Options holds the shared filter inputs used by bounded log operations such as
// query and summarise.
type Options struct {
	Environment string
	Streams     []string
	Timeframe   time.Duration
	From        string
	To          string
	Contains    []string
	Exclude     []string
}

// Build converts the options into a pb.LogFilter. An absolute time range
// (--from/--to) takes precedence over the relative --timeframe.
func (o Options) Build() (*pb.LogFilter, error) {
	filter := &pb.LogFilter{
		Environment: o.Environment,
		Streams:     o.Streams,
	}

	// Build the substring filters. An event must satisfy every entry.
	for _, value := range o.Contains {
		filter.Contains = append(filter.Contains, &pb.LogContainsFilter{
			Value: value,
		})
	}

	for _, value := range o.Exclude {
		filter.Contains = append(filter.Contains, &pb.LogContainsFilter{
			Value:   value,
			Exclude: true,
		})
	}

	if o.From != "" || o.To != "" {
		window, err := buildTimeRange(o.From, o.To)
		if err != nil {
			return nil, err
		}

		filter.Window = window
	} else {
		filter.Window = &pb.LogFilter_Timeframe{
			Timeframe: durationpb.New(o.Timeframe),
		}
	}

	return filter, nil
}

// buildTimeRange constructs an absolute time range window from the provided
// --from and --to values. Either bound may be omitted.
func buildTimeRange(from, to string) (*pb.LogFilter_TimeRange, error) {
	timeRange := &pb.LogTimeRange{}

	if from != "" {
		t, err := skprtime.ParseString(from)
		if err != nil {
			return nil, fmt.Errorf("failed to parse --from: %w", err)
		}

		timeRange.From = timestamppb.New(t)
	}

	if to != "" {
		t, err := skprtime.ParseString(to)
		if err != nil {
			return nil, fmt.Errorf("failed to parse --to: %w", err)
		}

		timeRange.To = timestamppb.New(t)
	}

	return &pb.LogFilter_TimeRange{TimeRange: timeRange}, nil
}

// AddFlags registers the shared filter flags on a command's flag set.
func AddFlags(flags flagSet, o *Options) {
	flags.DurationVar(&o.Timeframe, "timeframe", time.Hour, "Window duration relative to now (ignored when --from or --to is set)")
	flags.StringVar(&o.From, "from", "", "Start of an absolute time range (e.g. \"2 days ago\", RFC3339 timestamp)")
	flags.StringVar(&o.To, "to", "", "End of an absolute time range (e.g. \"now\", RFC3339 timestamp)")
	flags.StringArrayVar(&o.Contains, "contains", nil, "Only include events whose message contains this substring (repeatable)")
	flags.StringArrayVar(&o.Exclude, "exclude", nil, "Exclude events whose message contains this substring (repeatable)")
}

// flagSet is the subset of *pflag.FlagSet used to register filter flags.
type flagSet interface {
	DurationVar(p *time.Duration, name string, value time.Duration, usage string)
	StringVar(p *string, name, value, usage string)
	StringArrayVar(p *[]string, name string, value []string, usage string)
}
