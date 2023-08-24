package slog

import (
	"context"
	"log/slog"
	"strings"

	"github.com/juju/loggo"
)

type slogWriter struct {
	writer slog.Handler
}

// NewSlowWriter will write out slog severity levels.
func NewSlowWriter(writer slog.Handler) loggo.Writer {
	return &slogWriter{writer}
}

// Write implements Writer.
func (w *slogWriter) Write(entry loggo.Entry) {
	record := slog.NewRecord(
		entry.Timestamp,
		level(entry.Level),
		entry.Message,
		// TODO (stickupkid): Add a way to log the caller ptr in the
		// loggo.Entry. That way we can push the information directly into
		// the slog.Record.
		0,
	)

	record.AddAttrs(
		slog.String("module", entry.Module),
		slog.String("filename", entry.Filename),
		slog.Int("line", entry.Line),
	)
	if len(entry.Labels) > 0 {
		record.AddAttrs(slog.String("labels", strings.Join(entry.Labels, ",")))
	}

	w.writer.Handle(context.Background(), record)
}

// The level function allows levels to be mapped to slog levels. Although,
// slog doesn't explicitly implement all the levels that we require for mapping
// it does allow for custom levels to be added. This is done by using the
// slog.Level type as an int64.
// Reading the documentation https://pkg.go.dev/log/slog#Level explains how
// to insert custom levels.
func level(level loggo.Level) slog.Level {
	switch level {
	case loggo.TRACE:
		return slog.LevelDebug - 1
	case loggo.DEBUG:
		return slog.LevelDebug
	case loggo.INFO:
		return slog.LevelInfo
	case loggo.WARNING:
		return slog.LevelInfo + 1
	case loggo.ERROR:
		return slog.LevelError
	case loggo.CRITICAL:
		return slog.LevelError + 1
	default:
		panic("unknown level")
	}
}
