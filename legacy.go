// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"io"
	"path"
	"sync"
	"time"
)

// TODO(ericsnow) The contents of this file can go away as soon as we
// stop worrying about breaking compatibility.

// ParseConfigurationString parses a logger configuration string into a map of
// logger names and their associated log level. This method is provided to
// allow other programs to pre-validate a configuration string rather than
// just calling ConfigureLoggers.
//
// Loggers are colon- or semicolon-separated; each module is specified as
// <modulename>=<level>.  White space outside of module names and levels is
// ignored.  The root module is specified with the name "<root>".
//
// As a special case, a log level may be specified on its own.
// This is equivalent to specifying the level of the root module,
// so "DEBUG" is equivalent to `<root>=DEBUG`
//
// An example specification:
//	`<root>=ERROR; foo.bar=WARNING`
//
// ParseConfigurationString is deprecated. Use ParseLoggersConfig instead.
func ParseConfigurationString(specification string) (map[string]Level, error) {
	configs, err := ParseLoggersConfig(specification)
	if err != nil {
		return nil, err
	}
	levels := make(map[string]Level)
	for name, cfg := range configs {
		levels[name] = cfg.Level
	}
	return levels, nil
}

// LegacyFormatter defines the single method Format, which takes the logging
// information, and converts it to a string.
//
// LegacyFormatter is deprecated. Use Formatter instead.
type LegacyFormatter interface {
	Format(level Level, loggerName, filename string, line int, timestamp time.Time, message string) string
}

type legacyAdaptingFormatter struct {
	legacy LegacyFormatter
}

func (f *legacyAdaptingFormatter) Format(rec Record) string {
	return f.legacy.Format(rec.Level, rec.LoggerName, rec.Filename, rec.Line, rec.Timestamp, rec.Message)
}

// DefaultFormatter provides a simple concatenation of all the components.
//
// DefaultFormatter is deprecated. Pass nil to NewFormattingWriter() instead.
type DefaultFormatter struct{}

// Format returns the parameters separated by spaces except for filename and
// line which are separated by a colon.  The timestamp is shown to second
// resolution in UTC.
func (*DefaultFormatter) Format(level Level, loggerName, filename string, line int, timestamp time.Time, message string) string {
	rec := Record{
		Level:      level,
		LoggerName: loggerName,
		Filename:   filename,
		Line:       line,
		Timestamp:  timestamp,
		Message:    message,
	}
	return rec.String()
}

// Writer is implemented by any recipient of log messages.
//
// Writer is deprecated. Use RecordWriter instead.
type Writer legacyWriter

type legacyWriter interface {
	// Write writes a message to the Writer with the given
	// level and logger name. The filename and line hold
	// the file name and line number of the code that is
	// generating the log message; the time stamp holds
	// the time the log message was generated, and
	// message holds the log message itself.
	Write(level Level, name, filename string, line int, timestamp time.Time, message string)
}

// LegacyCompatibleWriter is a shim to temporarily support both interfaces.
//
// LegacyCompatibleWriter only exists to support transitioning to RecordWriter.
type LegacyCompatibleWriter interface {
	RecordWriter
	legacyWriter
}

type LegacyWriterShim struct {
	RecordWriter
}

func (shim LegacyWriterShim) Write(level Level, loggerName, filename string, line int, timestamp time.Time, message string) {
	shim.WriteRecord(Record{
		Level:      level,
		LoggerName: loggerName,
		Filename:   filename,
		Line:       line,
		Timestamp:  timestamp,
		Message:    message,
	})
}

type legacyAdaptingWriter struct {
	legacyWriter
}

func (law *legacyAdaptingWriter) WriteRecord(rec Record) {
	law.Write(rec.Level, rec.LoggerName, rec.Filename, rec.Line, rec.Timestamp, rec.Message)
}

// NewSimpleWriter returns a new writer that writes
// log messages to the given io.Writer formatting the
// messages with the given formatter.
//
// NewSimpleWriter is deprecated. Use NewFormattingWriter.
func NewSimpleWriter(writer io.Writer, formatter LegacyFormatter) LegacyCompatibleWriter {
	return &LegacyWriterShim{
		&formattingWriter{writer, &legacyAdaptingFormatter{formatter}},
	}
}

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
//
// LogLevel is deprecated. Use MinLogLevel() instead.
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
//
// LogLevel is deprecated. Use loggo.EffectiveMinLevel() instead.
func (logger SubLogger) EffectiveLogLevel() Level {
	return EffectiveMinLevel(logger.getModule())
}

// IsLevelEnabled returns whether debugging is enabled
// for the given log level.
//
// IsLevelEnabled is deprecated. Use loggo.IsLevelEnabled() instead.
func (logger SubLogger) IsLevelEnabled(level Level) bool {
	return IsLevelEnabled(logger.getModule(), level)
}

// IsErrorEnabled returns whether debugging is enabled
// at error level.
//
// IsErrorEnabled is deprecated. Use loggo.IsLevelEnabled() instead.
func (logger SubLogger) IsErrorEnabled() bool {
	return logger.IsLevelEnabled(ERROR)
}

// IsWarningEnabled returns whether debugging is enabled
// at warning level.
//
// IsWarningEnabled is deprecated. Use loggo.IsLevelEnabled() instead.
func (logger SubLogger) IsWarningEnabled() bool {
	return logger.IsLevelEnabled(WARNING)
}

// IsInfoEnabled returns whether debugging is enabled
// at info level.
//
// IsInfoEnabled is deprecated. Use loggo.IsLevelEnabled() instead.
func (logger SubLogger) IsInfoEnabled() bool {
	return logger.IsLevelEnabled(INFO)
}

// IsDebugEnabled returns whether debugging is enabled
// at debug level.
//
// IsDebugEnabled is deprecated. Use loggo.IsLevelEnabled() instead.
func (logger SubLogger) IsDebugEnabled() bool {
	return logger.IsLevelEnabled(DEBUG)
}

// IsTraceEnabled returns whether debugging is enabled
// at trace level.
//
// IsTraceEnabled is deprecated. Use loggo.IsLevelEnabled() instead.
func (logger SubLogger) IsTraceEnabled() bool {
	return logger.IsLevelEnabled(TRACE)
}

// TestLogValues represents a single logging call.
type TestLogValues struct {
	Level     Level
	Module    string
	Filename  string
	Line      int
	Timestamp time.Time
	Message   string
}

// TestWriter is a useful Writer for testing purposes.  Each component of the
// logging message is stored in the Log array.
type TestWriter struct {
	mu  sync.Mutex
	log []TestLogValues
}

// Write saves the params as members in the TestLogValues struct appended to the Log array.
func (writer *TestWriter) Write(level Level, module, filename string, line int, timestamp time.Time, message string) {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	writer.log = append(writer.log, TestLogValues{
		Level:     level,
		Module:    module,
		Filename:  path.Base(filename),
		Line:      line,
		Timestamp: timestamp,
		Message:   message,
	})
}

// Clear removes any saved log messages.
func (writer *TestWriter) Clear() {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	writer.log = nil
}

// Log returns a copy of the current logged values.
func (writer *TestWriter) Log() []TestLogValues {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	v := make([]TestLogValues, len(writer.log))
	copy(v, writer.log)
	return v
}

// TODO(ericsnow) ...also, the Write() method of loggotest.Writer.
