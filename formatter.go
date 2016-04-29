// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"fmt"
	"path/filepath"
	"time"
)

// Formatter defines the single method Format, which takes the logging
// record and converts it to a string.
type Formatter interface {
	Format(Record) string
}

// DefaultFormatter provides a simple concatenation of all the components.
type DefaultFormatter struct{}

// Format returns the parameters separated by spaces except for filename and
// line which are separated by a colon.  The timestamp is shown to second
// resolution in UTC.
func (*DefaultFormatter) Format(rec Record) string {
	ts := rec.Timestamp.In(time.UTC).Format("2006-01-02 15:04:05")
	// Just get the basename from the filename
	filename := filepath.Base(rec.Filename)
	return fmt.Sprintf("%s %s %s %s:%d %s", ts, rec.Level, rec.LoggerName, filename, rec.Line, rec.Message)
}

// TODO(ericsnow) The remainder of this file can go away when we fix
// the NewSimpleWriter() signature.

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
