package loggocolor

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/juju/ansiterm"
	"github.com/juju/loggo"
	"strconv"
)

var (
	// SeverityColor defines the colors for the levels output by the ColorWriter.
	SeverityColor = map[loggo.Level]*ansiterm.Context{
		loggo.TRACE:   ansiterm.Foreground(ansiterm.Gray),
		loggo.DEBUG:   ansiterm.Foreground(ansiterm.Green),
		loggo.INFO:    ansiterm.Foreground(ansiterm.BrightBlue),
		loggo.NOTICE:  ansiterm.Foreground(ansiterm.BrightGreen),
		loggo.WARNING: ansiterm.Foreground(ansiterm.Yellow),
		loggo.ERROR:   ansiterm.Foreground(ansiterm.BrightRed),
		loggo.CRITICAL: {
			Foreground: ansiterm.White,
			Background: ansiterm.Red,
		},
		loggo.ALERT: {
			Foreground: ansiterm.White,
			Background: ansiterm.Red,
		},
		loggo.EMERGENCY: {
			Foreground: ansiterm.White,
			Background: ansiterm.Red,
		},
	}
	// LocationColor defines the colors for the location output by the ColorWriter.
	LocationColor = ansiterm.Foreground(ansiterm.BrightBlue)
	// TimeStampColor defines the colors of timestamps
	TimeStampColor = ansiterm.Foreground(ansiterm.Yellow)
	// ModuleColor defines the color of module name
	ModuleColor = ansiterm.Foreground(ansiterm.Blue)

	// How long (padded) should be the module name
	ModuleLength = 35

	// How long (padded) should be the location name
	LocationLength = 25
)

type colorWriter struct {
	writer *ansiterm.Writer
}

// NewColorWriter will write out colored severity levels if the writer is
// outputting to a terminal.
func NewWriter(writer io.Writer) loggo.Writer {
	return &colorWriter{ansiterm.NewWriter(writer)}
}

// Write implements Writer.
func (w *colorWriter) Write(entry loggo.Entry) {
	ts := entry.Timestamp.Format(loggo.TimeFormat)
	// Just get the basename from the filename
	filename := filepath.Base(entry.Filename)

	TimeStampColor.Fprintf(w.writer, "%s", ts)
	fmt.Fprintf(w.writer, " ")

	SeverityColor[entry.Level].Fprintf(w.writer, "%5s", entry.Level.Short())
	fmt.Fprintf(w.writer, " ")

	module := entry.Module
	if len(module) > ModuleLength {
		module = "..." + module[len(module)-(ModuleLength-3):]
	}

	ModuleColor.Fprintf(w.writer, "%-"+strconv.Itoa(ModuleLength)+"s", module)
	fmt.Fprintf(w.writer, " ")

	line := fmt.Sprintf("%s:%d", filename, entry.Line)
	LocationColor.Fprintf(w.writer, "%-"+strconv.Itoa(LocationLength)+"s ", line)

	fmt.Fprintln(w.writer, entry.Message)
}
