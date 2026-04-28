// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"testing"
	"time"

	"github.com/juju/loggo/v3"
	"github.com/juju/loggo/v3/attrs"
	"github.com/juju/tc"
)

type LoggerSuite struct{}

func TestLoggerSuite(t *testing.T) {
	tc.Run(t, &LoggerSuite{})
}

func (*LoggerSuite) SetUpTest(c *tc.C) {
	loggo.ResetDefaultContext()
}

func (s *LoggerSuite) TestRootLogger(c *tc.C) {
	root := loggo.Logger{}.WithCallDepth(2)
	c.Check(root.Name(), tc.Equals, "<root>")
	c.Check(root.LogLevel(), tc.Equals, loggo.WARNING)
	c.Check(root.IsErrorEnabled(), tc.Equals, true)
	c.Check(root.IsWarningEnabled(), tc.Equals, true)
	c.Check(root.IsInfoEnabled(), tc.Equals, false)
	c.Check(root.IsDebugEnabled(), tc.Equals, false)
	c.Check(root.IsTraceEnabled(), tc.Equals, false)
}

func (s *LoggerSuite) TestWithLabels(c *tc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, tc.IsNil)

	logger := context.GetLogger("testing")
	loggerWithLabels := logger.WithLabels(loggo.Labels{"foo": "bar"})
	loggerWithTagsAndLabels := logger.
		ChildWithTags("withTags", "tag1", "tag2").
		WithLabels(loggo.Labels{"hello": "world"})

	_ = logger.Logf(c.Context(), loggo.INFO, "without labels")
	_ = loggerWithLabels.Logf(c.Context(), loggo.INFO, "with labels")
	_ = loggerWithTagsAndLabels.Logf(c.Context(), loggo.INFO, "with tags and labels")

	logs := writer.Log()
	c.Assert(logs, tc.HasLen, 3)
	c.Check(logs[0].Message, tc.Equals, "without labels")
	c.Check(logs[0].Labels, tc.HasLen, 0)
	c.Check(logs[1].Message, tc.Equals, "with labels")
	c.Check(logs[1].Labels, tc.DeepEquals, loggo.Labels{"foo": "bar"})
	c.Check(logs[2].Message, tc.Equals, "with tags and labels")
	c.Check(logs[2].Labels, tc.DeepEquals, loggo.Labels{
		"logger-tags": "tag1,tag2",
		"hello":       "world",
	})
}

func (s *LoggerSuite) TestNonInheritedLabels(c *tc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, tc.IsNil)

	logger := context.GetLogger("testing").
		WithLabels(loggo.Labels{"hello": "world"})

	inheritedLoggerWithLabels := logger.
		ChildWithLabels("inherited", loggo.Labels{"foo": "bar"})

	_ = logger.Logf(c.Context(), loggo.INFO, "with labels")
	_ = inheritedLoggerWithLabels.Logf(c.Context(), loggo.INFO, "with inherited labels")

	logs := writer.Log()
	c.Assert(logs, tc.HasLen, 2)

	// The second log message should _only_ have the inherited labels.

	c.Check(logs[0].Message, tc.Equals, "with labels")
	c.Check(logs[0].Labels, tc.DeepEquals, loggo.Labels{"hello": "world"})

	c.Check(logs[1].Message, tc.Equals, "with inherited labels")
	c.Check(logs[1].Labels, tc.DeepEquals, loggo.Labels{"foo": "bar"})
}

