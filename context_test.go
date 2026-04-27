// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"context"
	"testing"

	"github.com/juju/loggo/v3"
	"github.com/juju/tc"
)

type ContextSuite struct{}

func TestContextSuite(t *testing.T) {
	tc.Run(t, &ContextSuite{})
}

func (*ContextSuite) TestNewContextRootLevel(c *tc.C) {
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
		c.Check(cfg, tc.HasLen, 1)
		value, found := cfg[""]
		c.Check(found, tc.Equals, true)
		c.Check(value, tc.Equals, test.expected)
	}
}

func logAllSeverities(logger loggo.Logger) {
	_ = logger.Criticalf(context.Background(), "something critical")
	_ = logger.Errorf(context.Background(), "an error")
	_ = logger.Warningf(context.Background(), "a warning message")
	_ = logger.Infof(context.Background(), "an info message")
	_ = logger.Debugf(context.Background(), "a debug message")
	_ = logger.Tracef(context.Background(), "a trace message")
}

func checkLogEntry(c *tc.C, entry, expected loggo.Entry) {
	c.Check(entry.Level, tc.Equals, expected.Level)
	c.Check(entry.Module, tc.Equals, expected.Module)
	c.Check(entry.Message, tc.Equals, expected.Message)
}

func checkLogEntries(c *tc.C, obtained, expected []loggo.Entry) {
	if c.Check(len(obtained), tc.Equals, len(expected)) {
		for i := range obtained {
			checkLogEntry(c, obtained[i], expected[i])
		}
	}
}

func (*ContextSuite) TestGetLoggerRoot(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	blank := context.GetLogger("")
	root := context.GetLogger("<root>")
	c.Assert(blank, tc.DeepEquals, root)
}

func (*ContextSuite) TestGetLoggerCase(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	upper := context.GetLogger("TEST")
	lower := context.GetLogger("test")
	c.Assert(upper, tc.DeepEquals, lower)
	c.Assert(upper.Name(), tc.Equals, "test")
}

func (*ContextSuite) TestGetLoggerSpace(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	space := context.GetLogger(" test ")
	lower := context.GetLogger("test")
	c.Assert(space, tc.DeepEquals, lower)
	c.Assert(space.Name(), tc.Equals, "test")
}

func (*ContextSuite) TestNewContextNoWriter(c *tc.C) {
	// Should be no output.
	context := loggo.NewContext(loggo.DEBUG)
	logger := context.GetLogger("test")
	logAllSeverities(logger)
}

func (*ContextSuite) newContextWithTestWriter(c *tc.C, level loggo.Level) (*loggo.Context, *loggo.TestWriter) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(level)
	err := context.AddWriter("test", writer)
	c.Assert(err, tc.IsNil)
	return context, writer
}

func (s *ContextSuite) TestNewContextRootSeverityWarning(c *tc.C) {
	context, writer := s.newContextWithTestWriter(c, loggo.WARNING)
	logger := context.GetLogger("test")
	logAllSeverities(logger)
	checkLogEntries(c, writer.Log(), []loggo.Entry{
		{Level: loggo.CRITICAL, Module: "test", Message: "something critical"},
		{Level: loggo.ERROR, Module: "test", Message: "an error"},
		{Level: loggo.WARNING, Module: "test", Message: "a warning message"},
	})
}

