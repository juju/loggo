// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	gc "gopkg.in/check.v1"

	"github.com/juju/loggo/v2"
)

type LoggerSuite struct{}

var _ = gc.Suite(&LoggerSuite{})

func (*LoggerSuite) SetUpTest(c *gc.C) {
	loggo.ResetDefaultContext()
}

func (s *LoggerSuite) TestRootLogger(c *gc.C) {
	root := loggo.Logger{}
	c.Check(root.Name(), gc.Equals, "<root>")
	c.Check(root.LogLevel(), gc.Equals, loggo.WARNING)
	c.Check(root.IsErrorEnabled(), gc.Equals, true)
	c.Check(root.IsWarningEnabled(), gc.Equals, true)
	c.Check(root.IsInfoEnabled(), gc.Equals, false)
	c.Check(root.IsDebugEnabled(), gc.Equals, false)
	c.Check(root.IsTraceEnabled(), gc.Equals, false)
}

func (s *LoggerSuite) TestInheritedLabels(c *gc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, gc.IsNil)

	logger := context.GetLogger("testing")

	nestedLoggerWithLabels := logger.
		ChildWithLabels("nested", loggo.Labels{"foo": "bar"})
	deepNestedLoggerWithLabels := nestedLoggerWithLabels.
		ChildWithLabels("nested", loggo.Labels{"foo": "baz"}).
		ChildWithLabels("deepnested", loggo.Labels{"fred": "tim"})

	loggerWithTagsAndLabels := logger.
		ChildWithLabels("nested-labels", loggo.Labels{"hello": "world"}).
		ChildWithTags("nested-tag", "tag1", "tag2")

	c.Check(nestedLoggerWithLabels.Labels(), gc.DeepEquals, loggo.Labels{"foo": "bar"})
	c.Check(deepNestedLoggerWithLabels.Labels(), gc.DeepEquals, loggo.Labels{"foo": "baz", "fred": "tim"})
	c.Check(loggerWithTagsAndLabels.Labels(), gc.DeepEquals, loggo.Labels{"hello": "world"})
}

func (s *LoggerSuite) TestInheritedLabelsInLogs(c *gc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, gc.IsNil)

	logger := context.GetLogger("testing")

	nestedLoggerWithLabels := logger.
		ChildWithLabels("nested", loggo.Labels{"foo": "bar"})
	deepNestedLoggerWithLabels := nestedLoggerWithLabels.
		ChildWithLabels("nested", loggo.Labels{"foo": "bar"}).
		ChildWithLabels("deepnested", loggo.Labels{"fred": "tim"})

	loggerWithTagsAndLabels := logger.
		ChildWithLabels("nested-labels", loggo.Labels{"hello": "world"}).
		ChildWithTags("nested-tag", "tag1", "tag2")

	logger.Logf(loggo.INFO, "without labels")
	nestedLoggerWithLabels.Logf(loggo.INFO, "with nested labels")
	deepNestedLoggerWithLabels.Logf(loggo.INFO, "with deep nested labels")
	loggerWithTagsAndLabels.Logf(loggo.INFO, "with tags and labels")

	logs := writer.Log()
	c.Assert(logs, gc.HasLen, 4)
	c.Check(logs[0].Message, gc.Equals, "without labels")
	c.Check(logs[0].Labels, gc.HasLen, 0)

	c.Check(logs[1].Message, gc.Equals, "with nested labels")
	c.Check(logs[1].Labels, gc.DeepEquals, loggo.Labels{"foo": "bar"})

	c.Check(logs[2].Message, gc.Equals, "with deep nested labels")
	c.Check(logs[2].Labels, gc.DeepEquals, loggo.Labels{
		"foo":  "bar",
		"fred": "tim",
	})

	c.Check(logs[3].Message, gc.Equals, "with tags and labels")
	c.Check(logs[3].Labels, gc.DeepEquals, loggo.Labels{
		"logger-tags": "tag1,tag2",
		"hello":       "world",
	})
}

