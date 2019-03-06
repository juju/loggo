package loggocolor

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/juju/ansiterm"
	"github.com/juju/loggo"
)

var (
	// SeverityColor defines the colors for the levels output by the ColorWriter.
	SeverityColor = map[loggo.Level]*ansiterm.Context{
		loggo.TRACE:   ansiterm.Foreground(ansiterm.Default),
		loggo.DEBUG:   ansiterm.Foreground(ansiterm.Green),
		loggo.INFO:    ansiterm.Foreground(ansiterm.BrightBlue),
		loggo.WARNING: ansiterm.Foreground(ansiterm.Yellow),
		loggo.ERROR:   ansiterm.Foreground(ansiterm.BrightRed),
		loggo.CRITICAL: &ansiterm.Context{
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

// NewWriter will write out colored severity levels if the writer is
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
func (w *colorWriter) Write(entry loggo.Entry) {
	ts := entry.Timestamp.Format(loggo.TimeFormat)
	// Just get the basename from the filename
	filename := filepath.Base(entry.Filename)

	fmt.Fprintf(w.writer, "%s ", ts)
	SeverityColor[entry.Level].Fprintf(w.writer, entry.Level.Short())
	fmt.Fprintf(w.writer, " %s ", entry.Module)
	LocationColor.Fprintf(w.writer, "%s:%d ", filename, entry.Line)
	fmt.Fprintln(w.writer, entry.Message)
}
