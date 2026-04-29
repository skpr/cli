// Package logs contains shared helpers for log subcommands.
package logs

import (
	"fmt"
	"time"

	"github.com/skpr/api/pb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	skprtime "github.com/skpr/cli/internal/time"
)

// FilterParams holds the user-supplied filter inputs for query/summarise.
type FilterParams struct {
	Environment string
	Streams     []string
	Since       time.Duration
	From        string
	To          string
	Contains    []string
	Exclude     []string
}

// BuildFilter converts user inputs into a *pb.LogFilter.
func BuildFilter(p FilterParams) (*pb.LogFilter, error) {
	filter := &pb.LogFilter{
		Environment: p.Environment,
		Streams:     p.Streams,
	}

	hasFrom := p.From != ""
	hasTo := p.To != ""

	switch {
	case hasFrom != hasTo:
		return nil, fmt.Errorf("--from and --to must be provided together")
	case hasFrom && hasTo:
		from, err := skprtime.ParseString(p.From)
		if err != nil {
			return nil, fmt.Errorf("failed to parse --from: %w", err)
		}

		to, err := skprtime.ParseString(p.To)
		if err != nil {
			return nil, fmt.Errorf("failed to parse --to: %w", err)
		}

		filter.Window = &pb.LogFilter_TimeRange{
			TimeRange: &pb.LogTimeRange{
				From: timestamppb.New(from),
				To:   timestamppb.New(to),
			},
		}
	default:
		filter.Window = &pb.LogFilter_Timeframe{
			Timeframe: durationpb.New(p.Since),
		}
	}

	for _, v := range p.Contains {
		filter.Contains = append(filter.Contains, &pb.LogContainsFilter{
			Value:   v,
			Exclude: false,
		})
	}

	for _, v := range p.Exclude {
		filter.Contains = append(filter.Contains, &pb.LogContainsFilter{
			Value:   v,
			Exclude: true,
		})
	}

	return filter, nil
}
