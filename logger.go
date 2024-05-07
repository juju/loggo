// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// A Logger represents a logging module. It has an associated logging
// level which can be changed; messages of lesser severity will
// be dropped. Loggers have a hierarchical relationship - see
// the package documentation.
//
// The zero Logger value is usable - any messages logged
// to it will be sent to the root Logger.
type Logger struct {
	impl *module
}

func (logger Logger) getModule() *module {
	if logger.impl == nil {
		return defaultContext.root
	}
	return logger.impl
}

// Root returns the root logger for the Logger's context.
func (logger Logger) Root() Logger {
	module := logger.getModule()
	return module.context.GetLogger("")
}

// Parent returns the Logger whose module name is the same
// as this logger without the last period and suffix.
// For example the parent of the logger that has the module
// "a.b.c" is "a.b".
// The Parent of the root logger is still the root logger.
func (logger Logger) Parent() Logger {
	return Logger{impl: logger.getModule().parent}
}

// Child returns the Logger whose module name is the composed of this
// Logger's name and the specified name.
func (logger Logger) Child(name string) Logger {
	module := logger.getModule()
	path := module.name
	if path == "" {
		path = name
	} else {
		path += "." + name
	}
	return module.context.GetLogger(path)
}

// ChildWithTags returns the Logger whose module name is the composed of this
// Logger's name and the specified name with the correct associated tags.
func (logger Logger) ChildWithTags(name string, tags ...string) Logger {
	module := logger.getModule()
	path := module.name
	if path == "" {
		path = name
	} else {
		path += "." + name
	}
	return module.context.GetLogger(path, tags...)
}

// ChildWithLabels returns the Logger whose module name is the composed of this
// Logger's name and the specified name with the correct associated labels.
// Adding labels to the child logger will cause all child loggers to also
// inherit the labels of the parent(s) loggers.
// For targeting a singular logger with labels, use WithLabels which are not
// inherited by child loggers.
func (logger Logger) ChildWithLabels(name string, labels Labels) Logger {
	module := logger.getModule()
	path := module.name
	if path == "" {
		path = name
	} else {
		path += "." + name
	}

	merged := make(Labels)
	for k, v := range logger.impl.labels {
		merged[k] = v
	}
	for k, v := range labels {
		merged[k] = v
	}

	result := module.context.GetLogger(path)
	result.impl.labels = merged
	return result
}

// Name returns the logger's module name.
func (logger Logger) Name() string {
	return logger.getModule().Name()
}

// LogLevel returns the configured min log level of the logger.
func (logger Logger) LogLevel() Level {
	return logger.getModule().level
}

// Tags returns the configured tags of the logger's module.
func (logger Logger) Tags() []string {
	return logger.getModule().tags
}

// Labels returns the configured labels of the logger's module.
func (logger Logger) Labels() Labels {
	return logger.getModule().labels
}

// EffectiveLogLevel returns the effective min log level of
// the receiver - that is, messages with a lesser severity
// level will be discarded.
//
// If the log level of the receiver is unspecified,
// it will be taken from the effective log level of its
// parent.
func (logger Logger) EffectiveLogLevel() Level {
	return logger.getModule().getEffectiveLogLevel()
}

// SetLogLevel sets the severity level of the given logger.
// The root logger cannot be set to UNSPECIFIED level.
// See EffectiveLogLevel for how this affects the
// actual messages logged.
func (logger Logger) SetLogLevel(level Level) {
	logger.getModule().setLevel(level)
}

// Logf logs a printf-formatted message at the given level.
// A message will be discarded if level is less than the
// the effective log level of the logger.
// Note that the writers may also filter out messages that
// are less than their registered minimum severity level.
func (logger Logger) Logf(level Level, message string, args ...interface{}) {
	logger.LogCallf(2, level, message, args...)
}

