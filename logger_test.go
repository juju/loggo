// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"time"

	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggotest"
)

type LoggerSuite struct{}

var _ = gc.Suite(&LoggerSuite{})

func (*LoggerSuite) SetUpTest(c *gc.C) {
	loggo.ResetLoggers()
}

func (s *LoggerSuite) TearDownTest(c *gc.C) {
	loggo.ResetWriters()
}

func (s *LoggerSuite) TearDownSuite(c *gc.C) {
	loggo.ResetLoggers()
}

func (s *LoggerSuite) TestSetLevel(c *gc.C) {
	logger := loggo.NewLogger(nil)

	c.Assert(logger.MinLogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(loggo.EffectiveMinLevel(logger), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.ERROR), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.WARNING), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.INFO), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.DEBUG), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.TRACE), gc.Equals, true)
	logger.SetLogLevel(loggo.TRACE)
	c.Assert(logger.MinLogLevel(), gc.Equals, loggo.TRACE)
	c.Assert(loggo.EffectiveMinLevel(logger), gc.Equals, loggo.TRACE)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.ERROR), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.WARNING), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.INFO), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.DEBUG), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.TRACE), gc.Equals, true)
	logger.SetLogLevel(loggo.DEBUG)
	c.Assert(logger.MinLogLevel(), gc.Equals, loggo.DEBUG)
	c.Assert(loggo.EffectiveMinLevel(logger), gc.Equals, loggo.DEBUG)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.ERROR), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.WARNING), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.INFO), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.DEBUG), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.TRACE), gc.Equals, false)
	logger.SetLogLevel(loggo.INFO)
	c.Assert(logger.MinLogLevel(), gc.Equals, loggo.INFO)
	c.Assert(loggo.EffectiveMinLevel(logger), gc.Equals, loggo.INFO)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.ERROR), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.WARNING), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.INFO), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.DEBUG), gc.Equals, false)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.TRACE), gc.Equals, false)
	logger.SetLogLevel(loggo.WARNING)
	c.Assert(logger.MinLogLevel(), gc.Equals, loggo.WARNING)
	c.Assert(loggo.EffectiveMinLevel(logger), gc.Equals, loggo.WARNING)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.ERROR), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.WARNING), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.INFO), gc.Equals, false)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.DEBUG), gc.Equals, false)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.TRACE), gc.Equals, false)
	logger.SetLogLevel(loggo.ERROR)
	c.Assert(logger.MinLogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(loggo.EffectiveMinLevel(logger), gc.Equals, loggo.ERROR)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.ERROR), gc.Equals, true)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.WARNING), gc.Equals, false)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.INFO), gc.Equals, false)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.DEBUG), gc.Equals, false)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.TRACE), gc.Equals, false)
	// This is added for completeness, but not really expected to be used.
	logger.SetLogLevel(loggo.CRITICAL)
	c.Assert(logger.MinLogLevel(), gc.Equals, loggo.CRITICAL)
	c.Assert(loggo.EffectiveMinLevel(logger), gc.Equals, loggo.CRITICAL)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.ERROR), gc.Equals, false)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.WARNING), gc.Equals, false)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.INFO), gc.Equals, false)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.DEBUG), gc.Equals, false)
	c.Assert(loggo.IsLevelEnabled(logger, loggo.TRACE), gc.Equals, false)
	logger.SetLogLevel(loggo.UNSPECIFIED)
	c.Assert(logger.MinLogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(loggo.EffectiveMinLevel(logger), gc.Equals, loggo.UNSPECIFIED)
}

func (s *LoggerSuite) TestLoggingStrings(c *gc.C) {
	logger, writer := loggotest.TraceLogger()

	logger.Infof("simple")
	loggotest.CheckLastMessage(c, writer, "simple")

	logger.Infof("with args %d", 42)
	loggotest.CheckLastMessage(c, writer, "with args 42")

	logger.Infof("working 100%")
	loggotest.CheckLastMessage(c, writer, "working 100%")

	logger.Infof("missing %s")
	loggotest.CheckLastMessage(c, writer, "missing %s")
}

