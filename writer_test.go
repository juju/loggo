// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"bytes"
	"time"

	gc "gopkg.in/check.v1"
)

type SimpleWriterSuite struct{}

var _ = gc.Suite(&SimpleWriterSuite{})

func (s *SimpleWriterSuite) TestNewSimpleWriter(c *gc.C) {
	now := time.Now()
	formatter := func(entry Entry) string {
		return "<< " + entry.Message + " >>"
	}
	buf := &bytes.Buffer{}

	writer := NewSimpleWriter(buf, formatter)
	writer.Write(Entry{
		Level:     INFO,
		Module:    "test",
		Filename:  "somefile.go",
		Line:      12,
		Timestamp: now,
		Message:   "a message",
		Labels:    nil,
	})

	c.Check(buf.String(), gc.Equals, "<< a message >>\n")
}

func (s *SimpleWriterSuite) TestNewSimpleWriterWithLabels(c *gc.C) {
	now := time.Now()
	formatter := func(entry Entry) string {
		return "<< " + entry.Message + " >>"
	}
	buf := &bytes.Buffer{}

	writer := NewSimpleWriter(buf, formatter)
	writer.Write(Entry{
		Level:     INFO,
		Module:    "test",
		Filename:  "somefile.go",
		Line:      12,
		Timestamp: now,
		Message:   "a message",
		Labels:    Labels{LoggerTags: "ONE,TWO"},
	})

	c.Check(buf.String(), gc.Equals, "<< a message >>\n")
}
