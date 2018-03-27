// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"github.com/juju/loggo"
	gc "gopkg.in/check.v1"
)

type ContextSuite struct{}

var _ = gc.Suite(&ContextSuite{})

func (*ContextSuite) TestNewContextRootLevel(c *gc.C) {
	for i, test := range []struct {
		level    loggo.Level
		expected loggo.Level
	}{{
		level:    loggo.UNSPECIFIED,
		expected: loggo.WARNING,
	}, {
		level:    loggo.DEBUG,
		expected: loggo.DEBUG,
	}, {
		level:    loggo.INFO,
		expected: loggo.INFO,
	}, {
		level:    loggo.WARNING,
		expected: loggo.WARNING,
	}, {
		level:    loggo.ERROR,
		expected: loggo.ERROR,
	}, {
		level:    loggo.CRITICAL,
		expected: loggo.CRITICAL,
	}, {
		level:    loggo.Level(42),
		expected: loggo.WARNING,
	}} {
		c.Log("%d: %s", i, test.level)
		context := loggo.NewContext(test.level)
		cfg := context.Config()
		c.Check(cfg, gc.HasLen, 1)
		value, found := cfg[""]
		c.Check(found, gc.Equals, true)
		c.Check(value, gc.Equals, test.expected)
	}
}

func logAllSeverities(logger loggo.Logger) {
	logger.Criticalf("something critical")
	logger.Errorf("an error")
	logger.Warningf("a warning message")
	logger.Infof("an info message")
	logger.Debugf("a debug message")
	logger.Tracef("a trace message")
}

func checkLogEntry(c *gc.C, entry, expected loggo.Entry) {
	c.Check(entry.Level, gc.Equals, expected.Level)
	c.Check(entry.Module, gc.Equals, expected.Module)
	c.Check(entry.Message, gc.Equals, expected.Message)
}

func checkLogEntries(c *gc.C, obtained, expected []loggo.Entry) {
	if c.Check(len(obtained), gc.Equals, len(expected)) {
		for i := range obtained {
			checkLogEntry(c, obtained[i], expected[i])
		}
	}
}

func (*ContextSuite) TestGetLoggerRoot(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	blank := context.GetLogger("")
	root := context.GetLogger("<root>")
	c.Assert(blank, gc.Equals, root)
}

func (*ContextSuite) TestGetLoggerCase(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	upper := context.GetLogger("TEST")
	lower := context.GetLogger("test")
	c.Assert(upper, gc.Equals, lower)
	c.Assert(upper.Name(), gc.Equals, "test")
}

func (*ContextSuite) TestGetLoggerSpace(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	space := context.GetLogger(" test ")
	lower := context.GetLogger("test")
	c.Assert(space, gc.Equals, lower)
	c.Assert(space.Name(), gc.Equals, "test")
}

func (*ContextSuite) TestNewContextNoWriter(c *gc.C) {
	// Should be no output.
	context := loggo.NewContext(loggo.DEBUG)
	logger := context.GetLogger("test")
	logAllSeverities(logger)
}

func (*ContextSuite) newContextWithTestWriter(c *gc.C, level loggo.Level) (*loggo.Context, *loggo.TestWriter) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(level)
	context.AddWriter("test", writer)
	return context, writer
}

func (s *ContextSuite) TestNewContextRootSeverityWarning(c *gc.C) {
	context, writer := s.newContextWithTestWriter(c, loggo.WARNING)
	logger := context.GetLogger("test")
	logAllSeverities(logger)
	checkLogEntries(c, writer.Log(), []loggo.Entry{
		{Level: loggo.CRITICAL, Module: "test", Message: "something critical"},
		{Level: loggo.ERROR, Module: "test", Message: "an error"},
		{Level: loggo.WARNING, Module: "test", Message: "a warning message"},
	})
}

