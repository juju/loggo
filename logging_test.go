// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"testing"
	"time"

	"github.com/juju/loggo/v3"
	"github.com/juju/tc"
)

type LoggingSuite struct {
	context *loggo.Context
	writer  *writer
	logger  loggo.Logger

	// Test that labels get outputted to loggo.Entry
	Labels map[string]string
}

func TestLoggingSuite(t *testing.T) {
	tc.Run(t, &LoggingSuite{})
	tc.Run(t, &LoggingSuite{Labels: loggo.Labels{"logger-tags": "ONE,TWO"}})
}

func (s *LoggingSuite) SetUpTest(c *tc.C) {
	s.writer = &writer{}
	s.context = loggo.NewContext(loggo.TRACE)
	err := s.context.AddWriter("test", s.writer)
	c.Assert(err, tc.IsNil)
	s.logger = s.context.GetLogger("test", "ONE,TWO")
}

func (s *LoggingSuite) TestLoggingStrings(c *tc.C) {
	_ = s.logger.Infof(c.Context(), "simple")
	_ = s.logger.Infof(c.Context(), "with args %d", 42)
	_ = s.logger.Infof(c.Context(), "working 100%")
	_ = s.logger.Infof(c.Context(), "missing %s")

	checkLogEntries(c, s.writer.Log(), []loggo.Entry{
		{Level: loggo.INFO, Module: "test", Message: "simple", Labels: s.Labels},
		{Level: loggo.INFO, Module: "test", Message: "with args 42", Labels: s.Labels},
		{Level: loggo.INFO, Module: "test", Message: "working 100%", Labels: s.Labels},
		{Level: loggo.INFO, Module: "test", Message: "missing %s", Labels: s.Labels},
	})
}

func (s *LoggingSuite) TestLoggingLimitWarning(c *tc.C) {
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

func (s *LoggingSuite) TestLocationCapture(c *tc.C) {
	s.helperInfof(c, "helper message")                                              //tag helper-location
	_ = s.logger.Criticalf(c.Context(), "critical message")                         //tag critical-location
	_ = s.logger.Errorf(c.Context(), "error message")                               //tag error-location
	_ = s.logger.Warningf(c.Context(), "warning message")                           //tag warning-location
	_ = s.logger.Infof(c.Context(), "info message")                                 //tag info-location
	_ = s.logger.Debugf(c.Context(), "debug message")                               //tag debug-location
	_ = s.logger.Tracef(c.Context(), "trace message")                               //tag trace-location
	_ = s.logger.Logf(c.Context(), loggo.INFO, "logf msg")                          //tag logf-location
	_ = s.logger.LogCallf(c.Context(), 1, loggo.INFO, "logcallf msg")               //tag logcallf-location
	_ = s.logger.LogWithLabelsf(c.Context(), loggo.INFO, "logwithlabelsf msg", nil) //tag logwithlabelsf-location

	log := s.writer.Log()
	tags := []string{
		"helper-location",
		"critical-location",
		"error-location",
		"warning-location",
		"info-location",
		"debug-location",
		"trace-location",
		"logf-location",
		"logcallf-location",
		"logwithlabelsf-location",
	}
	c.Assert(log, tc.HasLen, len(tags))
	for x := range tags {
		assertLocation(c, log[x], tags[x])
	}
}

func (s *LoggingSuite) helperInfof(c *tc.C, format string, args ...any) {
	s.logger.Helper()
	_ = s.logger.Infof(c.Context(), format, args...)
}

func (s *LoggingSuite) TestLogDoesntLogWeirdLevels(c *tc.C) {
	_ = s.logger.Logf(c.Context(), loggo.UNSPECIFIED, "message")
	c.Assert(s.writer.Log(), tc.HasLen, 0)

	_ = s.logger.Logf(c.Context(), loggo.Level(42), "message")
	c.Assert(s.writer.Log(), tc.HasLen, 0)

	_ = s.logger.Logf(c.Context(), loggo.CRITICAL+loggo.Level(1), "message")
	c.Assert(s.writer.Log(), tc.HasLen, 0)
}
