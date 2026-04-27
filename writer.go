// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"context"
	"fmt"
	"io"
	"os"
)

// DefaultWriterName is the name of the default writer for
// a Context.
const DefaultWriterName = "default"

// Writer is implemented by any recipient of log messages.
type Writer interface {
	// Write writes a message to the Writer with the given level and module
	// name. The filename and line hold the file name and line number of the
	// code that is generating the log message; the time stamp holds the time
	// the log message was generated, and message holds the log message
	// itself.
	Write(ctx context.Context, entry Entry) error
}

// NewMinLevelWriter returns a Writer that will only pass on the Write calls
// to the provided writer if the log level is at or above the specified
// minimum level.
func NewMinimumLevelWriter(writer Writer, minLevel Level) Writer {
	return &minLevelWriter{
		writer: writer,
		level:  minLevel,
	}
}

type minLevelWriter struct {
	writer Writer
	level  Level
}

// Write writes the log record.
func (w minLevelWriter) Write(ctx context.Context, entry Entry) error {
	if entry.Level < w.level {
		return nil
	}
	return w.writer.Write(ctx, entry)
}

type simpleWriter struct {
	writer    io.Writer
	formatter func(entry Entry) string
}

// NewSimpleWriter returns a new writer that writes log messages to the given
// io.Writer formatting the messages with the given formatter.
func NewSimpleWriter(writer io.Writer, formatter func(entry Entry) string) Writer {
	if formatter == nil {
		formatter = DefaultFormatter
	}
	return &simpleWriter{writer: writer, formatter: formatter}
}

func (simple *simpleWriter) Write(ctx context.Context, entry Entry) error {
	logLine := simple.formatter(entry)
	_, err := fmt.Fprintln(simple.writer, logLine)
	return err
}

func defaultWriter() Writer {
	return NewSimpleWriter(os.Stderr, DefaultFormatter)
}
