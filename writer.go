// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"fmt"
	"io"
	"time"
)

// defaultWriterName is the name of the writer default writer.
const defaultWriterName = "default"

// TODO(ericsnow) Drop legacyWriter once it's no longer used.

// Writer is implemented by any recipient of log messages.
//
// Writer is deprecated. Use RecordWriter instead.
type Writer legacyWriter

type legacyWriter interface {
	// Write writes a message to the Writer with the given
	// level and module name. The filename and line hold
	// the file name and line number of the code that is
	// generating the log message; the time stamp holds
	// the time the log message was generated, and
	// message holds the log message itself.
	Write(level Level, name, filename string, line int, timestamp time.Time, message string)
}

// LegacyCompatibleWriter is a shim to temporarily support both interfaces.
type LegacyCompatibleWriter interface {
	RecordWriter
	legacyWriter
}

// RecordWriter is implemented by any recipient of log messages.
type RecordWriter interface {
	// WriteRecord writes a message to the Writer for the given
	// log record.
	WriteRecord(Record)
}

type legacyWriterShim struct {
	RecordWriter
}

func (shim legacyWriterShim) Write(level Level, loggerName, filename string, line int, timestamp time.Time, message string) {
	shim.WriteRecord(Record{
		Level:      level,
		LoggerName: loggerName,
		Filename:   filename,
		Line:       line,
		Timestamp:  timestamp,
		Message:    message,
	})
}

type legacyAdaptingWriter struct {
	legacyWriter
}

func (law *legacyAdaptingWriter) WriteRecord(rec Record) {
	law.Write(rec.Level, rec.LoggerName, rec.Filename, rec.Line, rec.Timestamp, rec.Message)
}

// MinLevelWriter is a writer that exposes its minimum log level.
type MinLevelWriter interface {
	RecordWriter
	HasMinLevel
}

type minLevelWriter struct {
	writer RecordWriter
	level  Level
}

// NewMinLevelWriter returns a MinLevelWriter that wraps the given
// writer with the provided min log level.
func NewMinLevelWriter(writer RecordWriter, minLevel Level) MinLevelWriter {
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
func (w minLevelWriter) WriteRecord(rec Record) {
	if !IsLevelEnabled(&w, rec.Level) {
		return
	}
	w.writer.WriteRecord(rec)
}

// TODO(ericsnow) Eliminate NewSimpleWriter().

// NewSimpleWriter returns a new writer that writes
// log messages to the given io.Writer formatting the
// messages with the given formatter.
func NewSimpleWriter(writer io.Writer, formatter LegacyFormatter) LegacyCompatibleWriter {
	return &legacyWriterShim{
		&formattingWriter{writer, &legacyAdaptingFormatter{formatter}},
	}
}

// formattingWriter is a log writer that writes
// log messages to the given io.Writer, formatting the
// messages with the given formatter.
type formattingWriter struct {
	writer    io.Writer
	formatter Formatter
}

// NewFormattingWriter returns a new writer that writes
// log messages to the given io.Writer, formatting the
// messages with the given formatter.
func NewFormattingWriter(writer io.Writer, formatter Formatter) RecordWriter {
	return &formattingWriter{
		writer:    writer,
		formatter: formatter,
	}
}

// Write formats the record and writes the result to the io.Writer.
func (fw *formattingWriter) WriteRecord(rec Record) {
	var logLine string
	if fw.formatter == nil {
		logLine = rec.String()
	} else {
		logLine = fw.formatter.Format(rec)
	}
	fmt.Fprintln(fw.writer, logLine)
}
