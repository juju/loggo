// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggotest

import (
	"path"
	"sync"
	"time"

	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
)

// TODO(ericsnow) LogValues can go away once we have a Record type.

// LogValues represents a single logging call.
type LogValues struct {
	Level     loggo.Level
	Module    string
	Filename  string
	Line      int
	Timestamp time.Time
	Message   string
}

// Writer is a useful loggo.Writer for testing purposes. Each component
// of the logging message is stored in the Log array.
type Writer struct {
	mu  sync.Mutex
	log []LogValues
}

// Write saves the params as members in the LogValues struct appended to the Log array.
func (writer *Writer) Write(level loggo.Level, module, filename string, line int, timestamp time.Time, message string) {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	writer.log = append(writer.log, LogValues{
		Level:     level,
		Module:    module,
		Filename:  path.Base(filename),
		Line:      line,
		Timestamp: timestamp,
		Message:   message,
	})
}

// Clear removes any saved log messages.
func (writer *Writer) Clear() {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	writer.log = nil
}

// Log returns a copy of the current logged values.
func (writer *Writer) Log() []LogValues {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	v := make([]LogValues, len(writer.log))
	copy(v, writer.log)
	return v
}

func CheckLastMessage(c *gc.C, writer *Writer, expected string) {
	log := writer.Log()
	writer.Clear()
	if c.Check(len(log) > 0, gc.Equals, true) {
		obtained := log[len(log)-1].Message
		c.Check(obtained, gc.Equals, expected)
	}
}
