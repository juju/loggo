// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"path"
	"sync"
	"time"
)

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
	writer.log = append(writer.log,
		TestLogValues{level, module, path.Base(filename), line, timestamp, message})
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

// TestFormatter is a useful Writer for testing purposes. Each component
// of the logging message is stored in the Log array.
type TestFormatter struct {
	TestWriter

	format func(level Level, module, filename string, line int, timestamp time.Time, message string) string
}

// NewTestFormatter returns a new TestFormatter that wraps the given
// format func. If the func is nil then Format() will return the message.
func NewTestFormatter(format func(level Level, module, filename string, line int, timestamp time.Time, message string) string) *TestFormatter {
	return &TestFormatter{
		format: format,
	}
}

// Format saves the params as members in the TestLogValues struct appended to the Log array.
func (tf *TestFormatter) Format(level Level, module, filename string, line int, timestamp time.Time, message string) string {
	tf.Write(level, module, filename, line, timestamp, message)
	if tf.format == nil {
		return message
	}
	return tf.format(level, module, filename, line, timestamp, message)
}