func (s *LoggerSuite) TestLoggingLimitWarning(c *gc.C) {
	logger, writer := loggotest.TraceLogger()
	logger.SetLogLevel(loggo.WARNING)

	start := time.Now()
	logger.Criticalf("Something critical.")
	logger.Errorf("An error.")
	logger.Warningf("A warning message")
	logger.Infof("Info message")
	logger.Tracef("Trace the function")
	end := time.Now()

	log := writer.Log()
	c.Assert(log, gc.HasLen, 3)
	c.Assert(log[0].Level, gc.Equals, loggo.CRITICAL)
	c.Assert(log[0].Message, gc.Equals, "Something critical.")
	c.Assert(log[0].Timestamp, loggotest.Between(start, end))

	c.Assert(log[1].Level, gc.Equals, loggo.ERROR)
	c.Assert(log[1].Message, gc.Equals, "An error.")
	c.Assert(log[1].Timestamp, loggotest.Between(start, end))

	c.Assert(log[2].Level, gc.Equals, loggo.WARNING)
	c.Assert(log[2].Message, gc.Equals, "A warning message")
	c.Assert(log[2].Timestamp, loggotest.Between(start, end))
}

func (s *LoggerSuite) TestLoggingLimitTrace(c *gc.C) {
	logger, writer := loggotest.TraceLogger()
	logger.SetLogLevel(loggo.TRACE)

	start := time.Now()
	logger.Criticalf("Something critical.")
	logger.Errorf("An error.")
	logger.Warningf("A warning message")
	logger.Infof("Info message")
	logger.Tracef("Trace the function")
	end := time.Now()

	log := writer.Log()
	c.Assert(log, gc.HasLen, 5)
	c.Assert(log[0].Level, gc.Equals, loggo.CRITICAL)
	c.Assert(log[0].Message, gc.Equals, "Something critical.")
	c.Assert(log[0].Timestamp, loggotest.Between(start, end))

	c.Assert(log[1].Level, gc.Equals, loggo.ERROR)
	c.Assert(log[1].Message, gc.Equals, "An error.")
	c.Assert(log[1].Timestamp, loggotest.Between(start, end))

	c.Assert(log[2].Level, gc.Equals, loggo.WARNING)
	c.Assert(log[2].Message, gc.Equals, "A warning message")
	c.Assert(log[2].Timestamp, loggotest.Between(start, end))

	c.Assert(log[3].Level, gc.Equals, loggo.INFO)
	c.Assert(log[3].Message, gc.Equals, "Info message")
	c.Assert(log[3].Timestamp, loggotest.Between(start, end))

	c.Assert(log[4].Level, gc.Equals, loggo.TRACE)
	c.Assert(log[4].Message, gc.Equals, "Trace the function")
	c.Assert(log[4].Timestamp, loggotest.Between(start, end))
}

func (s *LoggerSuite) TestLocationCapture(c *gc.C) {
	logger, writer := loggotest.TraceLogger()

	logger.Criticalf("critical message") //tag critical-location
	logger.Errorf("error message")       //tag error-location
	logger.Warningf("warning message")   //tag warning-location
	logger.Infof("info message")         //tag info-location
	logger.Debugf("debug message")       //tag debug-location
	logger.Tracef("trace message")       //tag trace-location

	log := writer.Log()
	tags := []string{
		"critical-location",
		"error-location",
		"warning-location",
		"info-location",
		"debug-location",
		"trace-location",
	}
	c.Assert(log, gc.HasLen, len(tags))
	for x := range tags {
		assertLocation(c, log[x], tags[x])
	}
}

func (s *LoggerSuite) TestLogDoesntLogWeirdLevels(c *gc.C) {
	logger, writer := loggotest.TraceLogger()

	logger.Logf(loggo.UNSPECIFIED, "message")
	c.Assert(writer.Log(), gc.HasLen, 0)

	logger.Logf(loggo.Level(42), "message")
	c.Assert(writer.Log(), gc.HasLen, 0)

	logger.Logf(loggo.CRITICAL+loggo.Level(1), "message")
	c.Assert(writer.Log(), gc.HasLen, 0)
}

func (s *LoggerSuite) TestMessageFormatting(c *gc.C) {
	logger, writer := loggotest.TraceLogger()

	logger.Logf(loggo.INFO, "some %s included", "formatting")

	log := writer.Log()
	c.Assert(log, gc.HasLen, 1)
	c.Assert(log[0].Message, gc.Equals, "some formatting included")
	c.Assert(log[0].Level, gc.Equals, loggo.INFO)
}

func (s *LoggerSuite) TestNoWriters(c *gc.C) {
	writer := &loggotest.Writer{}
	loggo.RemoveWriter("default")
	err := loggo.RegisterWriter("test", writer, loggo.TRACE)
	c.Assert(err, gc.IsNil)
	// Use a non-global logger with no writers set.
	logger := loggo.NewLogger(nil)
	logger.SetLogLevel(loggo.TRACE)

	logger.Warningf("just a simple warning")

	c.Check(writer.Log(), gc.HasLen, 0)
}

