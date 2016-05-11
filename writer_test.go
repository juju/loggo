// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"bytes"
	"time"

	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggotest"
)

type SimpleWriterSuite struct{}

var _ = gc.Suite(&SimpleWriterSuite{})

func (s *SimpleWriterSuite) TestNewSimpleWriter(c *gc.C) {
	now := time.Now()
	formatter := loggotest.NewFormatter(func(level loggo.Level, module, filename string, line int, timestamp time.Time, message string) string {
		return "<< " + message + " >>"
	})
	buf := &bytes.Buffer{}

	writer := loggo.NewSimpleWriter(buf, formatter)
	writer.Write(loggo.INFO, "test", "somefile.go", 12, now, "a message")

	log := formatter.Log()
	c.Check(log, gc.DeepEquals, []loggotest.LogValues{{
		Level:     loggo.INFO,
		Module:    "test",
		Filename:  "somefile.go",
		Line:      12,
		Timestamp: now,
		Message:   "a message",
	}})
	c.Check(buf.String(), gc.Equals, "<< a message >>\n")
}
