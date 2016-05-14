// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"log"
)

// LoggerAsGoLogger wraps the logger in a stdlib log.Logger. The
// messages are logged at the provided level.
func LoggerAsGoLogger(logger Logger, level Level) *log.Logger {
	w := LoggerAsIOWriter(logger, level)
	return log.New(w, "", 0)
}

// IOAdapter is an io.Writer that logs written messages.
type IOAdapter struct {
	logger Logger
	level  Level
}

// LoggerAsIOWriter returns a new io.Writer that logs the written messages
// at the given log level.
func LoggerAsIOWriter(logger Logger, level Level) *IOAdapter {
	return &IOAdapter{
		logger: logger,
		level:  level,
	}
}

// Write implements io.Writer, logging the message at the predefined log level.
func (w IOAdapter) Write(msg []byte) (int, error) {
	n := len(msg)
	// Same calldepth as in Logf + 2 for log.Logger.
	w.logger.LogCallf(5, w.level, string(msg))
	return n, nil
}
