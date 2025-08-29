// Package log is an implementation of https://dusted.codes/creating-a-pretty-console-logger-using-gos-slog-package
package log

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"

	"github.com/jwalton/gchalk"

	"github.com/skpr/cli/internal/color"
)

const (
	// Formats the time for our log event.
	timeFormat = "[15:04:05.000]"
)

// NewHandler for the Skpr CLI slog implementation.
func NewHandler(w io.Writer, opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	b := &bytes.Buffer{}

	return &Handler{
		w: w,
		b: b,
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
		m: &sync.Mutex{},
	}
}

// Handler for the Skpr CLI slog implementation.
type Handler struct {
	w io.Writer
	h slog.Handler
	b *bytes.Buffer
	m *sync.Mutex
}

// Enabled wraps the slog implementation.
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

// WithAttrs wraps the slog implementation.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m}
}

// WithGroup wraps the slog implementation.
func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name), b: h.b, m: h.m}
}

// Handle the log event.
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = gchalk.WithHex(color.HexBlue).Sprintf("%s", level)
	case slog.LevelInfo:
		level = gchalk.WithHex(color.HexBlue).Sprintf("%s", level)
	case slog.LevelWarn:
		level = gchalk.WithHex(color.HexYellow).Sprintf("%s", level)
	case slog.LevelError:
		level = gchalk.WithHex(color.HexRed).Sprintf("%s", level)
	}

	// @todo, Add attribute handling.
	_, err := fmt.Fprintln(h.w, r.Time.Format(timeFormat), level, gchalk.Bold(r.Message))
	if err != nil {
		return fmt.Errorf("failed to format log message: %w", err)
	}

	return nil
}

// Modify the default slog configuration.
func suppressDefaults(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}
