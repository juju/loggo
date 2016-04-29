// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"log"
)

const (
	defaultRootLevel = WARNING
	defaultLevel     = UNSPECIFIED
)

// CallLogger is the most basic kind logger in this package.
type CallLogger interface {
	// LogCallf logs a printf-formatted message at the given level.
	// The location of the call is indicated by the calldepth argument.
	// A calldepth of 1 means the function that called this function.
	// A message will be discarded if level is less than the
	// the effective log level of the logger.
	// Note that the writers may also filter out messages that
	// are less than their registered minimum severity level.
	LogCallf(calldepth int, level Level, message string, args ...interface{})
}

// A Logger represents a single logger. It has an associated
// logging level; messages of lesser severity will be dropped.
type Logger interface {
	HasMinLevel
	CallLogger

	// Logf logs a printf-formatted message at the given level.
	// A message will be discarded if level is less than the
	// the effective log level of the logger.
	// Note that the writers may also filter out messages that
	// are less than their registered minimum severity level.
	Logf(level Level, message string, args ...interface{})

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

// NewLogger returns a new Logger that will use the given writer.
func NewLogger(writer MinLevelWriter) ConfigurableLogger {
	st := &loggerState{
		level: defaultLevel,
	}
	return &simpleLogger{
		logger: logger{&callLogger{
			HasMinLevel: st,
			st:          st,
			writer:      writer,
		}},
		loggerState: st,
	}
}

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

type logger struct {
	CallLogger
}

// Logf logs a printf-formatted message at the given level.
// A message will be discarded if level is less than the
// the effective log level of the logger.
// Note that the writers may also filter out messages that
// are less than their registered minimum severity level.
func (logger *logger) Logf(level Level, message string, args ...interface{}) {
	if logger == nil || logger.CallLogger == nil {
		return
	}
	logger.LogCallf(3, level, message, args...)
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

// A simpleLogger represents an independent logger. It has an associated
// logging level which can be changed; messages of lesser severity
// will be dropped.
type simpleLogger struct {
	logger
	*loggerState
}

// SetLogLevel sets the severity level of the given logger.
func (logger *simpleLogger) SetLogLevel(level Level) {
	logger.setLevel(level)
}

// A callLogger represents the basic capability of a single logger.
// It has an associated logging level where messages of lesser severity
// will be dropped.
type callLogger struct {
	HasMinLevel
	st     *loggerState
	writer MinLevelWriter
}

// LogCallf logs a printf-formatted message at the given level.
// The location of the call is indicated by the calldepth argument.
// A calldepth of 1 means the function that called this function.
// A message will be discarded if level is less than the
// the effective log level of the logger.
// Note that the writers may also filter out messages that
// are less than their registered minimum severity level.
func (logger *callLogger) LogCallf(calldepth int, level Level, message string, args ...interface{}) {
	if logger == nil || !logger.willWrite(level) {
		return
	}
	loggerName := logger.st.name
	if loggerName == "" {
		loggerName = "<>"
	}
	rec := NewRecordf(calldepth+1, level, loggerName, message, args...)
	if logger.writer != nil {
		logger.writer.Write(rec)
	}
}

func (logger callLogger) willWrite(level Level) bool {
	if logger.HasMinLevel == nil {
		if !IsLevelEnabled(logger.st, level) {
			return false
		}
	} else if !IsLevelEnabled(logger, level) {
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

// loggerState holds the raw info about the logger.
//
// A nil loggerState pointer represents a read-only root logger.
type loggerState struct {
	name         string
	level        Level
	defaultLevel Level
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
