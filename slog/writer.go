package slog

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/juju/loggo"
	"github.com/juju/loggo/attrs"
)

type slogWriter struct {
	writer slog.Handler
}

// NewSlogWriter will write out slog severity levels.
func NewSlogWriter(writer slog.Handler) loggo.Writer {
	return &slogWriter{writer: writer}
}

// Write implements Writer.
func (w *slogWriter) Write(entry loggo.Entry) {
	record := slog.NewRecord(
		entry.Timestamp,
		Level(entry.Level),
		entry.Message,
		entry.PC,
	)

	record.AddAttrs(
		slog.String("module", entry.Module),
		slog.String("filename", entry.Filename),
		slog.Int("line", entry.Line),
	)
	if len(entry.Labels) > 0 {
		record.AddAttrs(slog.String("labels", strings.Join(entry.Labels, ",")))
	}
	for _, attr := range entry.Attrs {
		switch a := attr.(type) {
		case attrs.AttrValue[string]:
			record.AddAttrs(slog.String(a.Key(), a.Value()))
		case attrs.AttrValue[int]:
			record.AddAttrs(slog.Int(a.Key(), a.Value()))
		case attrs.AttrValue[int64]:
			record.AddAttrs(slog.Int64(a.Key(), a.Value()))
		case attrs.AttrValue[uint64]:
			record.AddAttrs(slog.Uint64(a.Key(), a.Value()))
		case attrs.AttrValue[float64]:
			record.AddAttrs(slog.Float64(a.Key(), a.Value()))
		case attrs.AttrValue[bool]:
			record.AddAttrs(slog.Bool(a.Key(), a.Value()))
		case attrs.AttrValue[time.Time]:
			record.AddAttrs(slog.Time(a.Key(), a.Value()))
		case attrs.AttrValue[time.Duration]:
			record.AddAttrs(slog.Duration(a.Key(), a.Value()))
		case attrs.AttrValue[any]:
			record.AddAttrs(slog.Any(a.Key(), a.Value()))
		}
	}

	w.writer.Handle(context.Background(), record)
}

// Level function allows levels to be mapped to slog levels. Although,
// slog doesn't explicitly implement all the levels that we require for mapping
// it does allow for custom levels to be added. This is done by using the
// slog.Level type as an int64.
// Reading the documentation https://pkg.go.dev/log/slog#Level explains how
// to insert custom levels.
func Level(level loggo.Level) slog.Level {
	switch level {
	case loggo.TRACE:
		return slog.LevelDebug - 1
	case loggo.DEBUG:
		return slog.LevelDebug
	case loggo.INFO:
		return slog.LevelInfo
	case loggo.WARNING:
		return slog.LevelWarn
	case loggo.ERROR:
		return slog.LevelError
	case loggo.CRITICAL:
		return slog.LevelError + 1
	default:
		panic("unknown level")
	}
}

// DefaultLevel returns the lowest level from the loggo config.
func DefaultLevel(v loggo.Config) slog.Level {
	lowest := loggo.CRITICAL
	for _, level := range v {
		if level < lowest {
			lowest = level
		}
	}
	return Level(lowest)
}
