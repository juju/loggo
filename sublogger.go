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
	writer Writer
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

// TODO(ericsnow) Everything below here is unnecessary now and should
// be deprecated.

func (logger SubLogger) getModule() *module {
	if logger.module == nil {
		return newRootModule()
	}
	return logger.module
}

// LogLevel returns the configured min log level of the logger.
//
// This is strictly an alias for MinLogLevel, intended to align with
// a more commonly used (but less specific) name.
func (logger SubLogger) LogLevel() Level {
	return logger.getModule().MinLogLevel()
}

// EffectiveLogLevel returns the effective min log level of
// the receiver - that is, messages with a lesser severity
// level will be discarded.
//
// If the log level of the receiver is unspecified,
// it will be taken from the effective log level of its
// parent.
func (logger SubLogger) EffectiveLogLevel() Level {
	return EffectiveMinLevel(logger.getModule())
}

// IsLevelEnabled returns whether debugging is enabled
// for the given log level.
func (logger SubLogger) IsLevelEnabled(level Level) bool {
	return IsLevelEnabled(logger.getModule(), level)
}

// IsErrorEnabled returns whether debugging is enabled
// at error level.
func (logger SubLogger) IsErrorEnabled() bool {
	return logger.IsLevelEnabled(ERROR)
}

// IsWarningEnabled returns whether debugging is enabled
// at warning level.
func (logger SubLogger) IsWarningEnabled() bool {
	return logger.IsLevelEnabled(WARNING)
}

// IsInfoEnabled returns whether debugging is enabled
// at info level.
func (logger SubLogger) IsInfoEnabled() bool {
	return logger.IsLevelEnabled(INFO)
}

// IsDebugEnabled returns whether debugging is enabled
// at debug level.
func (logger SubLogger) IsDebugEnabled() bool {
	return logger.IsLevelEnabled(DEBUG)
}

// IsTraceEnabled returns whether debugging is enabled
// at trace level.
func (logger SubLogger) IsTraceEnabled() bool {
	return logger.IsLevelEnabled(TRACE)
}