func (s *LoggerSuite) TestInheritedLabelsResetLabels(c *gc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, gc.IsNil)

	logger := context.GetLogger("testing")

	nestedLoggerWithLabels := logger.
		ChildWithLabels("nested", loggo.Labels{"foo": "bar"})
	deepNestedLoggerWithLabels := nestedLoggerWithLabels.
		ChildWithLabels("nested", loggo.Labels{"foo": "bar"}).
		ChildWithLabels("deepnested", loggo.Labels{"fred": "tim"})

	// Ensure we can set the logger level to error.

	deepNestedLoggerWithLabels.SetLogLevel(loggo.ERROR)
	c.Check(deepNestedLoggerWithLabels.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Check(deepNestedLoggerWithLabels.LogLevel(), gc.Equals, loggo.ERROR)

	logger.Logf(loggo.INFO, "without labels")
	nestedLoggerWithLabels.Logf(loggo.INFO, "with nested labels")
	deepNestedLoggerWithLabels.Logf(loggo.INFO, "with deep nested labels")

	logs := writer.Log()
	c.Assert(logs, gc.HasLen, 2, gc.Commentf("len = %d", len(logs)))
	c.Check(logs[0].Message, gc.Equals, "without labels")
	c.Check(logs[0].Labels, gc.HasLen, 0)

	c.Check(logs[1].Message, gc.Equals, "with nested labels")
	c.Check(logs[1].Labels, gc.DeepEquals, loggo.Labels{"foo": "bar"})

	writer.Clear()

	// Ensure that we can reset the logger levels.

	context.ResetLoggerLevels()

	c.Check(deepNestedLoggerWithLabels.EffectiveLogLevel(), gc.Equals, loggo.WARNING)
	c.Check(deepNestedLoggerWithLabels.LogLevel(), gc.Equals, loggo.UNSPECIFIED)

	logger.Logf(loggo.INFO, "without labels")
	nestedLoggerWithLabels.Logf(loggo.INFO, "with nested labels")
	deepNestedLoggerWithLabels.Logf(loggo.INFO, "with deep nested labels")

	logs = writer.Log()
	c.Assert(logs, gc.HasLen, 0, gc.Commentf("len = %d", len(logs)))

	writer.Clear()

	// Ensure that we can configure the logger levels again.

	context.ResetLoggerLevels()
	context.ConfigureLoggers("testing.nested.nested.deepnested=INFO")

	c.Check(deepNestedLoggerWithLabels.EffectiveLogLevel(), gc.Equals, loggo.INFO)
	c.Check(deepNestedLoggerWithLabels.LogLevel(), gc.Equals, loggo.INFO)

	logger.Logf(loggo.INFO, "without labels")
	nestedLoggerWithLabels.Logf(loggo.INFO, "with nested labels")
	deepNestedLoggerWithLabels.Logf(loggo.INFO, "with deep nested labels")

	logs = writer.Log()
	c.Assert(logs, gc.HasLen, 1, gc.Commentf("len = %d", len(logs)))

	c.Check(logs[0].Message, gc.Equals, "with deep nested labels")
	c.Check(logs[0].Labels, gc.DeepEquals, loggo.Labels{
		"foo":  "bar",
		"fred": "tim",
	})
}