func (s *ContextSuite) TestNewContextRootSeverityTrace(c *tc.C) {
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

func (*ContextSuite) TestNewContextConfig(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	config := context.Config()
	c.Assert(config, tc.DeepEquals, loggo.Config{"": loggo.DEBUG})
}

func (*ContextSuite) TestNewLoggerAddsConfig(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	_ = context.GetLogger("test.module")
	c.Assert(context.Config(), tc.DeepEquals, loggo.Config{
		"": loggo.DEBUG,
	})
	c.Assert(context.CompleteConfig(), tc.DeepEquals, loggo.Config{
		"":            loggo.DEBUG,
		"test":        loggo.UNSPECIFIED,
		"test.module": loggo.UNSPECIFIED,
	})
}

func (*ContextSuite) TestConfigureLoggers(c *tc.C) {
	context := loggo.NewContext(loggo.INFO)
	err := context.ConfigureLoggers("testing.module=debug")
	c.Assert(err, tc.IsNil)
	expected := "<root>=INFO;testing.module=DEBUG"
	c.Assert(context.Config().String(), tc.Equals, expected)
}

func (*ContextSuite) TestApplyNilConfig(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	context.ApplyConfig(nil)
	c.Assert(context.Config(), tc.DeepEquals, loggo.Config{"": loggo.DEBUG})
}

func (*ContextSuite) TestApplyConfigRootUnspecified(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	context.ApplyConfig(loggo.Config{"": loggo.UNSPECIFIED})
	c.Assert(context.Config(), tc.DeepEquals, loggo.Config{"": loggo.WARNING})
}

func (*ContextSuite) TestApplyConfigRootTrace(c *tc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.ApplyConfig(loggo.Config{"": loggo.TRACE})
	c.Assert(context.Config(), tc.DeepEquals, loggo.Config{"": loggo.TRACE})
}

func (*ContextSuite) TestApplyConfigCreatesModules(c *tc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.ApplyConfig(loggo.Config{"first.second": loggo.TRACE})
	c.Assert(context.Config(), tc.DeepEquals,
		loggo.Config{
			"":             loggo.WARNING,
			"first.second": loggo.TRACE,
		})
	c.Assert(context.CompleteConfig(), tc.DeepEquals,
		loggo.Config{
			"":             loggo.WARNING,
			"first":        loggo.UNSPECIFIED,
			"first.second": loggo.TRACE,
		})
}

func (*ContextSuite) TestApplyConfigAdditive(c *tc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.ApplyConfig(loggo.Config{"first.second": loggo.TRACE})
	context.ApplyConfig(loggo.Config{"other.module": loggo.DEBUG})
	c.Assert(context.Config(), tc.DeepEquals,
		loggo.Config{
			"":             loggo.WARNING,
			"first.second": loggo.TRACE,
			"other.module": loggo.DEBUG,
		})
	c.Assert(context.CompleteConfig(), tc.DeepEquals,
		loggo.Config{
			"":             loggo.WARNING,
			"first":        loggo.UNSPECIFIED,
			"first.second": loggo.TRACE,
			"other":        loggo.UNSPECIFIED,
			"other.module": loggo.DEBUG,
		})
}

func (*ContextSuite) TestGetAllLoggerTags(c *tc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.GetLogger("a.b", "one")
	context.GetLogger("c.d", "one")
	context.GetLogger("e", "two")

	labels := context.GetAllLoggerTags()
	c.Assert(labels, tc.DeepEquals, []string{"one", "two"})
}

func (*ContextSuite) TestGetAllLoggerTagsWithApplyConfig(c *tc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.ApplyConfig(loggo.Config{"#one": loggo.TRACE})

	labels := context.GetAllLoggerTags()
	c.Assert(labels, tc.DeepEquals, []string{})
}

func (*ContextSuite) TestApplyConfigTags(c *tc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.GetLogger("a.b", "one")
	context.GetLogger("c.d", "one")
	context.GetLogger("e", "two")

	context.ApplyConfig(loggo.Config{"#one": loggo.TRACE})
	context.ApplyConfig(loggo.Config{"#two": loggo.DEBUG})

	c.Assert(context.Config(), tc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a.b": loggo.TRACE,
			"c.d": loggo.TRACE,
			"e":   loggo.DEBUG,
		})
	c.Assert(context.CompleteConfig(), tc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a":   loggo.UNSPECIFIED,
			"a.b": loggo.TRACE,
			"c":   loggo.UNSPECIFIED,
			"c.d": loggo.TRACE,
			"e":   loggo.DEBUG,
		})
}

