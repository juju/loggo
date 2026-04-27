// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/juju/tc"
)

type SimpleWriterSuite struct{}

func TestSimpleWriterSuite(t *testing.T) {
	tc.Run(t, &SimpleWriterSuite{})
}

func (s *SimpleWriterSuite) TestNewSimpleWriter(c *tc.C) {
	now := time.Now()
	formatter := func(entry Entry) string {
		return "<< " + entry.Message + " >>"
	}
	buf := &bytes.Buffer{}

	writer := NewSimpleWriter(buf, formatter)
	_ = writer.Write(context.Background(), Entry{
		Level:     INFO,
		Module:    "test",
		Filename:  "somefile.go",
		Line:      12,
		Timestamp: now,
		Message:   "a message",
		Labels:    nil,
	})

	c.Check(buf.String(), tc.Equals, "<< a message >>\n")
}

func (s *SimpleWriterSuite) TestNewSimpleWriterWithLabels(c *tc.C) {
	now := time.Now()
	formatter := func(entry Entry) string {
		return "<< " + entry.Message + " >>"
	}
	buf := &bytes.Buffer{}

	writer := NewSimpleWriter(buf, formatter)
	_ = writer.Write(context.Background(), Entry{
		Level:     INFO,
		Module:    "test",
		Filename:  "somefile.go",
		Line:      12,
		Timestamp: now,
		Message:   "a message",
		Labels:    Labels{LoggerTags: "ONE,TWO"},
	})

	c.Check(buf.String(), tc.Equals, "<< a message >>\n")
}
