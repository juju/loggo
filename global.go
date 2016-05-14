// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"os"
)

func defaultWriters() map[string]MinLevelWriter {
	return map[string]MinLevelWriter{
		defaultWriterName: NewMinLevelWriter(
			NewFormattingWriter(os.Stderr, nil),
			TRACE,
		),
	}
}

var (
	globalWriters = NewWriters(defaultWriters())
	globalLoggers = NewLoggers(WARNING, globalWriters)
)

// LoggerInfo returns information about the configured loggers and their
// logging levels. The information is returned in the format expected by
// ConfigureLoggers. Loggers with UNSPECIFIED level will not
// be included.
func LoggerInfo() string {
	return globalLoggers.Config().String()
}

// Root returns the root logger.
func Root() SubLogger {
	return globalLoggers.Root()
}

// GetLogger returns a logger for the given module name,
// creating it and its parents if necessary.
func GetLogger(name string) SubLogger {
	return globalLoggers.Get(name)
}

// ResetLogging iterates through the known modules and sets the levels of all
// to UNSPECIFIED, except for <root> which is set to WARNING.
func ResetLoggers() {
	globalLoggers.resetLevels()
}

// ResetWriters puts the list of writers back into the initial state.
func ResetWriters() {
	globalWriters.reset(defaultWriters())
}

// ReplaceDefaultWriter is a convenience method that does the equivalent of
// RemoveWriter and then RegisterWriter with the name "default".  The previous
// default writer, if any is returned.
func ReplaceDefaultWriter(writer Writer) (LegacyCompatibleWriter, error) {
	var w RecordWriter
	if writer != nil {
		w = &legacyAdaptingWriter{writer}
	}
	w, err := globalWriters.replace(defaultWriterName, w)
	if err != nil {
		return nil, err
	}
	return &LegacyWriterShim{w}, nil
}

// RegisterWriter adds the writer to the list of writers that get notified
// when logging.  When registering, the caller specifies the minimum logging
// level that will be written, and a name for the writer.  If there is already
// a registered writer with that name, an error is returned.
func RegisterWriter(name string, writer Writer, minLevel Level) error {
	if w, ok := writer.(RecordWriter); ok {
		return globalWriters.AddWithLevel(name, w, minLevel)
	}
	var w RecordWriter
	if writer != nil {
		w = &legacyAdaptingWriter{writer}
	}
	return globalWriters.AddWithLevel(name, w, minLevel)
}

// RemoveWriter removes the Writer identified by 'name' and returns it.
// If the Writer is not found, an error is returned.
func RemoveWriter(name string) (Writer, Level, error) {
	registered, err := globalWriters.remove(name)
	if err != nil {
		return nil, UNSPECIFIED, err
	}
	return &LegacyWriterShim{registered}, registered.MinLogLevel(), nil
}

// WillWrite returns whether there are any writers registered
// at or above the given severity level. If it returns
// false, a log message at the given level will be discarded.
func WillWrite(level Level) bool {
	return IsLevelEnabled(globalWriters, level)
}

// ConfigureLoggers configures loggers according to the given string
// specification, which specifies a set of modules and their associated
// logging levels.  Loggers are colon- or semicolon-separated; each
// module is specified as <modulename>=<level>.  White space outside of
// module names and levels is ignored.  The root module is specified
// with the name "<root>".
//
// An example specification:
//	`<root>=ERROR; foo.bar=WARNING`
func ConfigureLoggers(specification string) error {
	configs, err := ParseLoggersConfig(specification)
	if err != nil {
		return err
	}
	globalLoggers.ApplyConfig(configs)
	return nil
}
