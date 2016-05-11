// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"io/ioutil"
	"os"

	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggotest"
)

type GlobalLoggersSuite struct{}

var _ = gc.Suite(&GlobalLoggersSuite{})

func (*GlobalLoggersSuite) SetUpTest(c *gc.C) {
	loggo.ResetLoggers()
}

func (*GlobalLoggersSuite) TestRootLogger(c *gc.C) {
	var root loggo.Logger

	got := loggo.GetLogger("")

	c.Check(got.Name(), gc.Equals, root.Name())
	c.Check(got.LogLevel(), gc.Equals, root.LogLevel())
}

func (*GlobalLoggersSuite) TestModuleName(c *gc.C) {
	logger := loggo.GetLogger("loggo.testing")

	c.Check(logger.Name(), gc.Equals, "loggo.testing")
}

func (*GlobalLoggersSuite) TestLevel(c *gc.C) {
	logger := loggo.GetLogger("testing")

	level := logger.LogLevel()

	c.Check(level, gc.Equals, loggo.UNSPECIFIED)
}

func (*GlobalLoggersSuite) TestEffectiveLevel(c *gc.C) {
	logger := loggo.GetLogger("testing")

	level := logger.EffectiveLogLevel()

	c.Check(level, gc.Equals, loggo.WARNING)
}

func (*GlobalLoggersSuite) TestLevelsSharedForSameModule(c *gc.C) {
	logger1 := loggo.GetLogger("testing.module")
	logger2 := loggo.GetLogger("testing.module")

	logger1.SetLogLevel(loggo.INFO)
	c.Assert(logger1.IsInfoEnabled(), gc.Equals, true)
	c.Assert(logger2.IsInfoEnabled(), gc.Equals, true)
}

func (*GlobalLoggersSuite) TestModuleLowered(c *gc.C) {
	logger1 := loggo.GetLogger("TESTING.MODULE")
	logger2 := loggo.GetLogger("Testing")

	c.Assert(logger1.Name(), gc.Equals, "testing.module")
	c.Assert(logger2.Name(), gc.Equals, "testing")
}

