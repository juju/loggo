// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
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

// defaultName is the name of a writer that is registered
// by default that writes to stderr.
const defaultName = "default"

var (
	writerMutex sync.Mutex
	writers     = map[string]*registeredWriter{
		defaultName: &registeredWriter{
			writer: NewSimpleWriter(os.Stderr, &DefaultFormatter{}),
			level:  TRACE,
		},
	}
	globalMinLevel = TRACE
)

// ResetWriters puts the list of writers back into the initial state.
func ResetWriters() {
	writerMutex.Lock()
	defer writerMutex.Unlock()
	writers = map[string]*registeredWriter{
		"default": &registeredWriter{
			writer: NewSimpleWriter(os.Stderr, &DefaultFormatter{}),
			level:  TRACE,
		},
	}
	findMinLevel()
}

// ReplaceDefaultWriter is a convenience method that does the equivalent of
// RemoveWriter and then RegisterWriter with the name "default".  The previous
// default writer, if any is returned.
func ReplaceDefaultWriter(writer Writer) (Writer, error) {
	if writer == nil {
		return nil, fmt.Errorf("Writer cannot be nil")
	}
	writerMutex.Lock()
	defer writerMutex.Unlock()
	reg, found := writers[defaultName]
	if !found {
		return nil, fmt.Errorf("there is no %q writer", defaultName)
	}
	oldWriter := reg.writer
	reg.writer = writer
	return oldWriter, nil

}

// RegisterWriter adds the writer to the list of writers that get notified
// when logging.  When registering, the caller specifies the minimum logging
// level that will be written, and a name for the writer.  If there is already
// a registered writer with that name, an error is returned.
func RegisterWriter(name string, writer Writer, minLevel Level) error {
	if writer == nil {
		return fmt.Errorf("Writer cannot be nil")
	}
	writerMutex.Lock()
	defer writerMutex.Unlock()
	if _, found := writers[name]; found {
		return fmt.Errorf("there is already a Writer registered with the name %q", name)
	}
	writers[name] = &registeredWriter{writer: writer, level: minLevel}
	findMinLevel()
	return nil
}

// RemoveWriter removes the Writer identified by 'name' and returns it.
// If the Writer is not found, an error is returned.
func RemoveWriter(name string) (Writer, Level, error) {
	writerMutex.Lock()
	defer writerMutex.Unlock()
	registered, found := writers[name]
	if !found {
		return nil, UNSPECIFIED, fmt.Errorf("Writer %q is not registered", name)
	}
	delete(writers, name)
	findMinLevel()
	return registered.writer, registered.level, nil
}

func findMinLevel() {
	// We assume the lock is already held
	minLevel := CRITICAL
	for _, registered := range writers {
		if registered.level < minLevel {
			minLevel = registered.level
		}
	}
	globalMinLevel.set(minLevel)
}

// WillWrite returns whether there are any writers registered
// at or above the given severity level. If it returns
// false, a log message at the given level will be discarded.
func WillWrite(level Level) bool {
	return level >= globalMinLevel.get()
}

func writeToWriters(level Level, module, filename string, line int, timestamp time.Time, message string) {
	writerMutex.Lock()
	defer writerMutex.Unlock()
	for _, registered := range writers {
		if level >= registered.level {
			registered.writer.Write(level, module, filename, line, timestamp, message)
		}
	}
}