func (s *ContextSuite) TestNewContextRootSeverityTrace(c *gc.C) {
	context, writer := s.newContextWithTestWriter(c, loggo.TRACE)
	logger := context.GetLogger("test")
	logAllSeverities(logger)
	checkLogEntries(c, writer.Log(), []loggo.Entry{
		{Level: loggo.CRITICAL, Module: "test", Message: "something critical"},
		{Level: loggo.ERROR, Module: "test", Message: "an error"},
		{Level: loggo.WARNING, Module: "test", Message: "a warning message"},
		{Level: loggo.INFO, Module: "test", Message: "an info message"},
		{Level: loggo.DEBUG, Module: "test", Message: "a debug message"},
		{Level: loggo.TRACE, Module: "test", Message: "a trace message"},
	})
}

func (*ContextSuite) TestNewContextConfig(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	config := context.Config()
	c.Assert(config, gc.DeepEquals, loggo.Config{"": loggo.DEBUG})
}

func (*ContextSuite) TestNewLoggerAddsConfig(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	_ = context.GetLogger("test.module")
	c.Assert(context.Config(), gc.DeepEquals, loggo.Config{
		"": loggo.DEBUG,
	})
	c.Assert(context.CompleteConfig(), gc.DeepEquals, loggo.Config{
		"":            loggo.DEBUG,
		"test":        loggo.UNSPECIFIED,
		"test.module": loggo.UNSPECIFIED,
	})
}

func (*ContextSuite) TestApplyNilConfig(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	context.ApplyConfig(nil)
	c.Assert(context.Config(), gc.DeepEquals, loggo.Config{"": loggo.DEBUG})
}

func (*ContextSuite) TestApplyConfigRootUnspecified(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	context.ApplyConfig(loggo.Config{"": loggo.UNSPECIFIED})
	c.Assert(context.Config(), gc.DeepEquals, loggo.Config{"": loggo.WARNING})
}

func (*ContextSuite) TestApplyConfigRootTrace(c *gc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.ApplyConfig(loggo.Config{"": loggo.TRACE})
	c.Assert(context.Config(), gc.DeepEquals, loggo.Config{"": loggo.TRACE})
}