func (*ContextSuite) TestApplyConfigLabelsAppliesToNewLoggers(c *tc.C) {
	context := loggo.NewContext(loggo.WARNING)

	context.ApplyConfig(loggo.Config{"#one": loggo.TRACE})
	context.ApplyConfig(loggo.Config{"#two": loggo.DEBUG})

	context.GetLogger("a.b", "one")
	context.GetLogger("c.d", "one")
	context.GetLogger("e", "two")

	c.Assert(context.Config(), tc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a.b": loggo.TRACE,
			"c.d": loggo.TRACE,
			"e":   loggo.DEBUG,
		})
	c.Assert(context.CompleteConfig(), tc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a":   loggo.UNSPECIFIED,
			"a.b": loggo.TRACE,
			"c":   loggo.UNSPECIFIED,
			"c.d": loggo.TRACE,
			"e":   loggo.DEBUG,
		})
}

func (*ContextSuite) TestApplyConfigLabelsAppliesToNewLoggersWithMultipleTags(c *tc.C) {
	context := loggo.NewContext(loggo.WARNING)

	// Invert the order here, to ensure that the config order doesn't matter,
	// but the way the tags are ordered in `GetLogger`.
	context.ApplyConfig(loggo.Config{"#two": loggo.DEBUG})
	context.ApplyConfig(loggo.Config{"#one": loggo.TRACE})

	context.GetLogger("a.b", "one", "two")

	c.Assert(context.Config(), tc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a.b": loggo.TRACE,
		})
	c.Assert(context.CompleteConfig(), tc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a":   loggo.UNSPECIFIED,
			"a.b": loggo.TRACE,
		})
}

func (*ContextSuite) TestApplyConfigLabelsResetLoggerLevels(c *tc.C) {
	context := loggo.NewContext(loggo.WARNING)

	context.ApplyConfig(loggo.Config{"#one": loggo.TRACE})
	context.ApplyConfig(loggo.Config{"#two": loggo.DEBUG})

	context.GetLogger("a.b", "one")
	context.GetLogger("c.d", "one")
	context.GetLogger("e", "two")

	context.ResetLoggerLevels()

	c.Assert(context.Config(), tc.DeepEquals,
		loggo.Config{
			"": loggo.WARNING,
		})
	c.Assert(context.CompleteConfig(), tc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a":   loggo.UNSPECIFIED,
			"a.b": loggo.UNSPECIFIED,
			"c":   loggo.UNSPECIFIED,
			"c.d": loggo.UNSPECIFIED,
			"e":   loggo.UNSPECIFIED,
		})
}

func (*ContextSuite) TestApplyConfigTagsAddative(c *tc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.ApplyConfig(loggo.Config{"#one": loggo.TRACE})
	context.ApplyConfig(loggo.Config{"#two": loggo.DEBUG})
	c.Assert(context.Config(), tc.DeepEquals,
		loggo.Config{
			"": loggo.WARNING,
		})
	c.Assert(context.CompleteConfig(), tc.DeepEquals,
		loggo.Config{
			"": loggo.WARNING,
		})
}

func (*ContextSuite) TestApplyConfigWithMalformedTag(c *tc.C) {
	context := loggo.NewContext(loggo.WARNING)
	context.GetLogger("a.b", "one")

	context.ApplyConfig(loggo.Config{"#ONE.1": loggo.TRACE})

	c.Assert(context.Config(), tc.DeepEquals,
		loggo.Config{
			"": loggo.WARNING,
		})
	c.Assert(context.CompleteConfig(), tc.DeepEquals,
		loggo.Config{
			"":    loggo.WARNING,
			"a":   loggo.UNSPECIFIED,
			"a.b": loggo.UNSPECIFIED,
		})
}

