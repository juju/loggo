// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
)

type SubLoggerSuite struct{}

var _ = gc.Suite(&SubLoggerSuite{})

func (*SubLoggerSuite) SetUpTest(c *gc.C) {
	loggo.ResetLoggers()
}

func (s *SubLoggerSuite) TearDownTest(c *gc.C) {
	loggo.ResetWriters()
}

func (s *SubLoggerSuite) TearDownSuite(c *gc.C) {
	loggo.ResetLoggers()
}

func (s *SubLoggerSuite) TestRootLogger(c *gc.C) {
	root := loggo.SubLogger{}
	c.Check(root.Name(), gc.Equals, "<root>")
	c.Assert(root.IsErrorEnabled(), gc.Equals, true)
	c.Assert(root.IsWarningEnabled(), gc.Equals, true)
	c.Assert(root.IsInfoEnabled(), gc.Equals, false)
	c.Assert(root.IsDebugEnabled(), gc.Equals, false)
	c.Assert(root.IsTraceEnabled(), gc.Equals, false)
}

func (s *SubLoggerSuite) TestModuleName(c *gc.C) {
	var parent loggo.SubLogger
	logger, _ := loggo.NewSubLogger("loggo.testing", parent)
	c.Assert(logger.Name(), gc.Equals, "loggo.testing")
}

func (s *SubLoggerSuite) TestModuleLowered(c *gc.C) {
	var parent loggo.SubLogger
	logger1, _ := loggo.NewSubLogger("TESTING.MODULE", parent)
	logger2, _ := loggo.NewSubLogger("Testing", parent)

	c.Assert(logger1.Name(), gc.Equals, "testing.module")
	c.Assert(logger2.Name(), gc.Equals, "testing")
}

func (s *SubLoggerSuite) TestUnspecifiedLevel(c *gc.C) {
	var parent loggo.SubLogger
	logger, _ := loggo.NewSubLogger("...", parent)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.WARNING)

	logger.SetLogLevel(loggo.UNSPECIFIED)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.UNSPECIFIED)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.WARNING)
}

func (s *SubLoggerSuite) TestRootUnspecifiedLevel(c *gc.C) {
	logger, _ := loggo.NewRootLogger()
	c.Assert(logger.LogLevel(), gc.Equals, loggo.WARNING)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.WARNING)

	logger.SetLogLevel(loggo.UNSPECIFIED)
	c.Assert(logger.LogLevel(), gc.Equals, loggo.WARNING)
	c.Assert(logger.EffectiveLogLevel(), gc.Equals, loggo.WARNING)
}

func (s *SubLoggerSuite) TestLevelsInherited(c *gc.C) {
	root, _ := loggo.NewRootLogger()
	first, _ := loggo.NewSubLogger("first", root)
	second, _ := loggo.NewSubLogger("first.second", first)

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
