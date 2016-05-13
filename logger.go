// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"fmt"
	"runtime"
	"time"
)

const (
	defaultRootLevel = WARNING
	defaultLevel     = UNSPECIFIED
)

// A Logger represents a single logger. It has an associated
// logging level; messages of lesser severity will be dropped.
type Logger interface {
	HasMinLevel

	// Logf logs a printf-formatted message at the given level.
	// A message will be discarded if level is less than the
	// the effective log level of the logger.
	// Note that the writers may also filter out messages that
	// are less than their registered minimum severity level.
	Logf(level Level, message string, args ...interface{})

	// LogCallf logs a printf-formatted message at the given level.
	// The location of the call is indicated by the calldepth argument.
	// A calldepth of 1 means the function that called this function.
	// A message will be discarded if level is less than the
	// the effective log level of the logger.
	// Note that the writers may also filter out messages that
	// are less than their registered minimum severity level.
	LogCallf(calldepth int, level Level, message string, args ...interface{})

	// Criticalf logs the printf-formatted message at critical level.
	Criticalf(message string, args ...interface{})

	// Errorf logs the printf-formatted message at error level.
	Errorf(message string, args ...interface{})

	// Warningf logs the printf-formatted message at warning level.
	Warningf(message string, args ...interface{})

	// Infof logs the printf-formatted message at info level.
	Infof(message string, args ...interface{})

	// Debugf logs the printf-formatted message at debug level.
	Debugf(message string, args ...interface{})

	// Tracef logs the printf-formatted message at trace level.
	Tracef(message string, args ...interface{})
}

// A ConfigurableLogger represents a single logger. It has an associated
// logging level which can be changed; messages of lesser severity
// will be dropped.
type ConfigurableLogger interface {
	Logger

	// SetLogLevel sets the severity level of the given logger.
	SetLogLevel(Level)
}

// loggerState holds the raw info about the logger.
//
// A nil loggerState pointer represents a read-only root logger.
type loggerState struct {
	name         string
	level        Level
	defaultLevel Level
	parent       *loggerState
}

// Name returns the logger's name.
func (st *loggerState) Name() string {
	if st == nil {
		return ""
	}
	return st.name
}

// MinLogLevel returns the configured minimum log level of the
// logger. This is the level at which messages with a lower level
// will be discarded.
func (st *loggerState) MinLogLevel() Level {
	if st == nil {
		return defaultRootLevel
	}
	return st.level.get()
}

// setLevel sets the severity level of the logger..
//
// This will panic if the loggerState pointer is nil.
func (st *loggerState) setLevel(level Level) {
	if level == UNSPECIFIED {
		level = st.defaultLevel
	}
	st.level.set(level)
}

// ParentWithMinLogLevel returns the logger's parent (or nil).
func (st *loggerState) ParentWithMinLogLevel() HasMinLevel {
	if st == nil {
		return nil
	}
	if st.parent == nil { // avoid double nil
		return nil
	}
	return st.parent
}

// config returns the current configuration for the logger.
func (st *loggerState) config() LoggerConfig {
	return LoggerConfig{
		Level: st.MinLogLevel(),
	}
}

// applyConfig configures the logger according to the provided config.
//
// This will panic if the loggerState pointer is nil.
func (st *loggerState) applyConfig(cfg LoggerConfig) {
	st.setLevel(cfg.Level)
}

// A logger represents an independent logger. It has an associated
// logging level which can be changed; messages of lesser severity
// will be dropped.
type logger struct {
	*loggerState
	writer MinLevelWriter
}

// NewLogger returns a new Logger that will use the given writer.
func NewLogger(writer MinLevelWriter) ConfigurableLogger {
	return &logger{
		loggerState: &loggerState{
			level: defaultLevel,
		},
		writer: writer,
	}
}

// SetLogLevel sets the severity level of the given logger.
func (logger *logger) SetLogLevel(level Level) {
	logger.setLevel(level)
}

// Logf logs a printf-formatted message at the given level.
// A message will be discarded if level is less than the
// the effective log level of the logger.
// Note that the writers may also filter out messages that
// are less than their registered minimum severity level.
func (logger logger) Logf(level Level, message string, args ...interface{}) {
	logger.LogCallf(2, level, message, args...)
}

// LogCallf logs a printf-formatted message at the given level.
// The location of the call is indicated by the calldepth argument.
// A calldepth of 1 means the function that called this function.
// A message will be discarded if level is less than the
// the effective log level of the logger.
// Note that the writers may also filter out messages that
// are less than their registered minimum severity level.
func (logger logger) LogCallf(calldepth int, level Level, message string, args ...interface{}) {
	if !logger.willWrite(level) {
		return
	}
	loggerName := logger.name
	if loggerName == "" {
		loggerName = "<>"
	}

	// Gather time, filename, and line number.
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
	if logger.writer != nil {
		logger.writer.Write(level, loggerName, file, line, now, formattedMessage)
	}
}

func (logger logger) willWrite(level Level) bool {
	if !IsLevelEnabled(logger, level) {
		return false
	}
	if !IsLevelEnabled(logger.writer, level) {
		return false
	}
	if level < TRACE || level > CRITICAL {
		return false
	}
	return true
}

// Criticalf logs the printf-formatted message at critical level.
func (logger logger) Criticalf(message string, args ...interface{}) {
	logger.Logf(CRITICAL, message, args...)
}

// Errorf logs the printf-formatted message at error level.
func (logger logger) Errorf(message string, args ...interface{}) {
	logger.Logf(ERROR, message, args...)
}

// Warningf logs the printf-formatted message at warning level.
func (logger logger) Warningf(message string, args ...interface{}) {
	logger.Logf(WARNING, message, args...)
}

// Infof logs the printf-formatted message at info level.
func (logger logger) Infof(message string, args ...interface{}) {
	logger.Logf(INFO, message, args...)
}

// Debugf logs the printf-formatted message at debug level.
func (logger logger) Debugf(message string, args ...interface{}) {
	logger.Logf(DEBUG, message, args...)
}

// Tracef logs the printf-formatted message at trace level.
func (logger logger) Tracef(message string, args ...interface{}) {
	logger.Logf(TRACE, message, args...)
}