func (s *LoggerSuite) TestNonInheritedWithInheritedLabels(c *tc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, tc.IsNil)

	logger := context.GetLogger("testing")

	inheritedLoggerWithLabels := logger.
		ChildWithLabels("inherited", loggo.Labels{"foo": "bar"})

	scopedLoggerWithLabels := inheritedLoggerWithLabels.
		WithLabels(loggo.Labels{"hello": "world"})

	_ = inheritedLoggerWithLabels.Logf(c.Context(), loggo.INFO, "with inherited labels")
	_ = scopedLoggerWithLabels.Logf(c.Context(), loggo.INFO, "with scoped labels")

	logs := writer.Log()
	c.Assert(logs, tc.HasLen, 2)

	// The second log message should have both the inherited labels and
	// scoped labels.

	c.Check(logs[0].Message, tc.Equals, "with inherited labels")
	c.Check(logs[0].Labels, tc.DeepEquals, loggo.Labels{"foo": "bar"})

	c.Check(logs[1].Message, tc.Equals, "with scoped labels")
	c.Check(logs[1].Labels, tc.DeepEquals, loggo.Labels{
		"foo":   "bar",
		"hello": "world",
	})
}

func (s *LoggerSuite) TestInheritedLabels(c *tc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, tc.IsNil)

	logger := context.GetLogger("testing")

	nestedLoggerWithLabels := logger.
		ChildWithLabels("nested", loggo.Labels{"foo": "bar"})
	deepNestedLoggerWithLabels := nestedLoggerWithLabels.
		ChildWithLabels("nested", loggo.Labels{"foo": "bar"}).
		ChildWithLabels("deepnested", loggo.Labels{"fred": "tim"})

	loggerWithTagsAndLabels := logger.
		ChildWithLabels("nested-labels", loggo.Labels{"hello": "world"}).
		ChildWithTags("nested-tag", "tag1", "tag2")

	_ = logger.Logf(c.Context(), loggo.INFO, "without labels")
	_ = nestedLoggerWithLabels.Logf(c.Context(), loggo.INFO, "with nested labels")
	_ = deepNestedLoggerWithLabels.Logf(c.Context(), loggo.INFO, "with deep nested labels")
	_ = loggerWithTagsAndLabels.Logf(c.Context(), loggo.INFO, "with tags and labels")

	logs := writer.Log()
	c.Assert(logs, tc.HasLen, 4)
	c.Check(logs[0].Message, tc.Equals, "without labels")
	c.Check(logs[0].Labels, tc.HasLen, 0)

	c.Check(logs[1].Message, tc.Equals, "with nested labels")
	c.Check(logs[1].Labels, tc.DeepEquals, loggo.Labels{"foo": "bar"})

	c.Check(logs[2].Message, tc.Equals, "with deep nested labels")
	c.Check(logs[2].Labels, tc.DeepEquals, loggo.Labels{
		"foo":  "bar",
		"fred": "tim",
	})

	c.Check(logs[3].Message, tc.Equals, "with tags and labels")
	c.Check(logs[3].Labels, tc.DeepEquals, loggo.Labels{
		"logger-tags": "tag1,tag2",
		"hello":       "world",
	})
}

func (s *LoggerSuite) TestLogWithStaticAndDynamicLabels(c *tc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, tc.IsNil)

	logger := context.GetLogger("testing")
	loggerWithLabels := logger.WithLabels(loggo.Labels{"foo": "bar"})

	_ = loggerWithLabels.LogWithLabelsf(c.Context(), loggo.INFO, "no extra labels", nil)
	_ = loggerWithLabels.LogWithLabelsf(c.Context(), loggo.INFO, "with extra labels", map[string]string{
		"domain": "status",
		"kind":   "machine",
		"id":     "0",
		"value":  "idle",
	})

	logs := writer.Log()
	c.Assert(logs, tc.HasLen, 2)
	c.Check(logs[0].Message, tc.Equals, "no extra labels")
	c.Check(logs[0].Labels, tc.DeepEquals, loggo.Labels{"foo": "bar"})
	c.Check(logs[1].Message, tc.Equals, "with extra labels")
	c.Check(logs[1].Labels, tc.DeepEquals, loggo.Labels{
		"foo": "bar", "domain": "status", "id": "0", "kind": "machine", "value": "idle"})
}

