// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggotest

import (
	"path"
	"sync"

	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
)

// Writer is a useful loggo.Writer for testing purposes. Each component
// of the logging message is stored in the Log array.
type Writer struct {
	mu  sync.Mutex
	log []loggo.Record
}

// Write saves the params as members in the LogValues struct appended to the Log array.
func (writer *Writer) Write(rec loggo.Record) {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	rec.Filename = path.Base(rec.Filename)
	writer.log = append(writer.log, rec)
}

// Clear removes any saved log messages.
func (writer *Writer) Clear() {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	writer.log = nil
}

// Log returns a copy of the current logged values.
func (writer *Writer) Log() []loggo.Record {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	v := make([]loggo.Record, len(writer.log))
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
