// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"os"
	"strings"
	"sync"
)

// defaultName is the name of a writer that is registered
// by default that writes to stderr.
const defaultName = "default"

func defaultWriters() map[string]*minLevelWriter {
	return map[string]*minLevelWriter{
		defaultName: &minLevelWriter{
			writer: NewSimpleWriter(os.Stderr, &DefaultFormatter{}),
			level:  TRACE,
		},
	}
}

var (
	globalWriters = newWriters(defaultWriters())
)

// Initially the modules map only contains the root module.
var (
	modulesMutex sync.Mutex
	modules      = map[string]*module{
		"": root,
	}
)

// LoggerInfo returns information about the configured loggers and their
// logging levels. The information is returned in the format expected by
// ConfigureLoggers. Loggers with UNSPECIFIED level will not
// be included.
func LoggerInfo() string {
	modulesMutex.Lock()
	defer modulesMutex.Unlock()

	return loggerInfo(modules)
}

// GetLogger returns a Logger for the given module name,
// creating it and its parents if necessary.
func GetLogger(name string) Logger {
	// Lowercase the module name, and look for it in the modules map.
	name = strings.ToLower(name)
	modulesMutex.Lock()
	defer modulesMutex.Unlock()
	return getLoggerInternal(name)
}

// getLoggerInternal assumes that the modulesMutex is locked.
func getLoggerInternal(name string) Logger {
	impl, found := modules[name]
	if found {
		return Logger{impl}
	}
	parentName := ""
	if i := strings.LastIndex(name, "."); i >= 0 {
		parentName = name[0:i]
	}
	parent := getLoggerInternal(parentName)
	logger := newLogger(name, parent.impl)
	modules[name] = logger.impl
	return logger
}

// ResetLogging iterates through the known modules and sets the levels of all
// to UNSPECIFIED, except for <root> which is set to WARNING.
func ResetLoggers() {
	modulesMutex.Lock()
	defer modulesMutex.Unlock()
	for name, module := range modules {
		if name == "" {
			module.level.set(WARNING)
		} else {
			module.level.set(UNSPECIFIED)
		}
	}
}

// ResetWriters puts the list of writers back into the initial state.
func ResetWriters() {
	globalWriters.reset(defaultWriters())
}

// ReplaceDefaultWriter is a convenience method that does the equivalent of
// RemoveWriter and then RegisterWriter with the name "default".  The previous
// default writer, if any is returned.
func ReplaceDefaultWriter(writer Writer) (Writer, error) {
	return globalWriters.replace(defaultName, writer)
}

// RegisterWriter adds the writer to the list of writers that get notified
// when logging.  When registering, the caller specifies the minimum logging
// level that will be written, and a name for the writer.  If there is already
// a registered writer with that name, an error is returned.
func RegisterWriter(name string, writer Writer, minLevel Level) error {
	return globalWriters.addWithLevel(name, writer, minLevel)
}

// RemoveWriter removes the Writer identified by 'name' and returns it.
// If the Writer is not found, an error is returned.
func RemoveWriter(name string) (Writer, Level, error) {
	registered, err := globalWriters.remove(name)
	if err != nil {
		return nil, UNSPECIFIED, err
	}
	return registered.writer, registered.level, nil
}

// WillWrite returns whether there are any writers registered
// at or above the given severity level. If it returns
// false, a log message at the given level will be discarded.
func WillWrite(level Level) bool {
	return globalWriters.willWrite(level)
}
