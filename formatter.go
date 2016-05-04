// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"time"
)

// Formatter defines the single method Format, which takes the logging
// record and converts it to a string.
type Formatter interface {
	Format(Record) string
}

// TODO(ericsnow) The remainder of this file can go away when
// NewSimpleWriter() does.

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
