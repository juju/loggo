// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"testing"

	"github.com/juju/loggo/v2"
	"github.com/juju/tc"
)

type GlobalSuite struct{}

func TestGlobalSuite(t *testing.T) {
	tc.Run(t, &GlobalSuite{})
}

func (*GlobalSuite) SetUpTest(c *tc.C) {
	loggo.ResetDefaultContext()
}

func (*GlobalSuite) TestRootLogger(c *tc.C) {
	var root loggo.Logger

	got := loggo.GetLogger("")

	c.Check(got.Name(), tc.Equals, root.Name())
	c.Check(got.LogLevel(), tc.Equals, root.LogLevel())
}

func (*GlobalSuite) TestModuleName(c *tc.C) {
	logger := loggo.GetLogger("loggo.testing")
	c.Check(logger.Name(), tc.Equals, "loggo.testing")
}

func (*GlobalSuite) TestLevel(c *tc.C) {
	logger := loggo.GetLogger("testing")
	level := logger.LogLevel()
	c.Check(level, tc.Equals, loggo.UNSPECIFIED)
}

func (*GlobalSuite) TestEffectiveLevel(c *tc.C) {
	logger := loggo.GetLogger("testing")
	level := logger.EffectiveLogLevel()
	c.Check(level, tc.Equals, loggo.WARNING)
}

func (*GlobalSuite) TestLevelsSharedForSameModule(c *tc.C) {
	logger1 := loggo.GetLogger("testing.module")
	logger2 := loggo.GetLogger("testing.module")

	logger1.SetLogLevel(loggo.INFO)
	c.Assert(logger1.IsInfoEnabled(), tc.Equals, true)
	c.Assert(logger2.IsInfoEnabled(), tc.Equals, true)
}

func (*GlobalSuite) TestModuleLowered(c *tc.C) {
	logger1 := loggo.GetLogger("TESTING.MODULE")
	logger2 := loggo.GetLogger("Testing")

	c.Assert(logger1.Name(), tc.Equals, "testing.module")
	c.Assert(logger2.Name(), tc.Equals, "testing")
}

func (s *GlobalSuite) TestConfigureLoggers(c *tc.C) {
	err := loggo.ConfigureLoggers("testing.module=debug")
	c.Assert(err, tc.IsNil)
	expected := "<root>=WARNING;testing.module=DEBUG"
	c.Assert(loggo.DefaultContext().Config().String(), tc.Equals, expected)
	c.Assert(loggo.LoggerInfo(), tc.Equals, expected)
}

func (*GlobalSuite) TestRegisterWriterExistingName(c *tc.C) {
	err := loggo.RegisterWriter("default", &writer{})
	c.Assert(err, tc.ErrorMatches, `context already has a writer named "default"`)
}

func (*GlobalSuite) TestReplaceDefaultWriter(c *tc.C) {
	oldWriter, err := loggo.ReplaceDefaultWriter(&writer{})
	c.Assert(oldWriter, tc.NotNil)
	c.Assert(err, tc.IsNil)
	c.Assert(loggo.DefaultContext().WriterNames(), tc.DeepEquals, []string{"default"})
}

func (*GlobalSuite) TestRemoveWriter(c *tc.C) {
	oldWriter, err := loggo.RemoveWriter("default")
	c.Assert(oldWriter, tc.NotNil)
	c.Assert(err, tc.IsNil)
	c.Assert(loggo.DefaultContext().WriterNames(), tc.HasLen, 0)
}

func (s *GlobalSuite) TestGetLoggerWithTags(c *tc.C) {
	logger := loggo.GetLoggerWithTags("parent", "labela", "labelb")
	c.Check(logger.Tags(), tc.DeepEquals, []string{"labela", "labelb"})
}
