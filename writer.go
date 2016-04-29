// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"fmt"
	"io"
)

// defaultWriterName is the name of the writer default writer.
const defaultWriterName = "default"

// Writer is implemented by any recipient of log messages.
type Writer interface {
	// Write writes a message to the Writer with the given
	// log record.
	Write(Record)
}

// MinLevelWriter is a writer that exposes its minimum log level.
type MinLevelWriter interface {
	Writer
	HasMinLevel
}

type minLevelWriter struct {
	writer Writer
	level  Level
}

// NewMinLevelWriter returns a MinLevelWriter that wraps the given
// writer with the provided min log level.
func NewMinLevelWriter(writer Writer, minLevel Level) MinLevelWriter {
	return &minLevelWriter{
		writer: writer,
		level:  minLevel,
	}
}

// MinLogLevel returns the writer's log level.
func (w minLevelWriter) MinLogLevel() Level {
	return w.level
}

// Write writes the log record.
func (w minLevelWriter) Write(rec Record) {
	if !IsLevelEnabled(&w, rec.Level) {
		return
	}
	w.writer.Write(rec)
}

// TODO(ericsnow) Eliminate NewSimpleWriter().

// NewSimpleWriter returns a new writer that writes
// log messages to the given io.Writer formatting the
// messages with the given formatter.
func NewSimpleWriter(writer io.Writer, formatter LegacyFormatter) Writer {
	return &formattingWriter{writer, &legacyAdaptingFormatter{formatter}}
}

type formattingWriter struct {
	writer    io.Writer
	formatter Formatter
}

// NewFormattingWriter returns a new writer that writes
// log messages to the given io.Writer formatting the
// messages with the given formatter.
func NewFormattingWriter(writer io.Writer, formatter Formatter) Writer {
	return &formattingWriter{
		writer:    writer,
		formatter: formatter,
	}
}

func (fw *formattingWriter) Write(rec Record) {
	logLine := fw.formatter.Format(rec)
	fmt.Fprintln(fw.writer, logLine)
}
