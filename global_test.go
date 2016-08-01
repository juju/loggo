// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"github.com/juju/loggo"
	gc "gopkg.in/check.v1"
)

type GlobalSuite struct{}

var _ = gc.Suite(&GlobalSuite{})

func (*GlobalSuite) SetUpTest(c *gc.C) {
	loggo.ResetDefaultContext()
}

func (*GlobalSuite) TestRootLogger(c *gc.C) {
	var root loggo.Logger

	got := loggo.GetLogger("")

	c.Check(got.Name(), gc.Equals, root.Name())
	c.Check(got.LogLevel(), gc.Equals, root.LogLevel())
}

func (*GlobalSuite) TestModuleName(c *gc.C) {
	logger := loggo.GetLogger("loggo.testing")
	c.Check(logger.Name(), gc.Equals, "loggo.testing")
}

func (*GlobalSuite) TestLevel(c *gc.C) {
	logger := loggo.GetLogger("testing")
	level := logger.LogLevel()
	c.Check(level, gc.Equals, loggo.UNSPECIFIED)
}

func (*GlobalSuite) TestEffectiveLevel(c *gc.C) {
	logger := loggo.GetLogger("testing")
	level := logger.EffectiveLogLevel()
	c.Check(level, gc.Equals, loggo.WARNING)
}

func (*GlobalSuite) TestLevelsSharedForSameModule(c *gc.C) {
	logger1 := loggo.GetLogger("testing.module")
	logger2 := loggo.GetLogger("testing.module")

	logger1.SetLogLevel(loggo.INFO)
	c.Assert(logger1.IsInfoEnabled(), gc.Equals, true)
	c.Assert(logger2.IsInfoEnabled(), gc.Equals, true)
}

func (*GlobalSuite) TestModuleLowered(c *gc.C) {
	logger1 := loggo.GetLogger("TESTING.MODULE")
	logger2 := loggo.GetLogger("Testing")

	c.Assert(logger1.Name(), gc.Equals, "testing.module")
	c.Assert(logger2.Name(), gc.Equals, "testing")
}

func (s *GlobalSuite) TestConfigureLoggers(c *gc.C) {
	err := loggo.ConfigureLoggers("testing.module=debug")
	c.Assert(err, gc.IsNil)
	expected := "<root>=WARNING;testing.module=DEBUG"
	c.Assert(loggo.DefaultContext().Config().String(), gc.Equals, expected)
	c.Assert(loggo.LoggerInfo(), gc.Equals, expected)
}

func (*GlobalSuite) TestRegisterWriterExistingName(c *gc.C) {
	err := loggo.RegisterWriter("default", &writer{})
	c.Assert(err, gc.ErrorMatches, `context already has a writer named "default"`)
}

func (*GlobalSuite) TestReplaceDefaultWriter(c *gc.C) {
	oldWriter, err := loggo.ReplaceDefaultWriter(&writer{})
	c.Assert(oldWriter, gc.NotNil)
	c.Assert(err, gc.IsNil)
	c.Assert(loggo.DefaultContext().WriterNames(), gc.DeepEquals, []string{"default"})
}

func (*GlobalSuite) TestRemoveWriter(c *gc.C) {
	oldWriter, err := loggo.RemoveWriter("default")
	c.Assert(oldWriter, gc.NotNil)
	c.Assert(err, gc.IsNil)
	c.Assert(loggo.DefaultContext().WriterNames(), gc.HasLen, 0)
}
