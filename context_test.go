// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"github.com/juju/loggo/v2"

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
		c.Logf("%d: %s", i, test.level)
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
	c.Assert(blank, gc.DeepEquals, root)
}

func (*ContextSuite) TestGetLoggerCase(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	upper := context.GetLogger("TEST")
	lower := context.GetLogger("test")
	c.Assert(upper, gc.DeepEquals, lower)
	c.Assert(upper.Name(), gc.Equals, "test")
}

func (*ContextSuite) TestGetLoggerSpace(c *gc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	space := context.GetLogger(" test ")
	lower := context.GetLogger("test")
	c.Assert(space, gc.DeepEquals, lower)
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
	err := context.AddWriter("test", writer)
	c.Assert(err, gc.IsNil)
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

func (*ContextSuite) TestConfigureLoggers(c *gc.C) {
	context := loggo.NewContext(loggo.INFO)
	err := context.ConfigureLoggers("testing.module=debug")
	c.Assert(err, gc.IsNil)
	expected := "<root>=INFO;testing.module=DEBUG"
	c.Assert(context.Config().String(), gc.Equals, expected)
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

func (*ContextSuite) TestGetAllLoggerTags(c *gc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.GetLogger("a.b", "one")
	context.GetLogger("c.d", "one")
	context.GetLogger("e", "two")

	labels := context.GetAllLoggerTags()
	c.Assert(labels, gc.DeepEquals, []string{"one", "two"})
}

func (*ContextSuite) TestGetAllLoggerTagsWithApplyConfig(c *gc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.ApplyConfig(loggo.Config{"#one": loggo.TRACE})

	labels := context.GetAllLoggerTags()
	c.Assert(labels, gc.DeepEquals, []string{})
}

func (*ContextSuite) TestApplyConfigTags(c *gc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.GetLogger("a.b", "one")
	context.GetLogger("c.d", "one")
	context.GetLogger("e", "two")

	context.ApplyConfig(loggo.Config{"#one": loggo.TRACE})
	context.ApplyConfig(loggo.Config{"#two": loggo.DEBUG})

	c.Assert(context.Config(), gc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a.b": loggo.TRACE,
			"c.d": loggo.TRACE,
			"e":   loggo.DEBUG,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a":   loggo.UNSPECIFIED,
			"a.b": loggo.TRACE,
			"c":   loggo.UNSPECIFIED,
			"c.d": loggo.TRACE,
			"e":   loggo.DEBUG,
		})
}

func (*ContextSuite) TestApplyConfigTagsAppliesToNewLoggers(c *gc.C) {
	context := loggo.NewContext(loggo.WARNING)

	context.ApplyConfig(loggo.Config{"#one": loggo.TRACE})
	context.ApplyConfig(loggo.Config{"#two": loggo.DEBUG})

	context.GetLogger("a.b", "one")
	context.GetLogger("c.d", "one")
	context.GetLogger("e", "two")

	c.Assert(context.Config(), gc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a.b": loggo.TRACE,
			"c.d": loggo.TRACE,
			"e":   loggo.DEBUG,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a":   loggo.UNSPECIFIED,
			"a.b": loggo.TRACE,
			"c":   loggo.UNSPECIFIED,
			"c.d": loggo.TRACE,
			"e":   loggo.DEBUG,
		})
}

func (*ContextSuite) TestApplyConfigTagsAppliesToNewLoggersWithMultipleTags(c *gc.C) {
	context := loggo.NewContext(loggo.WARNING)

	// Invert the order here, to ensure that the config order doesn't matter,
	// but the way the tags are ordered in `GetLogger`.
	context.ApplyConfig(loggo.Config{"#two": loggo.DEBUG})
	context.ApplyConfig(loggo.Config{"#one": loggo.TRACE})

	context.GetLogger("a.b", "one", "two")

	c.Assert(context.Config(), gc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a.b": loggo.TRACE,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a":   loggo.UNSPECIFIED,
			"a.b": loggo.TRACE,
		})
}

func (*ContextSuite) TestApplyConfigTagsResetLoggerLevels(c *gc.C) {
	context := loggo.NewContext(loggo.WARNING)

	context.ApplyConfig(loggo.Config{"#one": loggo.TRACE})
	context.ApplyConfig(loggo.Config{"#two": loggo.DEBUG})

	context.GetLogger("a.b", "one")
	context.GetLogger("c.d", "one")
	context.GetLogger("e", "two")

	context.ResetLoggerLevels()

	c.Assert(context.Config(), gc.DeepEquals,
		loggo.Config{
			"": loggo.WARNING,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a":   loggo.UNSPECIFIED,
			"a.b": loggo.UNSPECIFIED,
			"c":   loggo.UNSPECIFIED,
			"c.d": loggo.UNSPECIFIED,
			"e":   loggo.UNSPECIFIED,
		})
}

