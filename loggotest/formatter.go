// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggotest

import (
	"github.com/juju/loggo"
)

// Formatter is a useful Writer for testing purposes. Each component
// of the logging message is stored in the Log array.
type Formatter struct {
	Writer

	format func(loggo.Record) string
}

// NewFormatter returns a new Formatter that wraps the given
// format func. If the func is nil then Format() will return the message.
func NewFormatter(format func(loggo.Record) string) *Formatter {
	return &Formatter{
		format: format,
	}
}

// Format saves the params as members in the TestLogValues struct appended to the Log array.
func (f *Formatter) Format(rec loggo.Record) string {
	f.Write(rec)
	if f.format == nil {
		return rec.Message
	}
	return f.format(rec)
}
