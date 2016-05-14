// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"reflect"
)

// A SubLogger represents a logging module. It has an associated logging
// level which can be changed; messages of lesser severity will
// be dropped. SubLoggers have a hierarchical relationship - see
// the package documentation.
//
// The zero SubLogger value is usable - any messages logged
// to it will be sent to the global root logger.
type SubLogger struct {
	logger
	*module
	writer RecordWriter
}

// NewRootLogger creates a root logger and returns it, along with
// the writers that the logger will use.
func NewRootLogger() (SubLogger, *Writers) {
	writers := NewWriters(nil) // starts off empty
	module := newRootModule()
	root := newSubLogger(module, writers)
	return root, writers
}

// NewSubLogger creates a new SubLogger with the given name and parent.
// The logger uses the name to identify itself. The parent is used when
// determining the effective log level. The new logger is returned,
// along with the writers the logger will use.
func NewSubLogger(name string, parent SubLogger) (SubLogger, *Writers) {
	var writers *Writers
	if parent.isZero() {
		parent, writers = NewRootLogger()
	} else {
		writers = NewWriters(nil) // starts off empty
		if parent.writer != nil {
			// We set the level as low as possible in order to defer
			// strictly to the new logger's level.
			// TODO(ericsnow) Use parent.writer's level.
			writers.AddWithLevel(defaultWriterName, parent.writer, UNSPECIFIED)
		}
	}
	module := newSubmodule(name, parent.module, defaultLevel)
	logger := newSubLogger(module, writers)
	return logger, writers
}

func newSubLogger(module *module, writer MinLevelWriter) SubLogger {
	return SubLogger{
		logger: logger{&callLogger{
			HasMinLevel: &module.loggerState,
			st:          &module.loggerState,
			writer:      writer,
		}},
		module: module,
		writer: writer,
	}
}

func (logger SubLogger) isZero() bool {
	return reflect.DeepEqual(logger, SubLogger{})
}

// SetLogLevel sets the severity level of the given logger.
func (logger *SubLogger) SetLogLevel(level Level) {
	logger.setLevel(level)
}