func (s *LoggerSuite) TestLogWithExtraLabels(c *tc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, tc.IsNil)

	logger := context.GetLogger("testing")

	_ = logger.LogWithLabelsf(c.Context(), loggo.INFO, "no extra labels", nil)
	_ = logger.LogWithLabelsf(c.Context(), loggo.INFO, "with extra labels", map[string]string{
		"domain": "status",
		"kind":   "machine",
		"id":     "0",
		"value":  "idle",
	})

	logs := writer.Log()
	c.Assert(logs, tc.HasLen, 2)
	c.Check(logs[0].Message, tc.Equals, "no extra labels")
	c.Check(logs[0].Labels, tc.HasLen, 0)
	c.Check(logs[1].Message, tc.Equals, "with extra labels")
	c.Check(logs[1].Labels, tc.DeepEquals, loggo.Labels{
		"domain": "status", "id": "0", "kind": "machine", "value": "idle"})
}

func (s *LoggerSuite) TestSetLevel(c *tc.C) {
	logger := loggo.GetLogger("testing")

	c.Assert(logger.LogLevel(), tc.Equals, loggo.UNSPECIFIED)
	c.Assert(logger.EffectiveLogLevel(), tc.Equals, loggo.WARNING)
	c.Assert(logger.IsErrorEnabled(), tc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), tc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), tc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), tc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), tc.Equals, false)
	logger.SetLogLevel(loggo.TRACE)
	c.Assert(logger.LogLevel(), tc.Equals, loggo.TRACE)
	c.Assert(logger.EffectiveLogLevel(), tc.Equals, loggo.TRACE)
	c.Assert(logger.IsErrorEnabled(), tc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), tc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), tc.Equals, true)
	c.Assert(logger.IsDebugEnabled(), tc.Equals, true)
	c.Assert(logger.IsTraceEnabled(), tc.Equals, true)
	logger.SetLogLevel(loggo.DEBUG)
	c.Assert(logger.LogLevel(), tc.Equals, loggo.DEBUG)
	c.Assert(logger.EffectiveLogLevel(), tc.Equals, loggo.DEBUG)
	c.Assert(logger.IsErrorEnabled(), tc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), tc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), tc.Equals, true)
	c.Assert(logger.IsDebugEnabled(), tc.Equals, true)
	c.Assert(logger.IsTraceEnabled(), tc.Equals, false)
	logger.SetLogLevel(loggo.INFO)
	c.Assert(logger.LogLevel(), tc.Equals, loggo.INFO)
	c.Assert(logger.EffectiveLogLevel(), tc.Equals, loggo.INFO)
	c.Assert(logger.IsErrorEnabled(), tc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), tc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), tc.Equals, true)
	c.Assert(logger.IsDebugEnabled(), tc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), tc.Equals, false)
	logger.SetLogLevel(loggo.WARNING)
	c.Assert(logger.LogLevel(), tc.Equals, loggo.WARNING)
	c.Assert(logger.EffectiveLogLevel(), tc.Equals, loggo.WARNING)
	c.Assert(logger.IsErrorEnabled(), tc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), tc.Equals, true)
	c.Assert(logger.IsInfoEnabled(), tc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), tc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), tc.Equals, false)
	logger.SetLogLevel(loggo.ERROR)
	c.Assert(logger.LogLevel(), tc.Equals, loggo.ERROR)
	c.Assert(logger.EffectiveLogLevel(), tc.Equals, loggo.ERROR)
	c.Assert(logger.IsErrorEnabled(), tc.Equals, true)
	c.Assert(logger.IsWarningEnabled(), tc.Equals, false)
	c.Assert(logger.IsInfoEnabled(), tc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), tc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), tc.Equals, false)
	// This is added for completeness, but not really expected to be used.
	logger.SetLogLevel(loggo.CRITICAL)
	c.Assert(logger.LogLevel(), tc.Equals, loggo.CRITICAL)
	c.Assert(logger.EffectiveLogLevel(), tc.Equals, loggo.CRITICAL)
	c.Assert(logger.IsErrorEnabled(), tc.Equals, false)
	c.Assert(logger.IsWarningEnabled(), tc.Equals, false)
	c.Assert(logger.IsInfoEnabled(), tc.Equals, false)
	c.Assert(logger.IsDebugEnabled(), tc.Equals, false)
	c.Assert(logger.IsTraceEnabled(), tc.Equals, false)
	logger.SetLogLevel(loggo.UNSPECIFIED)
	c.Assert(logger.LogLevel(), tc.Equals, loggo.UNSPECIFIED)
	c.Assert(logger.EffectiveLogLevel(), tc.Equals, loggo.WARNING)
}