func (s *LoggerSuite) TestInheritedLabelsConfigByLabels(c *gc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, gc.IsNil)

	logger := context.GetLogger("testing")

	nestedLoggerWithLabels := logger.
		ChildWithLabels("nested", loggo.Labels{"foo": "bar"})
	deepNestedLoggerWithLabels := nestedLoggerWithLabels.
		ChildWithLabels("nested", loggo.Labels{"foo": "bar"}).
		ChildWithLabels("deepnested", loggo.Labels{"fred": "tim"})

	// Apply the ERROR level to the logger with the labels "foo=bar".

	context.ResetLoggerLevels()
	context.ConfigureLoggers("testing=INFO")
	context.ConfigureLoggers("testing.nested=ERROR", loggo.Labels{"foo": "bar"})

	c.Check(nestedLoggerWithLabels.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Check(nestedLoggerWithLabels.LogLevel(), gc.Equals, loggo.ERROR)

	c.Check(deepNestedLoggerWithLabels.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
	c.Check(deepNestedLoggerWithLabels.LogLevel(), gc.Equals, loggo.UNSPECIFIED)

	logger.Logf(loggo.INFO, "without labels")
	nestedLoggerWithLabels.Logf(loggo.INFO, "with nested labels")
	deepNestedLoggerWithLabels.Logf(loggo.INFO, "with deep nested labels")

	logs := writer.Log()
	c.Assert(logs, gc.HasLen, 1, gc.Commentf("len = %d", len(logs)))

	c.Check(logs[0].Message, gc.Equals, "without labels")
	c.Check(logs[0].Labels, gc.HasLen, 0)

	writer.Clear()

	// Apply the INFO level to the logger with the labels "foo=bar" and
	// "fred=tim".

	run := func(labels ...loggo.Labels) {
		context.ConfigureLoggers("testing.nested.nested.deepnested=INFO", labels...)

		c.Check(nestedLoggerWithLabels.EffectiveLogLevel(), gc.Equals, loggo.ERROR)
		c.Check(nestedLoggerWithLabels.LogLevel(), gc.Equals, loggo.ERROR)

		c.Check(deepNestedLoggerWithLabels.EffectiveLogLevel(), gc.Equals, loggo.INFO)
		c.Check(deepNestedLoggerWithLabels.LogLevel(), gc.Equals, loggo.INFO)

		logger.Logf(loggo.INFO, "without labels")
		nestedLoggerWithLabels.Logf(loggo.INFO, "with nested labels")
		deepNestedLoggerWithLabels.Logf(loggo.INFO, "with deep nested labels")

		logs = writer.Log()
		c.Assert(logs, gc.HasLen, 2, gc.Commentf("len = %d", len(logs)))
		c.Check(logs[0].Message, gc.Equals, "without labels")
		c.Check(logs[0].Labels, gc.HasLen, 0)

		c.Check(logs[1].Message, gc.Equals, "with deep nested labels")
		c.Check(logs[1].Labels, gc.DeepEquals, loggo.Labels{
			"foo":  "bar",
			"fred": "tim",
		})

		writer.Clear()
	}

	// Notice that the order of the labels are merged and last label wins.

	run(loggo.Labels{"foo": "bar", "fred": "tim"})
	run(loggo.Labels{"foo": "bar"}, loggo.Labels{"fred": "tim"})
	run(loggo.Labels{"foo": "bar", "fred": "john"}, loggo.Labels{"fred": "tim"})
}

func (s *LoggerSuite) TestInheritedLabelsConfigByLabelsIgnoresInvalidLabels(c *gc.C) {
	writer := &loggo.TestWriter{}
	context := loggo.NewContext(loggo.INFO)
	err := context.AddWriter("test", writer)
	c.Assert(err, gc.IsNil)

	logger := context.GetLogger("testing")

	nestedLoggerWithLabels := logger.
		ChildWithLabels("nested", loggo.Labels{"foo": "bar"})
	deepNestedLoggerWithLabels := nestedLoggerWithLabels.
		ChildWithLabels("nested", loggo.Labels{"foo": "bar"}).
		ChildWithLabels("deepnested", loggo.Labels{"fred": "tim"})

	// Apply the ERROR level to the logger with the labels "hello=world".
	// Nothing matches this, so it should have no effect.

	context.ConfigureLoggers("testing.nested.nested.deepnested=ERROR", loggo.Labels{
		"hello": "world",
	})

	logger.Logf(loggo.INFO, "without labels")
	nestedLoggerWithLabels.Logf(loggo.INFO, "with nested labels")
	deepNestedLoggerWithLabels.Logf(loggo.INFO, "with deep nested labels")

	logs := writer.Log()
	c.Assert(logs, gc.HasLen, 3, gc.Commentf("len = %d", len(logs)))
	c.Check(logs[0].Message, gc.Equals, "without labels")
	c.Check(logs[0].Labels, gc.HasLen, 0)

	c.Check(logs[1].Message, gc.Equals, "with nested labels")
	c.Check(logs[1].Labels, gc.DeepEquals, loggo.Labels{"foo": "bar"})

	c.Check(logs[2].Message, gc.Equals, "with deep nested labels")
	c.Check(logs[2].Labels, gc.DeepEquals, loggo.Labels{
		"foo":  "bar",
		"fred": "tim",
	})

	writer.Clear()

	// Apply the ERROR level to the logger with the labels "foo=bar",
	// "fred=tim" and "hello=world".
	// This should not match anything and have no effect, as we expect
	// to match ALL the labels.

	context.ConfigureLoggers("testing.nested.nested.deepnested=ERROR", loggo.Labels{
		"foo":   "bar",
		"fred":  "tim",
		"hello": "world",
	})

	logger.Logf(loggo.INFO, "without labels")
	nestedLoggerWithLabels.Logf(loggo.INFO, "with nested labels")
	deepNestedLoggerWithLabels.Logf(loggo.INFO, "with deep nested labels")

	logs = writer.Log()
	c.Assert(logs, gc.HasLen, 3, gc.Commentf("len = %d", len(logs)))
	c.Check(logs[0].Message, gc.Equals, "without labels")
	c.Check(logs[0].Labels, gc.HasLen, 0)

	c.Check(logs[1].Message, gc.Equals, "with nested labels")
	c.Check(logs[1].Labels, gc.DeepEquals, loggo.Labels{"foo": "bar"})

	c.Check(logs[2].Message, gc.Equals, "with deep nested labels")
	c.Check(logs[2].Labels, gc.DeepEquals, loggo.Labels{
		"foo":  "bar",
		"fred": "tim",
	})
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

func (s *LoggerSuite) TestParent(c *gc.C) {
	logger := loggo.GetLogger("a.b.c")
	b := logger.Parent()
	a := b.Parent()
	root := a.Parent()

	c.Check(b.Name(), gc.Equals, "a.b")
	c.Check(a.Name(), gc.Equals, "a")
	c.Check(root.Name(), gc.Equals, "<root>")
	c.Check(root.Parent(), gc.DeepEquals, root)
}

func (s *LoggerSuite) TestParentSameContext(c *gc.C) {
	ctx := loggo.NewContext(loggo.DEBUG)

	logger := ctx.GetLogger("a.b.c")
	b := logger.Parent()

	c.Check(b, gc.DeepEquals, ctx.GetLogger("a.b"))
	c.Check(b, gc.Not(gc.DeepEquals), loggo.GetLogger("a.b"))
}

func (s *LoggerSuite) TestChild(c *gc.C) {
	root := loggo.GetLogger("")

	a := root.Child("a")
	logger := a.Child("b.c")

	c.Check(a.Name(), gc.Equals, "a")
	c.Check(logger.Name(), gc.Equals, "a.b.c")
	c.Check(logger.Parent(), gc.DeepEquals, a.Child("b"))
}

func (s *LoggerSuite) TestChildSameContext(c *gc.C) {
	ctx := loggo.NewContext(loggo.DEBUG)

	logger := ctx.GetLogger("a")
	b := logger.Child("b")

	c.Check(b, gc.DeepEquals, ctx.GetLogger("a.b"))
	c.Check(b, gc.Not(gc.DeepEquals), loggo.GetLogger("a.b"))
}

func (s *LoggerSuite) TestChildSameContextWithTags(c *gc.C) {
	ctx := loggo.NewContext(loggo.DEBUG)

	logger := ctx.GetLogger("a", "parent")
	b := logger.ChildWithTags("b", "child")

	c.Check(ctx.GetAllLoggerTags(), gc.DeepEquals, []string{"child", "parent"})
	c.Check(logger.Tags(), gc.DeepEquals, []string{"parent"})
	c.Check(b.Tags(), gc.DeepEquals, []string{"child"})
}

func (s *LoggerSuite) TestRoot(c *gc.C) {
	logger := loggo.GetLogger("a.b.c")
	root := logger.Root()

	c.Check(root.Name(), gc.Equals, "<root>")
	c.Check(root.Child("a.b.c"), gc.DeepEquals, logger)
}

func (s *LoggerSuite) TestRootSameContext(c *gc.C) {
	ctx := loggo.NewContext(loggo.DEBUG)

	logger := ctx.GetLogger("a.b.c")
	root := logger.Root()

	c.Check(root.Name(), gc.Equals, "<root>")
	c.Check(root.Child("a.b.c"), gc.DeepEquals, logger)
	c.Check(root.Child("a.b.c"), gc.Not(gc.DeepEquals), loggo.GetLogger("a.b.c"))
}
