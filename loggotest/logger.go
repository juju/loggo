// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggotest

import (
	"github.com/juju/loggo"
)

// TraceLogger returns the named logger. It also sets the logger's
// writer and returns it.
func TraceLogger() (loggo.Logger, *Writer) {
	logger, writers := loggo.NewRootLogger()
	// Make it so the logger itself writes all messages.
	logger.SetLogLevel(loggo.TRACE)

	writer := &Writer{}
	writers.AddWithLevel("", writer, loggo.TRACE)

	return logger, writer
}
