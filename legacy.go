// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"io"
	"time"
)

// TODO(ericsnow) The contents of this file can go away as soon as we
// stop worrying about breaking compatibility.

// LegacyFormatter defines the single method Format, which takes the logging
// information, and converts it to a string.
type LegacyFormatter interface {
	Format(level Level, loggerName, filename string, line int, timestamp time.Time, message string) string
}

type legacyAdaptingFormatter struct {
	legacy LegacyFormatter
}

func (f *legacyAdaptingFormatter) Format(rec Record) string {
	return f.legacy.Format(rec.Level, rec.LoggerName, rec.Filename, rec.Line, rec.Timestamp, rec.Message)
}

// DefaultFormatter provides a simple concatenation of all the components.
//
// DefaultFormatter is deprecated. Pass nil to NewFormattingWriter() instead.
type DefaultFormatter struct{}

// Format returns the parameters separated by spaces except for filename and
// line which are separated by a colon.  The timestamp is shown to second
// resolution in UTC.
func (*DefaultFormatter) Format(level Level, loggerName, filename string, line int, timestamp time.Time, message string) string {
	rec := Record{
		Level:      level,
		LoggerName: loggerName,
		Filename:   filename,
		Line:       line,
		Timestamp:  timestamp,
		Message:    message,
	}
	return rec.String()
}

// Writer is implemented by any recipient of log messages.
//
// Writer is deprecated. Use RecordWriter instead.
type Writer legacyWriter

type legacyWriter interface {
	// Write writes a message to the Writer with the given
	// level and logger name. The filename and line hold
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

// NewSimpleWriter returns a new writer that writes
// log messages to the given io.Writer formatting the
// messages with the given formatter.
//
// NewSimpleWriter is deprecated. Use NewFormattingWriter.
func NewSimpleWriter(writer io.Writer, formatter LegacyFormatter) LegacyCompatibleWriter {
	return &legacyWriterShim{
		&formattingWriter{writer, &legacyAdaptingFormatter{formatter}},
	}
}

// TODO(ericsnow) ...also, the Write() method of loggotest.Writer.
