// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"io/ioutil"
	"os"

	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
)

type loggerSuite struct{}

var _ = gc.Suite(&loggerSuite{})

func (*loggerSuite) SetUpTest(c *gc.C) {
	loggo.ResetLoggers()
}

func (*loggerSuite) TestRootLogger(c *gc.C) {
	root := loggo.Logger{}
	c.Assert(root.Name(), gc.Equals, "<root>")
	c.Assert(root.IsErrorEnabled(), gc.Equals, true)
	c.Assert(root.IsWarningEnabled(), gc.Equals, true)
	c.Assert(root.IsInfoEnabled(), gc.Equals, false)
	c.Assert(root.IsDebugEnabled(), gc.Equals, false)
	c.Assert(root.IsTraceEnabled(), gc.Equals, false)
}

func (*loggerSuite) TestModuleName(c *gc.C) {
	logger := loggo.GetLogger("loggo.testing")
	c.Assert(logger.Name(), gc.Equals, "loggo.testing")
}

func (*loggerSuite) TestSetLevel(c *gc.C) {
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

func (*loggerSuite) TestLevelsSharedForSameModule(c *gc.C) {
	logger1 := loggo.GetLogger("testing.module")
	logger2 := loggo.GetLogger("testing.module")

	logger1.SetLogLevel(loggo.INFO)
	c.Assert(logger1.IsInfoEnabled(), gc.Equals, true)
	c.Assert(logger2.IsInfoEnabled(), gc.Equals, true)
}

func (*loggerSuite) TestModuleLowered(c *gc.C) {
	logger1 := loggo.GetLogger("TESTING.MODULE")
	logger2 := loggo.GetLogger("Testing")

	c.Assert(logger1.Name(), gc.Equals, "testing.module")
	c.Assert(logger2.Name(), gc.Equals, "testing")
}

func (*loggerSuite) TestLevelsInherited(c *gc.C) {
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

type stack_error struct {
	message string
	stack   []string
}

func (s *stack_error) Error() string {
	return s.message
}

func (s *stack_error) StackTrace() []string {
	return s.stack
}

func checkLastMessage(c *gc.C, writer *loggo.TestWriter, expected string) {
	log := writer.Log()
	writer.Clear()
	obtained := log[len(log)-1].Message
	c.Check(obtained, gc.Equals, expected)
}

func (*loggerSuite) TestLoggingStrings(c *gc.C) {
	writer := &loggo.TestWriter{}
	loggo.ReplaceDefaultWriter(writer)
	logger := loggo.GetLogger("test")
	logger.SetLogLevel(loggo.TRACE)

	logger.Infof("simple")
	checkLastMessage(c, writer, "simple")

	logger.Infof("with args %d", 42)
	checkLastMessage(c, writer, "with args 42")

	logger.Infof("working 100%")
	checkLastMessage(c, writer, "working 100%")

	logger.Infof("missing %s")
	checkLastMessage(c, writer, "missing %s")
}

func (*loggerSuite) TestLocationCapture(c *gc.C) {
	writer := &loggo.TestWriter{}
	loggo.ReplaceDefaultWriter(writer)
	logger := loggo.GetLogger("test")
	logger.SetLogLevel(loggo.TRACE)

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

type logwriterSuite struct {
	logger loggo.Logger
	writer *loggo.TestWriter
}

var _ = gc.Suite(&logwriterSuite{})

func (s *logwriterSuite) SetUpTest(c *gc.C) {
	loggo.ResetLoggers()
	loggo.RemoveWriter("default")
	s.writer = &loggo.TestWriter{}
	err := loggo.RegisterWriter("test", s.writer, loggo.TRACE)
	c.Assert(err, gc.IsNil)
	s.logger = loggo.GetLogger("test.writer")
	// Make it so the logger itself writes all messages.
	s.logger.SetLogLevel(loggo.TRACE)
}

func (s *logwriterSuite) TearDownTest(c *gc.C) {
	loggo.ResetWriters()
}

func (s *logwriterSuite) TestLogDoesntLogWeirdLevels(c *gc.C) {
	s.logger.Logf(loggo.UNSPECIFIED, "message")
	c.Assert(s.writer.Log(), gc.HasLen, 0)

	s.logger.Logf(loggo.Level(42), "message")
	c.Assert(s.writer.Log(), gc.HasLen, 0)

	s.logger.Logf(loggo.CRITICAL+loggo.Level(1), "message")
	c.Assert(s.writer.Log(), gc.HasLen, 0)
}

func (s *logwriterSuite) TestMessageFormatting(c *gc.C) {
	s.logger.Logf(loggo.INFO, "some %s included", "formatting")
	log := s.writer.Log()
	c.Assert(log, gc.HasLen, 1)
	c.Assert(log[0].Message, gc.Equals, "some formatting included")
	c.Assert(log[0].Level, gc.Equals, loggo.INFO)
}

func (s *logwriterSuite) BenchmarkLoggingNoWriters(c *gc.C) {
	// No writers
	loggo.RemoveWriter("test")
	for i := 0; i < c.N; i++ {
		s.logger.Warningf("just a simple warning for %d", i)
	}
}

func (s *logwriterSuite) BenchmarkLoggingNoWritersNoFormat(c *gc.C) {
	// No writers
	loggo.RemoveWriter("test")
	for i := 0; i < c.N; i++ {
		s.logger.Warningf("just a simple warning")
	}
}

func (s *logwriterSuite) BenchmarkLoggingTestWriters(c *gc.C) {
	for i := 0; i < c.N; i++ {
		s.logger.Warningf("just a simple warning for %d", i)
	}
	c.Assert(s.writer.Log, gc.HasLen, c.N)
}

func setupTempFileWriter(c *gc.C) (logFile *os.File, cleanup func()) {
	loggo.RemoveWriter("test")
	logFile, err := ioutil.TempFile("", "loggo-test")
	c.Assert(err, gc.IsNil)
	cleanup = func() {
		logFile.Close()
		os.Remove(logFile.Name())
	}
	writer := loggo.NewSimpleWriter(logFile, &loggo.DefaultFormatter{})
	err = loggo.RegisterWriter("testfile", writer, loggo.TRACE)
	c.Assert(err, gc.IsNil)
	return
}

func (s *logwriterSuite) BenchmarkLoggingDiskWriter(c *gc.C) {
	logFile, cleanup := setupTempFileWriter(c)
	defer cleanup()
	msg := "just a simple warning for %d"
	for i := 0; i < c.N; i++ {
		s.logger.Warningf(msg, i)
	}
	offset, err := logFile.Seek(0, os.SEEK_CUR)
	c.Assert(err, gc.IsNil)
	c.Assert((offset > int64(len(msg))*int64(c.N)), gc.Equals, true,
		gc.Commentf("Not enough data was written to the log file."))
}

func (s *logwriterSuite) BenchmarkLoggingDiskWriterNoMessages(c *gc.C) {
	logFile, cleanup := setupTempFileWriter(c)
	defer cleanup()
	// Change the log level
	writer, _, err := loggo.RemoveWriter("testfile")
	c.Assert(err, gc.IsNil)
	loggo.RegisterWriter("testfile", writer, loggo.WARNING)
	msg := "just a simple warning for %d"
	for i := 0; i < c.N; i++ {
		s.logger.Debugf(msg, i)
	}
	offset, err := logFile.Seek(0, os.SEEK_CUR)
	c.Assert(err, gc.IsNil)
	c.Assert(offset, gc.Equals, int64(0),
		gc.Commentf("Data was written to the log file."))
}

func (s *logwriterSuite) BenchmarkLoggingDiskWriterNoMessagesLogLevel(c *gc.C) {
	logFile, cleanup := setupTempFileWriter(c)
	defer cleanup()
	// Change the log level
	s.logger.SetLogLevel(loggo.WARNING)
	msg := "just a simple warning for %d"
	for i := 0; i < c.N; i++ {
		s.logger.Debugf(msg, i)
	}
	offset, err := logFile.Seek(0, os.SEEK_CUR)
	c.Assert(err, gc.IsNil)
	c.Assert(offset, gc.Equals, int64(0),
		gc.Commentf("Data was written to the log file."))
}
