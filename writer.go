// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"fmt"
	"io"
	"time"
)

// Writer is implemented by any recipient of log messages.
type Writer interface {
	// Write writes a message to the Writer with the given
	// level and module name. The filename and line hold
	// the file name and line number of the code that is
	// generating the log message; the time stamp holds
	// the time the log message was generated, and
	// message holds the log message itself.
	Write(level Level, name, filename string, line int, timestamp time.Time, message string)
}

type registeredWriter struct {
	writer Writer
	level  Level
}

type simpleWriter struct {
	writer    io.Writer
	formatter Formatter
}

// NewSimpleWriter returns a new writer that writes
// log messages to the given io.Writer formatting the
// messages with the given formatter.
func NewSimpleWriter(writer io.Writer, formatter Formatter) Writer {
	return &simpleWriter{writer, formatter}
}

func (simple *simpleWriter) Write(level Level, module, filename string, line int, timestamp time.Time, message string) {
	logLine := simple.formatter.Format(level, module, filename, line, timestamp, message)
	fmt.Fprintln(simple.writer, logLine)
}