func (*GlobalLoggersSuite) TestLevelsInherited(c *gc.C) {
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

var configureLoggersTests = []struct {
	spec string
	info string
	err  string
}{{
	spec: "",
	info: "<root>=WARNING",
}, {
	spec: "<root>=UNSPECIFIED",
	info: "<root>=WARNING",
}, {
	spec: "<root>=DEBUG",
	info: "<root>=DEBUG",
}, {
	spec: "TRACE",
	info: "<root>=TRACE",
}, {
	spec: "test.module=debug",
	info: "<root>=WARNING;test.module=DEBUG",
}, {
	spec: "module=info; sub.module=debug; other.module=warning",
	info: "<root>=WARNING;module=INFO;other.module=WARNING;sub.module=DEBUG",
}, {
	spec: "  foo.bar \t\r\n= \t\r\nCRITICAL \t\r\n; \t\r\nfoo \r\t\n = DEBUG",
	info: "<root>=WARNING;foo=DEBUG;foo.bar=CRITICAL",
}, {
	spec: "foo;bar",
	info: "<root>=WARNING",
	err:  `logger entry expected '=', found "foo"`,
}, {
	spec: "=foo",
	info: "<root>=WARNING",
	err:  `logger entry "=foo" has blank name`,
}, {
	spec: "foo=",
	info: "<root>=WARNING",
	err:  `logger entry "foo=" has blank config`,
}, {
	spec: "=",
	info: "<root>=WARNING",
	err:  `logger entry "=" has blank name`,
}, {
	spec: "foo=unknown",
	info: "<root>=WARNING",
	err:  `unknown log level "unknown"`,
}, {
	// Test that nothing is changed even when the
	// first part of the specification parses ok.
	spec: "module=info; foo=unknown",
	info: "<root>=WARNING",
	err:  `unknown log level "unknown"`,
}}

func (s *GlobalLoggersSuite) TestConfigureLoggers(c *gc.C) {
	for i, test := range configureLoggersTests {
		c.Logf("test %d: %q", i, test.spec)
		loggo.ResetLoggers()
		err := loggo.ConfigureLoggers(test.spec)
		c.Check(loggo.LoggerInfo(), gc.Equals, test.info)
		if test.err != "" {
			c.Assert(err, gc.ErrorMatches, test.err)
			continue
		}
		c.Assert(err, gc.IsNil)

		// Test that it's idempotent.
		err = loggo.ConfigureLoggers(test.spec)
		c.Assert(err, gc.IsNil)
		c.Assert(loggo.LoggerInfo(), gc.Equals, test.info)

		// Test that calling ConfigureLoggers with the
		// output of LoggerInfo works too.
		err = loggo.ConfigureLoggers(test.info)
		c.Assert(err, gc.IsNil)
		c.Assert(loggo.LoggerInfo(), gc.Equals, test.info)
	}
}

type GlobalWritersSuite struct {
	logger loggo.Logger
	writer *loggotest.Writer
}

var _ = gc.Suite(&GlobalWritersSuite{})

func (s *GlobalWritersSuite) SetUpTest(c *gc.C) {
	loggo.ResetLoggers()
}

func (s *GlobalWritersSuite) TearDownTest(c *gc.C) {
	loggo.ResetWriters()
}

func (s *GlobalWritersSuite) TestLoggerUsesDefault(c *gc.C) {
	writer := &loggotest.Writer{}
	_, err := loggo.ReplaceDefaultWriter(writer)
	c.Assert(err, gc.IsNil)
	logger := loggo.GetLogger("test.writer")
	logger.SetLogLevel(loggo.TRACE)

	logger.Infof("message")

	loggotest.CheckLastMessage(c, writer, "message")
}

func (s *GlobalWritersSuite) TestLoggerWriterRegisteredLater(c *gc.C) {
	writer := &loggotest.Writer{}
	logger := loggo.GetLogger("test.writer")
	logger.SetLogLevel(loggo.TRACE)
	_, err := loggo.ReplaceDefaultWriter(writer)
	c.Assert(err, gc.IsNil)

	logger.Infof("message")

	loggotest.CheckLastMessage(c, writer, "message")
}

func (s *GlobalWritersSuite) TestLoggerUsesMultiple(c *gc.C) {
	writer := &loggotest.Writer{}
	_, err := loggo.ReplaceDefaultWriter(writer)
	c.Assert(err, gc.IsNil)
	err = loggo.RegisterWriter("test", writer, loggo.TRACE)
	c.Assert(err, gc.IsNil)
	logger := loggo.GetLogger("test.writer")
	logger.SetLogLevel(loggo.TRACE)

	logger.Infof("message")

	log := writer.Log()
	c.Assert(log, gc.HasLen, 2)
	c.Check(log[0].Message, gc.Equals, "message")
	c.Check(log[1].Message, gc.Equals, "message")
}

func (s *GlobalWritersSuite) TestLoggerRespectsWriterLevel(c *gc.C) {
	logger := loggo.GetLogger("test.writer")
	loggo.RemoveWriter("default")
	writer := &loggotest.Writer{}
	err := loggo.RegisterWriter("test", writer, loggo.ERROR)
	c.Assert(err, gc.IsNil)

	logger.Infof("message")

	c.Check(writer.Log(), gc.HasLen, 0)
}

func (*GlobalWritersSuite) TestRemoveDefaultWriter(c *gc.C) {
	defaultWriter, level, err := loggo.RemoveWriter("default")
	c.Assert(err, gc.IsNil)
	c.Assert(level, gc.Equals, loggo.TRACE)
	c.Assert(defaultWriter, gc.NotNil)

	// Trying again fails.
	defaultWriter, level, err = loggo.RemoveWriter("default")
	c.Assert(err, gc.ErrorMatches, `Writer "default" is not recognized`)
	c.Assert(level, gc.Equals, loggo.UNSPECIFIED)
	c.Assert(defaultWriter, gc.IsNil)
}

func (*GlobalWritersSuite) TestRegisterWriterExistingName(c *gc.C) {
	err := loggo.RegisterWriter("default", &loggotest.Writer{}, loggo.INFO)
	c.Assert(err, gc.ErrorMatches, `there is already a Writer with the name "default"`)
}

func (*GlobalWritersSuite) TestRegisterNilWriter(c *gc.C) {
	err := loggo.RegisterWriter("nil", nil, loggo.INFO)
	c.Assert(err, gc.ErrorMatches, `Writer cannot be nil`)
}

func (*GlobalWritersSuite) TestRegisterWriterTypedNil(c *gc.C) {
	// If the interface is a typed nil, we have to trust the user.
	var writer *loggotest.Writer
	err := loggo.RegisterWriter("nil", writer, loggo.INFO)
	c.Assert(err, gc.IsNil)
}

func (*GlobalWritersSuite) TestReplaceDefaultWriter(c *gc.C) {
	oldWriter, err := loggo.ReplaceDefaultWriter(&loggotest.Writer{})
	c.Assert(oldWriter, gc.NotNil)
	c.Assert(err, gc.IsNil)
}

func (*GlobalWritersSuite) TestReplaceDefaultWriterWithNil(c *gc.C) {
	oldWriter, err := loggo.ReplaceDefaultWriter(nil)
	c.Assert(oldWriter, gc.IsNil)
	c.Assert(err, gc.ErrorMatches, "Writer cannot be nil")
}

func (*GlobalWritersSuite) TestReplaceDefaultWriterNoDefault(c *gc.C) {
	loggo.RemoveWriter("default")
	oldWriter, err := loggo.ReplaceDefaultWriter(&loggotest.Writer{})
	c.Assert(oldWriter, gc.IsNil)
	c.Assert(err, gc.ErrorMatches, `Writer "default" is not recognized`)
}

func (s *GlobalWritersSuite) TestWillWrite(c *gc.C) {
	// By default, the root logger watches TRACE messages
	c.Assert(loggo.WillWrite(loggo.TRACE), gc.Equals, true)
	// Note: ReplaceDefaultWriter doesn't let us change the default log
	//	 level :(
	writer, _, err := loggo.RemoveWriter("default")
	c.Assert(err, gc.IsNil)
	c.Assert(writer, gc.NotNil)
	err = loggo.RegisterWriter("default", writer, loggo.CRITICAL)
	c.Assert(err, gc.IsNil)
	c.Assert(loggo.WillWrite(loggo.TRACE), gc.Equals, false)
	c.Assert(loggo.WillWrite(loggo.DEBUG), gc.Equals, false)
	c.Assert(loggo.WillWrite(loggo.INFO), gc.Equals, false)
	c.Assert(loggo.WillWrite(loggo.WARNING), gc.Equals, false)
	c.Assert(loggo.WillWrite(loggo.CRITICAL), gc.Equals, true)
}

type GlobalBenchmarksSuite struct {
	logger loggo.Logger
	writer *loggotest.Writer
}

var _ = gc.Suite(&GlobalBenchmarksSuite{})

func (s *GlobalBenchmarksSuite) SetUpTest(c *gc.C) {
	loggo.ResetLoggers()
	s.logger = loggo.GetLogger("test.writer")
	s.writer = &loggotest.Writer{}
	loggo.RemoveWriter("default")
	err := loggo.RegisterWriter("test", s.writer, loggo.TRACE)
	c.Assert(err, gc.IsNil)
}

func (s *GlobalBenchmarksSuite) TearDownTest(c *gc.C) {
	loggo.ResetWriters()
}

func (s *GlobalBenchmarksSuite) BenchmarkLoggingNoWriters(c *gc.C) {
	// No writers
	loggo.RemoveWriter("test")
	for i := 0; i < c.N; i++ {
		s.logger.Warningf("just a simple warning for %d", i)
	}
}

func (s *GlobalBenchmarksSuite) BenchmarkLoggingNoWritersNoFormat(c *gc.C) {
	// No writers
	loggo.RemoveWriter("test")
	for i := 0; i < c.N; i++ {
		s.logger.Warningf("just a simple warning")
	}
}

func (s *GlobalBenchmarksSuite) BenchmarkLoggingTestWriters(c *gc.C) {
	for i := 0; i < c.N; i++ {
		s.logger.Warningf("just a simple warning for %d", i)
	}
	c.Assert(s.writer.Log, gc.HasLen, c.N)
}

func (s *GlobalBenchmarksSuite) BenchmarkLoggingDiskWriter(c *gc.C) {
	loggo.RemoveWriter("test")
	logFile := s.setupTempFileWriter(c)
	defer logFile.Close()
	msg := "just a simple warning for %d"
	for i := 0; i < c.N; i++ {
		s.logger.Warningf(msg, i)
	}
	offset, err := logFile.Seek(0, os.SEEK_CUR)
	c.Assert(err, gc.IsNil)
	c.Assert((offset > int64(len(msg))*int64(c.N)), gc.Equals, true,
		gc.Commentf("Not enough data was written to the log file."))
}

func (s *GlobalBenchmarksSuite) BenchmarkLoggingDiskWriterNoMessages(c *gc.C) {
	logFile := s.setupTempFileWriter(c)
	defer logFile.Close()
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

func (s *GlobalBenchmarksSuite) BenchmarkLoggingDiskWriterNoMessagesLogLevel(c *gc.C) {
	logFile := s.setupTempFileWriter(c)
	defer logFile.Close()
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

func (s *GlobalBenchmarksSuite) setupTempFileWriter(c *gc.C) *os.File {
	loggo.RemoveWriter("test")
	logFile, err := ioutil.TempFile(c.MkDir(), "loggo-test")
	c.Assert(err, gc.IsNil)
	writer := loggo.NewSimpleWriter(logFile, &loggo.DefaultFormatter{})
	err = loggo.RegisterWriter("testfile", writer, loggo.TRACE)
	c.Assert(err, gc.IsNil)
	return logFile
}