func (s *LoggerSuite) TestModuleLowered(c *tc.C) {
	logger1 := loggo.GetLogger("TESTING.MODULE")
	logger2 := loggo.GetLogger("Testing")

	c.Assert(logger1.Name(), tc.Equals, "testing.module")
	c.Assert(logger2.Name(), tc.Equals, "testing")
}

func (s *LoggerSuite) TestLevelsInherited(c *tc.C) {
	root := loggo.GetLogger("")
	first := loggo.GetLogger("first")
	second := loggo.GetLogger("first.second")

	root.SetLogLevel(loggo.ERROR)
	c.Assert(root.LogLevel(), tc.Equals, loggo.ERROR)
	c.Assert(root.EffectiveLogLevel(), tc.Equals, loggo.ERROR)
	c.Assert(first.LogLevel(), tc.Equals, loggo.UNSPECIFIED)
	c.Assert(first.EffectiveLogLevel(), tc.Equals, loggo.ERROR)
	c.Assert(second.LogLevel(), tc.Equals, loggo.UNSPECIFIED)
	c.Assert(second.EffectiveLogLevel(), tc.Equals, loggo.ERROR)

	first.SetLogLevel(loggo.DEBUG)
	c.Assert(root.LogLevel(), tc.Equals, loggo.ERROR)
	c.Assert(root.EffectiveLogLevel(), tc.Equals, loggo.ERROR)
	c.Assert(first.LogLevel(), tc.Equals, loggo.DEBUG)
	c.Assert(first.EffectiveLogLevel(), tc.Equals, loggo.DEBUG)
	c.Assert(second.LogLevel(), tc.Equals, loggo.UNSPECIFIED)
	c.Assert(second.EffectiveLogLevel(), tc.Equals, loggo.DEBUG)

	second.SetLogLevel(loggo.INFO)
	c.Assert(root.LogLevel(), tc.Equals, loggo.ERROR)
	c.Assert(root.EffectiveLogLevel(), tc.Equals, loggo.ERROR)
	c.Assert(first.LogLevel(), tc.Equals, loggo.DEBUG)
	c.Assert(first.EffectiveLogLevel(), tc.Equals, loggo.DEBUG)
	c.Assert(second.LogLevel(), tc.Equals, loggo.INFO)
	c.Assert(second.EffectiveLogLevel(), tc.Equals, loggo.INFO)

	first.SetLogLevel(loggo.UNSPECIFIED)
	c.Assert(root.LogLevel(), tc.Equals, loggo.ERROR)
	c.Assert(root.EffectiveLogLevel(), tc.Equals, loggo.ERROR)
	c.Assert(first.LogLevel(), tc.Equals, loggo.UNSPECIFIED)
	c.Assert(first.EffectiveLogLevel(), tc.Equals, loggo.ERROR)
	c.Assert(second.LogLevel(), tc.Equals, loggo.INFO)
	c.Assert(second.EffectiveLogLevel(), tc.Equals, loggo.INFO)
}

func (s *LoggerSuite) TestParent(c *tc.C) {
	logger := loggo.GetLogger("a.b.c")
	b := logger.Parent()
	a := b.Parent()
	root := a.Parent()

	c.Check(b.Name(), tc.Equals, "a.b")
	c.Check(a.Name(), tc.Equals, "a")
	c.Check(root.Name(), tc.Equals, "<root>")
	c.Check(root.Parent(), tc.DeepEquals, root)
}