func (s *LoggerSuite) TestWritingLimitWarning(c *gc.C) {
	writers := loggo.NewWriters(nil)
	logger := loggo.NewLogger(writers)
	logger.SetLogLevel(loggo.TRACE)
	writer := &loggotest.Writer{}
	err := writers.AddWithLevel("test", writer, loggo.WARNING)
	c.Assert(err, gc.IsNil)

	start := time.Now()
	logger.Criticalf("Something critical.")
	logger.Errorf("An error.")
	logger.Warningf("A warning message")
	logger.Infof("Info message")
	logger.Tracef("Trace the function")
	end := time.Now()

	log := writer.Log()
	c.Assert(log, gc.HasLen, 3)
	c.Assert(log[0].Level, gc.Equals, loggo.CRITICAL)
	c.Assert(log[0].Message, gc.Equals, "Something critical.")
	c.Assert(log[0].Timestamp, loggotest.Between(start, end))

	c.Assert(log[1].Level, gc.Equals, loggo.ERROR)
	c.Assert(log[1].Message, gc.Equals, "An error.")
	c.Assert(log[1].Timestamp, loggotest.Between(start, end))

	c.Assert(log[2].Level, gc.Equals, loggo.WARNING)
	c.Assert(log[2].Message, gc.Equals, "A warning message")
	c.Assert(log[2].Timestamp, loggotest.Between(start, end))
}

func (s *LoggerSuite) TestWritingLimitTrace(c *gc.C) {
	writers := loggo.NewWriters(nil)
	logger := loggo.NewLogger(writers)
	logger.SetLogLevel(loggo.TRACE)
	writer := &loggotest.Writer{}
	err := writers.AddWithLevel("test", writer, loggo.TRACE)
	c.Assert(err, gc.IsNil)

	start := time.Now()
	logger.Criticalf("Something critical.")
	logger.Errorf("An error.")
	logger.Warningf("A warning message")
	logger.Infof("Info message")
	logger.Tracef("Trace the function")
	end := time.Now()

	log := writer.Log()
	c.Assert(log, gc.HasLen, 5)
	c.Assert(log[0].Level, gc.Equals, loggo.CRITICAL)
	c.Assert(log[0].Message, gc.Equals, "Something critical.")
	c.Assert(log[0].Timestamp, loggotest.Between(start, end))

	c.Assert(log[1].Level, gc.Equals, loggo.ERROR)
	c.Assert(log[1].Message, gc.Equals, "An error.")
	c.Assert(log[1].Timestamp, loggotest.Between(start, end))

	c.Assert(log[2].Level, gc.Equals, loggo.WARNING)
	c.Assert(log[2].Message, gc.Equals, "A warning message")
	c.Assert(log[2].Timestamp, loggotest.Between(start, end))

	c.Assert(log[3].Level, gc.Equals, loggo.INFO)
	c.Assert(log[3].Message, gc.Equals, "Info message")
	c.Assert(log[3].Timestamp, loggotest.Between(start, end))

	c.Assert(log[4].Level, gc.Equals, loggo.TRACE)
	c.Assert(log[4].Message, gc.Equals, "Trace the function")
	c.Assert(log[4].Timestamp, loggotest.Between(start, end))
}

func (s *LoggerSuite) TestMultipleWriters(c *gc.C) {
	writers := loggo.NewWriters(nil)
	logger := loggo.NewLogger(writers)
	logger.SetLogLevel(loggo.TRACE)
	errorWriter := &loggotest.Writer{}
	err := writers.AddWithLevel("error", errorWriter, loggo.ERROR)
	c.Assert(err, gc.IsNil)
	warningWriter := &loggotest.Writer{}
	err = writers.AddWithLevel("warning", warningWriter, loggo.WARNING)
	c.Assert(err, gc.IsNil)
	infoWriter := &loggotest.Writer{}
	err = writers.AddWithLevel("info", infoWriter, loggo.INFO)
	c.Assert(err, gc.IsNil)
	traceWriter := &loggotest.Writer{}
	err = writers.AddWithLevel("trace", traceWriter, loggo.TRACE)
	c.Assert(err, gc.IsNil)

	logger.Errorf("An error.")
	logger.Warningf("A warning message")
	logger.Infof("Info message")
	logger.Tracef("Trace the function")

	c.Assert(errorWriter.Log(), gc.HasLen, 1)
	c.Assert(warningWriter.Log(), gc.HasLen, 2)
	c.Assert(infoWriter.Log(), gc.HasLen, 3)
	c.Assert(traceWriter.Log(), gc.HasLen, 4)
}

func (s *LoggerSuite) TestLoggerAsIOWriter(c *gc.C) {
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
