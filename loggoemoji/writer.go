// Copyright 2021 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggoemoji

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/juju/ansiterm"
	"github.com/juju/loggo"
)

type levelContext struct {
	Emoji string
	Style *ansiterm.Context
}

var (
	// SeverityEmoji defines the colors for the levels output by the ColorWriter.
	SeverityEmoji = map[loggo.Level]levelContext{
		loggo.TRACE:   {Emoji: "✏️", Style: ansiterm.Foreground(ansiterm.Default)},
		loggo.DEBUG:   {Emoji: "🐞", Style: ansiterm.Foreground(ansiterm.Green)},
		loggo.INFO:    {Emoji: "🧐", Style: ansiterm.Foreground(ansiterm.BrightBlue)},
		loggo.WARNING: {Emoji: "⚠️ ", Style: ansiterm.Foreground(ansiterm.Yellow)},
		loggo.ERROR:   {Emoji: "😱", Style: ansiterm.Foreground(ansiterm.BrightRed)},
		loggo.CRITICAL: {Emoji: "💥", Style: &ansiterm.Context{
			Foreground: ansiterm.White,
			Background: ansiterm.Red,
		}},
	}
	// LocationColor defines the colors for the location output by the ColorWriter.
	LocationColor = ansiterm.Foreground(ansiterm.BrightBlue)
)

type emojiWriter struct {
	writer *ansiterm.Writer
}

// NewWriter will write out colored severity levels if the writer is
// outputting to a terminal.
func NewWriter(writer io.Writer) loggo.Writer {
	return &emojiWriter{ansiterm.NewWriter(writer)}
}

// NewColorWriter will write out colored severity levels whether or not the
// writer is outputting to a terminal.
func NewColorWriter(writer io.Writer) loggo.Writer {
	w := ansiterm.NewWriter(writer)
	w.SetColorCapable(true)
	return &emojiWriter{w}
}

// Write implements Writer.
func (w *emojiWriter) Write(entry loggo.Entry) {
	ts := entry.Timestamp.Format(loggo.TimeFormat)
	// Just get the basename from the filename
	filename := filepath.Base(entry.Filename)

	fmt.Fprintf(w.writer, "%s ", ts)
	SeverityEmoji[entry.Level].Style.Fprintf(w.writer, entry.Level.Short())
	fmt.Fprintf(w.writer, " %s", SeverityEmoji[entry.Level].Emoji)
	fmt.Fprintf(w.writer, " %s ", entry.Module)
	LocationColor.Fprintf(w.writer, "%s:%d ", filename, entry.Line)
	fmt.Fprintln(w.writer, entry.Message)
}
