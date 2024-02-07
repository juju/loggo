// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"time"

	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
)

type LoggingSuite struct {
	context *loggo.Context
	writer  *writer
	logger  loggo.Logger

	// Test that labels get outputted to loggo.Entry
	Labels map[string]string
}

var _ = gc.Suite(&LoggingSuite{})
var _ = gc.Suite(&LoggingSuite{Labels: loggo.Labels{"logger-tags": "ONE,TWO"}})

func (s *LoggingSuite) SetUpTest(c *gc.C) {
	s.writer = &writer{}
	s.context = loggo.NewContext(loggo.TRACE)
	err := s.context.AddWriter("test", s.writer)
	c.Assert(err, gc.IsNil)
	s.logger = s.context.GetLogger("test", "ONE,TWO")
}

func (s *LoggingSuite) TestLoggingStrings(c *gc.C) {
	s.logger.Infof("simple")
	s.logger.Infof("with args %d", 42)
	s.logger.Infof("working 100%")
	s.logger.Infof("missing %s")

	checkLogEntries(c, s.writer.Log(), []loggo.Entry{
		{Level: loggo.INFO, Module: "test", Message: "simple", Labels: s.Labels},
		{Level: loggo.INFO, Module: "test", Message: "with args 42", Labels: s.Labels},
		{Level: loggo.INFO, Module: "test", Message: "working 100%", Labels: s.Labels},
		{Level: loggo.INFO, Module: "test", Message: "missing %s", Labels: s.Labels},
	})
}

func (s *LoggingSuite) TestLoggingLimitWarning(c *gc.C) {
	s.logger.SetLogLevel(loggo.WARNING)
	start := time.Now()
	logAllSeverities(s.logger)
	end := time.Now()
	entries := s.writer.Log()
	checkLogEntries(c, entries, []loggo.Entry{
		{Level: loggo.CRITICAL, Module: "test", Message: "something critical", Labels: s.Labels},
		{Level: loggo.ERROR, Module: "test", Message: "an error", Labels: s.Labels},
		{Level: loggo.WARNING, Module: "test", Message: "a warning message", Labels: s.Labels},
	})

	for _, entry := range entries {
		c.Check(entry.Timestamp, Between(start, end))
	}
}

func (s *LoggingSuite) TestLocationCapture(c *gc.C) {
	s.logger.Criticalf("critical message") //tag critical-location
	s.logger.Errorf("error message")       //tag error-location
	s.logger.Warningf("warning message")   //tag warning-location
	s.logger.Infof("info message")         //tag info-location
	s.logger.Debugf("debug message")       //tag debug-location
	s.logger.Tracef("trace message")       //tag trace-location

	log := s.writer.Log()
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

func (s *LoggingSuite) TestLogDoesntLogWeirdLevels(c *gc.C) {
	s.logger.Logf(loggo.UNSPECIFIED, "message")
	c.Assert(s.writer.Log(), gc.HasLen, 0)

	s.logger.Logf(loggo.Level(42), "message")
	c.Assert(s.writer.Log(), gc.HasLen, 0)

	s.logger.Logf(loggo.CRITICAL+loggo.Level(1), "message")
	c.Assert(s.writer.Log(), gc.HasLen, 0)
}
