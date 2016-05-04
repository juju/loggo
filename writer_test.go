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
	formatter := loggotest.NewFormatter(func(rec loggo.Record) string {
		return "<< " + rec.Message + " >>"
	})
	buf := &bytes.Buffer{}
	rec := loggo.Record{
		Level:      loggo.INFO,
		LoggerName: "test",
		Filename:   "somefile.go",
		Line:       12,
		Timestamp:  now,
		Message:    "a message",
	}

	writer := loggo.NewFormattingWriter(buf, formatter)
	writer.WriteRecord(rec)

	log := formatter.Log()
	c.Check(log, gc.DeepEquals, []loggo.Record{rec})
	c.Check(buf.String(), gc.Equals, "<< a message >>\n")
}