func (s *LoggerSuite) TestParentSameContext(c *tc.C) {
	ctx := loggo.NewContext(loggo.DEBUG)

	logger := ctx.GetLogger("a.b.c")
	b := logger.Parent()

	c.Check(b, tc.DeepEquals, ctx.GetLogger("a.b"))
	c.Check(b, tc.Not(tc.DeepEquals), loggo.GetLogger("a.b"))
}

func (s *LoggerSuite) TestChild(c *tc.C) {
	root := loggo.GetLogger("")

	a := root.Child("a")
	logger := a.Child("b.c")

	c.Check(a.Name(), tc.Equals, "a")
	c.Check(logger.Name(), tc.Equals, "a.b.c")
	c.Check(logger.Parent(), tc.DeepEquals, a.Child("b"))
}

func (s *LoggerSuite) TestChildSameContext(c *tc.C) {
	ctx := loggo.NewContext(loggo.DEBUG)

	logger := ctx.GetLogger("a")
	b := logger.Child("b")

	c.Check(b, tc.DeepEquals, ctx.GetLogger("a.b"))
	c.Check(b, tc.Not(tc.DeepEquals), loggo.GetLogger("a.b"))
}

func (s *LoggerSuite) TestChildSameContextWithTags(c *tc.C) {
	ctx := loggo.NewContext(loggo.DEBUG)

	logger := ctx.GetLogger("a", "parent")
	b := logger.ChildWithTags("b", "child")

	c.Check(ctx.GetAllLoggerTags(), tc.DeepEquals, []string{"child", "parent"})
	c.Check(logger.Tags(), tc.DeepEquals, []string{"parent"})
	c.Check(b.Tags(), tc.DeepEquals, []string{"child"})
}

func (s *LoggerSuite) TestRoot(c *tc.C) {
	logger := loggo.GetLogger("a.b.c")
	root := logger.Root()

	c.Check(root.Name(), tc.Equals, "<root>")
	c.Check(root.Child("a.b.c"), tc.DeepEquals, logger)
}

func (s *LoggerSuite) TestRootSameContext(c *tc.C) {
	ctx := loggo.NewContext(loggo.DEBUG)

	logger := ctx.GetLogger("a.b.c")
	root := logger.Root()

	c.Check(root.Name(), tc.Equals, "<root>")
	c.Check(root.Child("a.b.c"), tc.DeepEquals, logger)
	c.Check(root.Child("a.b.c"), tc.Not(tc.DeepEquals), loggo.GetLogger("a.b.c"))
}

func (s *LoggerSuite) TestAttrsIgnoredForMessageFormatting(c *tc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, tc.IsNil)

	logger := context.GetLogger("testing")

	// Attrs should not be consumed as format arguments.
	_ = logger.Logf(c.Context(), loggo.INFO, "hello world", attrs.String("key", "value"))

	logs := writer.Log()
	c.Assert(logs, tc.HasLen, 1)
	c.Check(logs[0].Message, tc.Equals, "hello world")
	c.Assert(logs[0].Attrs, tc.HasLen, 1)
	c.Check(logs[0].Attrs[0].(attrs.AttrValue[string]).Key(), tc.Equals, "key")
	c.Check(logs[0].Attrs[0].(attrs.AttrValue[string]).Value(), tc.Equals, "value")
}

func (s *LoggerSuite) TestAttrsIgnoredWithFormatArgs(c *tc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, tc.IsNil)

	logger := context.GetLogger("testing")

	// Format args should be used for formatting, attrs should be separated.
	_ = logger.Logf(c.Context(), loggo.INFO, "count: %d", 42, attrs.Int("extra", 7))

	logs := writer.Log()
	c.Assert(logs, tc.HasLen, 1)
	c.Check(logs[0].Message, tc.Equals, "count: 42")
	c.Assert(logs[0].Attrs, tc.HasLen, 1)
	c.Check(logs[0].Attrs[0].(attrs.AttrValue[int]).Key(), tc.Equals, "extra")
	c.Check(logs[0].Attrs[0].(attrs.AttrValue[int]).Value(), tc.Equals, 7)
}

