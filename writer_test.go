// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"bytes"
	"time"

	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
)

type SimpleWriterSuite struct{}

var _ = gc.Suite(&SimpleWriterSuite{})

func (s *SimpleWriterSuite) TestNewSimpleWriter(c *gc.C) {
	now := time.Now()
	formatter := func(entry loggo.Entry) string {
		return "<< " + entry.Message + " >>"
	}
	buf := &bytes.Buffer{}

	writer := loggo.NewSimpleWriter(buf, formatter)
	writer.Write(loggo.Entry{
		Level:     loggo.INFO,
		Module:    "test",
		Filename:  "somefile.go",
		Line:      12,
		Timestamp: now,
		Message:   "a message",
	})

	c.Check(buf.String(), gc.Equals, "<< a message >>\n")
}
