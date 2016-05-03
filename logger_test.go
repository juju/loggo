// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"time"

	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
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

func (s *LoggerSuite) TestRootLogger(c *gc.C) {
	root := loggo.Logger{}
	c.Check(root.Name(), gc.Equals, "<root>")
	c.Assert(root.IsErrorEnabled(), gc.Equals, true)
	c.Assert(root.IsWarningEnabled(), gc.Equals, true)
	c.Assert(root.IsInfoEnabled(), gc.Equals, false)
	c.Assert(root.IsDebugEnabled(), gc.Equals, false)
	c.Assert(root.IsTraceEnabled(), gc.Equals, false)
}

func (s *LoggerSuite) TestModuleName(c *gc.C) {
	logger := loggo.GetLogger("loggo.testing")
	c.Assert(logger.Name(), gc.Equals, "loggo.testing")
}

func (s *LoggerSuite) TestSetLevel(c *gc.C) {
	logger := loggo.GetLogger("testing")

	c.Assert(logger.LogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.WARNING)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	logger.SetLogLevel(loggo.TRACE)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.TRACE)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.TRACE)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, true)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, true)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, true)
	logger.SetLogLevel(loggo.DEBUG)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.DEBUG)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.DEBUG)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, true)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, true)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	logger.SetLogLevel(loggo.INFO)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.INFO)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.INFO)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, true)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	logger.SetLogLevel(loggo.WARNING)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.WARNING)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.WARNING)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	logger.SetLogLevel(loggo.ERROR)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, false)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	// This is added for completeness, but not really expected to be used.
	logger.SetLogLevel(loggo.CRITICAL)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.CRITICAL)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.CRITICAL)
	c.Assert(logger.IsErrorEnabled(), gc.Equals, false)
	c.Assert(logger.IsWarningEnabled(), gc.Equals, false)
	c.Assert(logger.IsInfoEnabled(), gc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), gc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), gc.Equals, false)
	logger.SetLogLevel(loggo.UNSPECIFIED)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.WARNING)
}

func (s *LoggerSuite) TestModuleLowered(c *gc.C) {
	logger1 := loggo.GetLogger("TESTING.MODULE")
	logger2 := loggo.GetLogger("Testing")

	c.Assert(logger1.Name(), gc.Equals, "testing.module")
	c.Assert(logger2.Name(), gc.Equals, "testing")
}

