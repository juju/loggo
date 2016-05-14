// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggotest"
)

type LoggerHelpersSuite struct{}

var _ = gc.Suite(&LoggerHelpersSuite{})

func (*LoggerHelpersSuite) SetUpTest(c *gc.C) {
	loggo.ResetLoggers()
}

func (s *LoggerHelpersSuite) TearDownTest(c *gc.C) {
	loggo.ResetWriters()
}

func (s *LoggerHelpersSuite) TearDownSuite(c *gc.C) {
	loggo.ResetLoggers()
}

func (s *LoggerHelpersSuite) TestLoggerAsGoLogger(c *gc.C) {
	writer := &loggotest.Writer{}
	logger := loggo.NewLogger(loggo.NewMinLevelWriter(writer, loggo.TRACE))
	logger.SetLogLevel(loggo.TRACE)
	log := loggo.LoggerAsGoLogger(logger, loggo.WARNING)

	log.Print("raw message")

	records := writer.Log()
	c.Assert(records, gc.HasLen, 1)
	c.Assert(records[0].Level, gc.Equals, loggo.WARNING)
	c.Assert(records[0].LoggerName, gc.Equals, "<>")
	c.Assert(records[0].Filename, gc.Equals, "logger_helpers_test.go")
	c.Assert(records[0].Message, gc.Equals, "raw message")
}

func (s *LoggerHelpersSuite) TestLoggerAsIOWriter(c *gc.C) {
	writer := &loggotest.Writer{}
	logger := loggo.NewLogger(loggo.NewMinLevelWriter(writer, loggo.TRACE))
	logger.SetLogLevel(loggo.TRACE)
	ioWriter := loggo.LoggerAsIOWriter(logger, loggo.WARNING)

	logger.Errorf("Error message")
	_, err := ioWriter.Write([]byte("raw message"))
	c.Assert(err, gc.IsNil)
	logger.Infof("Info message")
	logger.Warningf("Warning message")

	log := writer.Log()
	c.Assert(log, gc.HasLen, 4)
	c.Assert(log[0].Level, gc.Equals, loggo.ERROR)
	c.Assert(log[0].Message, gc.Equals, "Error message")
	c.Assert(log[1].Level, gc.Equals, loggo.WARNING)
	c.Assert(log[1].Message, gc.Equals, "raw message")
	c.Assert(log[2].Level, gc.Equals, loggo.INFO)
	c.Assert(log[2].Message, gc.Equals, "Info message")
	c.Assert(log[3].Level, gc.Equals, loggo.WARNING)
	c.Assert(log[3].Message, gc.Equals, "Warning message")
}