func (s *LoggerSuite) TestMultipleAttrsIgnoredForMessageFormatting(c *tc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, tc.IsNil)

	logger := context.GetLogger("testing")

	// Multiple attrs should all be separated from the message.
	_ = logger.Logf(c.Context(), loggo.INFO, "msg %s", "formatted",
		attrs.String("s", "val"),
		attrs.Int("i", 1),
		attrs.Bool("b", true),
	)

	logs := writer.Log()
	c.Assert(logs, tc.HasLen, 1)
	c.Check(logs[0].Message, tc.Equals, "msg formatted")
	c.Assert(logs[0].Attrs, tc.HasLen, 3)
	c.Check(logs[0].Attrs[0].(attrs.AttrValue[string]).Key(), tc.Equals, "s")
	c.Check(logs[0].Attrs[1].(attrs.AttrValue[int]).Key(), tc.Equals, "i")
	c.Check(logs[0].Attrs[2].(attrs.AttrValue[bool]).Key(), tc.Equals, "b")
}

func (s *LoggerSuite) TestAttrsOnlyNoFormatArgs(c *tc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, tc.IsNil)

	logger := context.GetLogger("testing")

	// When only attrs are passed with no format verbs, message stays as-is.
	_ = logger.Logf(c.Context(), loggo.INFO, "plain message",
		attrs.Float64("latency", 1.23),
		attrs.Duration("elapsed", 5*time.Second),
	)

	logs := writer.Log()
	c.Assert(logs, tc.HasLen, 1)
	c.Check(logs[0].Message, tc.Equals, "plain message")
	c.Assert(logs[0].Attrs, tc.HasLen, 2)
	c.Check(logs[0].Attrs[0].(attrs.AttrValue[float64]).Key(), tc.Equals, "latency")
	c.Check(logs[0].Attrs[0].(attrs.AttrValue[float64]).Value(), tc.Equals, 1.23)
	c.Check(logs[0].Attrs[1].(attrs.AttrValue[time.Duration]).Key(), tc.Equals, "elapsed")
	c.Check(logs[0].Attrs[1].(attrs.AttrValue[time.Duration]).Value(), tc.Equals, 5*time.Second)
}

func (s *LoggerSuite) TestNoAttrsFormatsNormally(c *tc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, tc.IsNil)

	logger := context.GetLogger("testing")

	// Without attrs, formatting works as usual.
	_ = logger.Logf(c.Context(), loggo.INFO, "hello %s %d", "world", 42)

	logs := writer.Log()
	c.Assert(logs, tc.HasLen, 1)
	c.Check(logs[0].Message, tc.Equals, "hello world 42")
	c.Check(logs[0].Attrs, tc.HasLen, 0)
}

func (s *LoggerSuite) TestAttrsInterleavedWithFormatArgs(c *tc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, tc.IsNil)

	logger := context.GetLogger("testing")

	// Attrs interleaved with format args: attrs are separated out,
	// format args are used in order.
	_ = logger.Logf(c.Context(), loggo.INFO, "%s=%d",
		attrs.String("before", "x"),
		"key",
		attrs.Int("middle", 5),
		100,
		attrs.Bool("after", false),
	)

	logs := writer.Log()
	c.Assert(logs, tc.HasLen, 1)
	c.Check(logs[0].Message, tc.Equals, "key=100")
	c.Assert(logs[0].Attrs, tc.HasLen, 3)
	c.Check(logs[0].Attrs[0].(attrs.AttrValue[string]).Key(), tc.Equals, "before")
	c.Check(logs[0].Attrs[1].(attrs.AttrValue[int]).Key(), tc.Equals, "middle")
	c.Check(logs[0].Attrs[2].(attrs.AttrValue[bool]).Key(), tc.Equals, "after")
}