func (*ContextSuite) TestResetLoggerTags(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	context.ApplyConfig(loggo.Config{"first.second": loggo.TRACE})
	context.ResetLoggerLevels()
	c.Assert(context.Config(), tc.DeepEquals,
		loggo.Config{
			"": loggo.WARNING,
		})
	c.Assert(context.CompleteConfig(), tc.DeepEquals,
		loggo.Config{
			"":             loggo.WARNING,
			"first":        loggo.UNSPECIFIED,
			"first.second": loggo.UNSPECIFIED,
		})
}

func (*ContextSuite) TestWriterNamesNone(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	writers := context.WriterNames()
	c.Assert(writers, tc.HasLen, 0)
}

func (*ContextSuite) TestAddWriterNoName(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	err := context.AddWriter("", nil)
	c.Assert(err.Error(), tc.Equals, "name cannot be empty")
}

func (*ContextSuite) TestAddWriterNil(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	err := context.AddWriter("foo", nil)
	c.Assert(err.Error(), tc.Equals, "writer cannot be nil")
}

func (*ContextSuite) TestNamedAddWriter(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	err := context.AddWriter("foo", &writer{name: "foo"})
	c.Assert(err, tc.IsNil)
	err = context.AddWriter("foo", &writer{name: "foo"})
	c.Assert(err.Error(), tc.Equals, `context already has a writer named "foo"`)

	writers := context.WriterNames()
	c.Assert(writers, tc.DeepEquals, []string{"foo"})
}

func (*ContextSuite) TestRemoveWriter(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	w, err := context.RemoveWriter("unknown")
	c.Assert(err.Error(), tc.Equals, `context has no writer named "unknown"`)
	c.Assert(w, tc.IsNil)
}

func (*ContextSuite) TestRemoveWriterFound(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	original := &writer{name: "foo"}
	err := context.AddWriter("foo", original)
	c.Assert(err, tc.IsNil)
	existing, err := context.RemoveWriter("foo")
	c.Assert(err, tc.IsNil)
	c.Assert(existing, tc.Equals, original)

	writers := context.WriterNames()
	c.Assert(writers, tc.HasLen, 0)
}

func (*ContextSuite) TestReplaceWriterNoName(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	existing, err := context.ReplaceWriter("", nil)
	c.Assert(err.Error(), tc.Equals, "name cannot be empty")
	c.Assert(existing, tc.IsNil)
}

func (*ContextSuite) TestReplaceWriterNil(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	existing, err := context.ReplaceWriter("foo", nil)
	c.Assert(err.Error(), tc.Equals, "writer cannot be nil")
	c.Assert(existing, tc.IsNil)
}

func (*ContextSuite) TestReplaceWriterNotFound(c *tc.C) {
	context := loggo.NewContext(loggo.DEBUG)
	existing, err := context.ReplaceWriter("foo", &writer{})
	c.Assert(err.Error(), tc.Equals, `context has no writer named "foo"`)
	c.Assert(existing, tc.IsNil)
}

func (*ContextSuite) TestMultipleWriters(c *tc.C) {
	first := &writer{}
	second := &writer{}
	third := &writer{}
	context := loggo.NewContext(loggo.TRACE)
	err := context.AddWriter("first", first)
	c.Assert(err, tc.IsNil)
	err = context.AddWriter("second", second)
	c.Assert(err, tc.IsNil)
	err = context.AddWriter("third", third)
	c.Assert(err, tc.IsNil)

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

func (*ContextSuite) TestWriter(c *tc.C) {
	first := &writer{name: "first"}
	second := &writer{name: "second"}
	context := loggo.NewContext(loggo.TRACE)
	err := context.AddWriter("first", first)
	c.Assert(err, tc.IsNil)
	err = context.AddWriter("second", second)
	c.Assert(err, tc.IsNil)

	c.Check(context.Writer("unknown"), tc.IsNil)
	c.Check(context.Writer("first"), tc.Equals, first)
	c.Check(context.Writer("second"), tc.Equals, second)

	c.Check(first, tc.Not(tc.Equals), second)
}

type writer struct {
	loggo.TestWriter
	// The name exists to discriminate writer equality.
	name string
}
