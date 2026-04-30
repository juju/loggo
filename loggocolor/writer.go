package loggocolor

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/juju/ansiterm"
	"github.com/juju/loggo/v3"
	"github.com/juju/loggo/v3/attrs"
)

var (
	// SeverityColor defines the colors for the levels output by the ColorWriter.
	SeverityColor = map[loggo.Level]*ansiterm.Context{
		loggo.TRACE:   ansiterm.Foreground(ansiterm.Default),
		loggo.DEBUG:   ansiterm.Foreground(ansiterm.Green),
		loggo.INFO:    ansiterm.Foreground(ansiterm.BrightBlue),
		loggo.WARNING: ansiterm.Foreground(ansiterm.Yellow),
		loggo.ERROR:   ansiterm.Foreground(ansiterm.BrightRed),
		loggo.CRITICAL: {
			Foreground: ansiterm.White,
			Background: ansiterm.Red,
		},
	}
	// LocationColor defines the colors for the location output by the ColorWriter.
	LocationColor = ansiterm.Foreground(ansiterm.BrightBlue)
)

type colorWriter struct {
	writer *ansiterm.Writer
}

// NewColorWriter will write out colored severity levels if the writer is
// outputting to a terminal.
func NewWriter(writer io.Writer) loggo.Writer {
	return &colorWriter{ansiterm.NewWriter(writer)}
}

// NewcolorWriter will write out colored severity levels whether or not the
// writer is outputting to a terminal.
func NewColorWriter(writer io.Writer) loggo.Writer {
	w := ansiterm.NewWriter(writer)
	w.SetColorCapable(true)
	return &colorWriter{w}
}

// Write implements Writer.
func (w *colorWriter) Write(_ context.Context, entry loggo.Entry) error {
	ts := entry.Timestamp.Format(loggo.TimeFormat)
	// Just get the basename from the filename
	filename := filepath.Base(entry.Filename)

	if _, err := fmt.Fprintf(w.writer, "%s ", ts); err != nil {
		return err
	}

	SeverityColor[entry.Level].Fprintf(w.writer, "%s", entry.Level.Short())
	if _, err := fmt.Fprintf(w.writer, " %s ", entry.Module); err != nil {
		return err
	}
	LocationColor.Fprintf(w.writer, "%s:%d ", filename, entry.Line)
	if _, err := fmt.Fprintln(w.writer, entry.Message); err != nil {
		return err
	}

	for _, attr := range entry.Attrs {
		switch a := attr.(type) {
		case attrs.AttrValue[string]:
			if _, err := fmt.Fprintf(w.writer, "  %s=%s\n", a.Key(), a.Value()); err != nil {
				return err
			}
		case attrs.AttrValue[int]:
			if _, err := fmt.Fprintf(w.writer, "  %s=%d\n", a.Key(), a.Value()); err != nil {
				return err
			}
		case attrs.AttrValue[int64]:
			if _, err := fmt.Fprintf(w.writer, "  %s=%d\n", a.Key(), a.Value()); err != nil {
				return err
			}
		case attrs.AttrValue[uint64]:
			if _, err := fmt.Fprintf(w.writer, "  %s=%d\n", a.Key(), a.Value()); err != nil {
				return err
			}
		case attrs.AttrValue[float64]:
			if _, err := fmt.Fprintf(w.writer, "  %s=%f\n", a.Key(), a.Value()); err != nil {
				return err
			}
		case attrs.AttrValue[bool]:
			if _, err := fmt.Fprintf(w.writer, "  %s=%t\n", a.Key(), a.Value()); err != nil {
				return err
			}
		case attrs.AttrValue[time.Time]:
			if _, err := fmt.Fprintf(w.writer, "  %s=%v\n", a.Key(), a.Value()); err != nil {
				return err
			}
		case attrs.AttrValue[time.Duration]:
			if _, err := fmt.Fprintf(w.writer, "  %s=%v\n", a.Key(), a.Value()); err != nil {
				return err
			}
		case attrs.AttrValue[any]:
			if _, err := fmt.Fprintf(w.writer, "  %s=%v\n", a.Key(), a.Value()); err != nil {
				return err
			}
		}
	}
	return nil
}