func (*ContextSuite) TestApplyConfigCreatesModules(c *gc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.ApplyConfig(loggo.Config{"first.second": loggo.TRACE})
	c.Assert(context.Config(), gc.DeepEquals,
		loggo.Config{
			"":             loggo.WARNING,
			"first.second": loggo.TRACE,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		loggo.Config{
			"":             loggo.WARNING,
			"first":        loggo.UNSPECIFIED,
			"first.second": loggo.TRACE,
		})
}

func (*ContextSuite) TestApplyConfigAdditive(c *gc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.ApplyConfig(loggo.Config{"first.second": loggo.TRACE})
	context.ApplyConfig(loggo.Config{"other.module": loggo.DEBUG})
	c.Assert(context.Config(), gc.DeepEquals,
		loggo.Config{
			"":             loggo.WARNING,
			"first.second": loggo.TRACE,
			"other.module": loggo.DEBUG,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		loggo.Config{
			"":             loggo.WARNING,
			"first":        loggo.UNSPECIFIED,
			"first.second": loggo.TRACE,
			"other":        loggo.UNSPECIFIED,
			"other.module": loggo.DEBUG,
		})
}

func (*ContextSuite) TestResetLoggerLevels(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	context.ApplyConfig(loggo.Config{"first.second": loggo.TRACE})
	context.ResetLoggerLevels()
	c.Assert(context.Config(), gc.DeepEquals,
		loggo.Config{
			"": loggo.WARNING,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		loggo.Config{
			"":             loggo.WARNING,
			"first":        loggo.UNSPECIFIED,
			"first.second": loggo.UNSPECIFIED,
		})
}

func (*ContextSuite) TestWriterNamesNone(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	writers := context.WriterNames()
	c.Assert(writers, gc.HasLen, 0)
}

func (*ContextSuite) TestAddWriterNoName(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	err := context.AddWriter("", nil)
	c.Assert(err.Error(), gc.Equals, "name cannot be empty")
}

func (*ContextSuite) TestAddWriterNil(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	err := context.AddWriter("foo", nil)
	c.Assert(err.Error(), gc.Equals, "writer cannot be nil")
}

func (*ContextSuite) TestNamedAddWriter(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	err := context.AddWriter("foo", &writer{name: "foo"})
	c.Assert(err, gc.IsNil)
	err = context.AddWriter("foo", &writer{name: "foo"})
	c.Assert(err.Error(), gc.Equals, `context already has a writer named "foo"`)

	writers := context.WriterNames()
	c.Assert(writers, gc.DeepEquals, []string{"foo"})
}

func (*ContextSuite) TestRemoveWriter(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	w, err := context.RemoveWriter("unknown")
	c.Assert(err.Error(), gc.Equals, `context has no writer named "unknown"`)
	c.Assert(w, gc.IsNil)
}

func (*ContextSuite) TestRemoveWriterFound(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	original := &writer{name: "foo"}
	err := context.AddWriter("foo", original)
	c.Assert(err, gc.IsNil)
	existing, err := context.RemoveWriter("foo")
	c.Assert(err, gc.IsNil)
	c.Assert(existing, gc.Equals, original)

	writers := context.WriterNames()
	c.Assert(writers, gc.HasLen, 0)
}

func (*ContextSuite) TestReplaceWriterNoName(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	existing, err := context.ReplaceWriter("", nil)
	c.Assert(err.Error(), gc.Equals, "name cannot be empty")
	c.Assert(existing, gc.IsNil)
}

func (*ContextSuite) TestReplaceWriterNil(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	existing, err := context.ReplaceWriter("foo", nil)
	c.Assert(err.Error(), gc.Equals, "writer cannot be nil")
	c.Assert(existing, gc.IsNil)
}

func (*ContextSuite) TestReplaceWriterNotFound(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	existing, err := context.ReplaceWriter("foo", &writer{})
	c.Assert(err.Error(), gc.Equals, `context has no writer named "foo"`)
	c.Assert(existing, gc.IsNil)
}

func (*ContextSuite) TestMultipleWriters(c *gc.C) {
	first := &writer{}
	second := &writer{}
	third := &writer{}
	context := loggo.NewContext(loggo.TRACE)
	err := context.AddWriter("first", first)
	c.Assert(err, gc.IsNil)
	err = context.AddWriter("second", second)
	c.Assert(err, gc.IsNil)
	err = context.AddWriter("third", third)
	c.Assert(err, gc.IsNil)

	logger := context.GetLogger("test")
	logAllSeverities(logger)

	expected := []loggo.Entry{
		{Level: loggo.CRITICAL, Module: "test", Message: "something critical"},
		{Level: loggo.ERROR, Module: "test", Message: "an error"},
		{Level: loggo.WARNING, Module: "test", Message: "a warning message"},
		{Level: loggo.INFO, Module: "test", Message: "an info message"},
		{Level: loggo.DEBUG, Module: "test", Message: "a debug message"},
		{Level: loggo.TRACE, Module: "test", Message: "a trace message"},
	}

	checkLogEntries(c, first.Log(), expected)
	checkLogEntries(c, second.Log(), expected)
	checkLogEntries(c, third.Log(), expected)
}

func (*ContextSuite) TestWriter(c *gc.C) {
	first := &writer{name: "first"}
	second := &writer{name: "second"}
	context := loggo.NewContext(loggo.TRACE)
	err := context.AddWriter("first", first)
	c.Assert(err, gc.IsNil)
	err = context.AddWriter("second", second)
	c.Assert(err, gc.IsNil)

	c.Check(context.Writer("unknown"), gc.IsNil)
	c.Check(context.Writer("first"), gc.Equals, first)
	c.Check(context.Writer("second"), gc.Equals, second)

	c.Check(first, gc.Not(gc.Equals), second)

}

type writer struct {
	loggo.TestWriter
	// The name exists to discriminate writer equality.
	name string
}
