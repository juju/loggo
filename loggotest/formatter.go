// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggotest

import (
	"time"

	"github.com/juju/loggo"
)

// Formatter is a useful Writer for testing purposes. Each component
// of the logging message is stored in the Log array.
type Formatter struct {
	Writer

	format func(level loggo.Level, module, filename string, line int, timestamp time.Time, message string) string
}

// NewFormatter returns a new Formatter that wraps the given
// format func. If the func is nil then Format() will return the message.
func NewFormatter(format func(level loggo.Level, module, filename string, line int, timestamp time.Time, message string) string) *Formatter {
	return &Formatter{
		format: format,
	}
}

// Format saves the params as members in the TestLogValues struct appended to the Log array.
func (f *Formatter) Format(level loggo.Level, module, filename string, line int, timestamp time.Time, message string) string {
	f.Write(level, module, filename, line, timestamp, message)
	if f.format == nil {
		return message
	}
	return f.format(level, module, filename, line, timestamp, message)
}
