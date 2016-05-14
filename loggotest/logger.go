// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggotest

import (
	"github.com/juju/loggo"
)

// Logger returns the named logger. It also sets the logger's
// writer and returns it.
func Logger(level loggo.Level) (loggo.ConfigurableLogger, *Writer) {
	writer := &Writer{}

	logger := loggo.NewLogger(loggo.NewMinLevelWriter(writer, level))
	// Make it so the logger itself writes all messages.
	logger.SetLogLevel(loggo.TRACE)

	return logger, writer
}

// TraceLogger returns the named logger. It also sets the logger's
// writer and returns it.
func TraceLogger() (loggo.ConfigurableLogger, *Writer) {
	return Logger(loggo.TRACE)
}