func (s *LoggerSuite) TestLevelsInherited(c *gc.C) {
	root := loggo.GetLogger("")
	first := loggo.GetLogger("first")
	second := loggo.GetLogger("first.second")

	root.SetLogLevel(loggo.ERROR)
	c.Assert(root.LogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(root.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(first.LogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(first.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(second.LogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(second.EffectiveLogLevel(), gc.Equals, loggo.ERROR)

	first.SetLogLevel(loggo.DEBUG)
	c.Assert(root.LogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(root.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(first.LogLevel(), gc.Equals, loggo.DEBUG)
	c.Assert(first.EffectiveLogLevel(), gc.Equals, loggo.DEBUG)
	c.Assert(second.LogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(second.EffectiveLogLevel(), gc.Equals, loggo.DEBUG)

	second.SetLogLevel(loggo.INFO)
	c.Assert(root.LogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(root.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(first.LogLevel(), gc.Equals, loggo.DEBUG)
	c.Assert(first.EffectiveLogLevel(), gc.Equals, loggo.DEBUG)
	c.Assert(second.LogLevel(), gc.Equals, loggo.INFO)
	c.Assert(second.EffectiveLogLevel(), gc.Equals, loggo.INFO)

	first.SetLogLevel(loggo.UNSPECIFIED)
	c.Assert(root.LogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(root.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(first.LogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(first.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Assert(second.LogLevel(), gc.Equals, loggo.INFO)
	c.Assert(second.EffectiveLogLevel(), gc.Equals, loggo.INFO)
}

func (s *LoggerSuite) TestLoggingStrings(c *gc.C) {
	logger, writer := newTraceLogger("test")

	logger.Infof("simple")
	checkLastMessage(c, writer, "simple")

	logger.Infof("with args %d", 42)
	checkLastMessage(c, writer, "with args 42")

	logger.Infof("working 100%")
	checkLastMessage(c, writer, "working 100%")

	logger.Infof("missing %s")
	checkLastMessage(c, writer, "missing %s")
}

func (s *LoggerSuite) TestLocationCapture(c *gc.C) {
	logger, writer := newTraceLogger("test")

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
	logger, writer := newTraceLogger("test.writer")

	logger.Logf(loggo.UNSPECIFIED, "message")
	c.Assert(writer.Log(), gc.HasLen, 0)

	logger.Logf(loggo.Level(42), "message")
	c.Assert(writer.Log(), gc.HasLen, 0)

	logger.Logf(loggo.CRITICAL+loggo.Level(1), "message")
	c.Assert(writer.Log(), gc.HasLen, 0)
}

func (s *LoggerSuite) TestMessageFormatting(c *gc.C) {
	logger, writer := newTraceLogger("test.writer")

	logger.Logf(loggo.INFO, "some %s included", "formatting")

	log := writer.Log()
	c.Assert(log, gc.HasLen, 1)
	c.Assert(log[0].Message, gc.Equals, "some formatting included")
	c.Assert(log[0].Level, gc.Equals, loggo.INFO)
}

func (s *LoggerSuite) TestNoWriters(c *gc.C) {
	logger, writer := newTraceLogger("test.writer")
	loggo.RemoveWriter("default")

	logger.Warningf("just a simple warning")

	c.Check(writer.Log(), gc.HasLen, 0)
}

func (s *LoggerSuite) TestWritingLimitWarning(c *gc.C) {
	logger, writer := newTraceLogger("test.writer")
	loggo.RemoveWriter("default")
	err := loggo.RegisterWriter("test", writer, loggo.WARNING)
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
	c.Assert(log[0].Timestamp, Between(start, end))

	c.Assert(log[1].Level, gc.Equals, loggo.ERROR)
	c.Assert(log[1].Message, gc.Equals, "An error.")
	c.Assert(log[1].Timestamp, Between(start, end))

	c.Assert(log[2].Level, gc.Equals, loggo.WARNING)
	c.Assert(log[2].Message, gc.Equals, "A warning message")
	c.Assert(log[2].Timestamp, Between(start, end))
}

func (s *LoggerSuite) TestWritingLimitTrace(c *gc.C) {
	logger, writer := newTraceLogger("test.writer")
	loggo.RemoveWriter("default")
	err := loggo.RegisterWriter("test", writer, loggo.TRACE)
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
	c.Assert(log[0].Timestamp, Between(start, end))

	c.Assert(log[1].Level, gc.Equals, loggo.ERROR)
	c.Assert(log[1].Message, gc.Equals, "An error.")
	c.Assert(log[1].Timestamp, Between(start, end))

	c.Assert(log[2].Level, gc.Equals, loggo.WARNING)
	c.Assert(log[2].Message, gc.Equals, "A warning message")
	c.Assert(log[2].Timestamp, Between(start, end))

	c.Assert(log[3].Level, gc.Equals, loggo.INFO)
	c.Assert(log[3].Message, gc.Equals, "Info message")
	c.Assert(log[3].Timestamp, Between(start, end))

	c.Assert(log[4].Level, gc.Equals, loggo.TRACE)
	c.Assert(log[4].Message, gc.Equals, "Trace the function")
	c.Assert(log[4].Timestamp, Between(start, end))
}

func (s *LoggerSuite) TestMultipleWriters(c *gc.C) {
	logger, _ := newTraceLogger("test.writer")
	loggo.RemoveWriter("default")
	errorWriter := &loggo.TestWriter{}
	err := loggo.RegisterWriter("error", errorWriter, loggo.ERROR)
	c.Assert(err, gc.IsNil)
	warningWriter := &loggo.TestWriter{}
	err = loggo.RegisterWriter("warning", warningWriter, loggo.WARNING)
	c.Assert(err, gc.IsNil)
	infoWriter := &loggo.TestWriter{}
	err = loggo.RegisterWriter("info", infoWriter, loggo.INFO)
	c.Assert(err, gc.IsNil)
	traceWriter := &loggo.TestWriter{}
	err = loggo.RegisterWriter("trace", traceWriter, loggo.TRACE)
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
