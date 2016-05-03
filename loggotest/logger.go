// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggotest

import (
	"github.com/juju/loggo"
)

// TraceLogger returns the named logger. It also sets the logger's
// writer and returns it.
func TraceLogger(name string) (loggo.Logger, *Writer) {
	writer := &Writer{}
	loggo.ReplaceDefaultWriter(writer)
	logger := loggo.GetLogger(name)
	// Make it so the logger itself writes all messages.
	logger.SetLogLevel(loggo.TRACE)
	return logger, writer
}