func (*ContextSuite) TestApplyConfigTagsResetLoggerLevelsUsingLabels(c *gc.C) {
	context := loggo.NewContext(loggo.WARNING)

	context.ApplyConfig(loggo.Config{"#one": loggo.TRACE})
	context.ApplyConfig(loggo.Config{"#two": loggo.DEBUG})

	context.GetLogger("a", "one").ChildWithLabels("b", loggo.Labels{"x": "y"})
	context.GetLogger("c.d", "one")
	context.GetLogger("e", "two")

	// If a label is available on a logger, then resetting the levels should
	// not remove the label.

	context.ResetLoggerLevels()

	c.Assert(context.Config(), gc.DeepEquals,
		loggo.Config{
			"": loggo.WARNING,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a":   loggo.UNSPECIFIED,
			"a.b": loggo.UNSPECIFIED,
			"c":   loggo.UNSPECIFIED,
			"c.d": loggo.UNSPECIFIED,
			"e":   loggo.UNSPECIFIED,
		})
}

func (*ContextSuite) TestApplyConfigTagsResetLoggerLevelsUsingLabelsRemoval(c *gc.C) {
	context := loggo.NewContext(loggo.WARNING)

	context.ApplyConfig(loggo.Config{"#one": loggo.TRACE})
	context.ApplyConfig(loggo.Config{"#two": loggo.DEBUG})

	context.GetLogger("a", "one").ChildWithLabels("b", loggo.Labels{"x": "y"}).ChildWithTags("g", "one")
	context.GetLogger("c.d", "one")
	context.GetLogger("e", "two")
	context.GetLogger("f")

	// Ensure that the logger that matches exactly the label is removed,
	// including it's children. So we observe hierarchy in the removal.

	c.Assert(context.Config(), gc.DeepEquals,
		loggo.Config{
			"":      loggo.WARNING,
			"a":     loggo.TRACE,
			"a.b.g": loggo.TRACE,
			"c.d":   loggo.TRACE,
			"e":     loggo.DEBUG,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		loggo.Config{
			"":      loggo.WARNING,
			"a":     loggo.TRACE,
			"a.b":   loggo.UNSPECIFIED,
			"a.b.g": loggo.TRACE,
			"c":     loggo.UNSPECIFIED,
			"c.d":   loggo.TRACE,
			"e":     loggo.DEBUG,
			"f":     loggo.UNSPECIFIED,
		})

	context.ResetLoggerLevels(loggo.Labels{"x": "y"})

	c.Assert(context.Config(), gc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a":   loggo.TRACE,
			"c.d": loggo.TRACE,
			"e":   loggo.DEBUG,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		loggo.Config{
			"":      loggo.WARNING,
			"a":     loggo.TRACE,
			"a.b":   loggo.UNSPECIFIED,
			"a.b.g": loggo.UNSPECIFIED,
			"c":     loggo.UNSPECIFIED,
			"c.d":   loggo.TRACE,
			"e":     loggo.DEBUG,
			"f":     loggo.UNSPECIFIED,
		})
}

func (*ContextSuite) TestApplyConfigTagsAdditive(c *gc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.ApplyConfig(loggo.Config{"#one": loggo.TRACE})
	context.ApplyConfig(loggo.Config{"#two": loggo.DEBUG})
	c.Assert(context.Config(), gc.DeepEquals,
		loggo.Config{
			"": loggo.WARNING,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		loggo.Config{
			"": loggo.WARNING,
		})
}

func (*ContextSuite) TestApplyConfigWithMalformedTag(c *gc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.GetLogger("a.b", "one")

	context.ApplyConfig(loggo.Config{"#ONE.1": loggo.TRACE})

	c.Assert(context.Config(), gc.DeepEquals,
		loggo.Config{
			"": loggo.WARNING,
		})
	c.Assert(context.CompleteConfig(), gc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a":   loggo.UNSPECIFIED,
			"a.b": loggo.UNSPECIFIED,
		})
}

func (*ContextSuite) TestResetLoggerTags(c *gc.C) {
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