// LogCallf logs a printf-formatted message at the given level.
// The location of the call is indicated by the calldepth argument.
// A calldepth of 1 means the function that called this function.
// A message will be discarded if level is less than the
// the effective log level of the logger.
// Note that the writers may also filter out messages that
// are less than their registered minimum severity level.
func (logger Logger) LogCallf(calldepth int, level Level, message string, args ...interface{}) {
	logger.logCallf(calldepth+1, level, message, args...)
}

// logCallf is a private method for logging a printf-formatted message at the
// given level. Used by LogWithLabelsf and LogCallf.
func (logger Logger) logCallf(calldepth int, level Level, message string, args ...interface{}) {
	module := logger.getModule()
	if !module.willWrite(level) {
		return
	}
	// Gather time, and filename, line number.
	now := time.Now() // get this early.
	// Param to Caller is the call depth.  Since this method is called from
	// the Logger methods, we want the place that those were called from.
	_, file, line, ok := runtime.Caller(calldepth + 1)
	if !ok {
		file = "???"
		line = 0
	}
	// Trim newline off format string, following usual
	// Go logging conventions.
	if len(message) > 0 && message[len(message)-1] == '\n' {
		message = message[0 : len(message)-1]
	}

	// To avoid having a proliferation of Info/Infof methods,
	// only use Sprintf if there are any args, and rely on the
	// `go vet` tool for the obvious cases where someone has forgotten
	// to provide an arg.
	formattedMessage := message
	if len(args) > 0 {
		formattedMessage = fmt.Sprintf(message, args...)
	}

	entry := Entry{
		Level:     level,
		Filename:  file,
		Line:      line,
		Timestamp: now,
		Message:   formattedMessage,
	}
	entry.Labels = make(Labels)
	if len(module.tags) > 0 {
		entry.Labels[LoggerTags] = strings.Join(module.tags, ",")
	}
	for k, v := range logger.impl.labels {
		entry.Labels[k] = v
	}
	module.write(entry)
}

// Criticalf logs the printf-formatted message at critical level.
func (logger Logger) Criticalf(message string, args ...interface{}) {
	logger.Logf(CRITICAL, message, args...)
}

// Errorf logs the printf-formatted message at error level.
func (logger Logger) Errorf(message string, args ...interface{}) {
	logger.Logf(ERROR, message, args...)
}

// Warningf logs the printf-formatted message at warning level.
func (logger Logger) Warningf(message string, args ...interface{}) {
	logger.Logf(WARNING, message, args...)
}

// Infof logs the printf-formatted message at info level.
func (logger Logger) Infof(message string, args ...interface{}) {
	logger.Logf(INFO, message, args...)
}

// Debugf logs the printf-formatted message at debug level.
func (logger Logger) Debugf(message string, args ...interface{}) {
	logger.Logf(DEBUG, message, args...)
}

// Tracef logs the printf-formatted message at trace level.
func (logger Logger) Tracef(message string, args ...interface{}) {
	logger.Logf(TRACE, message, args...)
}

// IsLevelEnabled returns whether debugging is enabled
// for the given log level.
func (logger Logger) IsLevelEnabled(level Level) bool {
	return logger.getModule().willWrite(level)
}

// IsErrorEnabled returns whether debugging is enabled
// at error level.
func (logger Logger) IsErrorEnabled() bool {
	return logger.IsLevelEnabled(ERROR)
}

// IsWarningEnabled returns whether debugging is enabled
// at warning level.
func (logger Logger) IsWarningEnabled() bool {
	return logger.IsLevelEnabled(WARNING)
}

// IsInfoEnabled returns whether debugging is enabled
// at info level.
func (logger Logger) IsInfoEnabled() bool {
	return logger.IsLevelEnabled(INFO)
}

// IsDebugEnabled returns whether debugging is enabled
// at debug level.
func (logger Logger) IsDebugEnabled() bool {
	return logger.IsLevelEnabled(DEBUG)
}

// IsTraceEnabled returns whether debugging is enabled
// at trace level.
func (logger Logger) IsTraceEnabled() bool {
	return logger.IsLevelEnabled(TRACE)
}
