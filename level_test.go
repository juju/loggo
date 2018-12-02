// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
)

type LevelSuite struct{}

var _ = gc.Suite(&LevelSuite{})

var parseLevelTests = []struct {
	str   string
	level loggo.Level
	fail  bool
}{{
	str:   "trace",
	level: loggo.TRACE,
}, {
	str:   "TrAce",
	level: loggo.TRACE,
}, {
	str:   "TRACE",
	level: loggo.TRACE,
}, {
	str:   "debug",
	level: loggo.DEBUG,
}, {
	str:   "DEBUG",
	level: loggo.DEBUG,
}, {
	str:   "info",
	level: loggo.INFO,
}, {
	str:   "INFO",
	level: loggo.INFO,
}, {
	str:   "notE",
	level: loggo.NOTICE,
}, {
	str:   "NoTICE",
	level: loggo.NOTICE,
}, {
	str:   "warn",
	level: loggo.WARNING,
}, {
	str:   "WARN",
	level: loggo.WARNING,
}, {
	str:   "warning",
	level: loggo.WARNING,
}, {
	str:   "WARNING",
	level: loggo.WARNING,
}, {
	str:   "error",
	level: loggo.ERROR,
}, {
	str:   "ERROR",
	level: loggo.ERROR,
}, {
	str:   "critical",
	level: loggo.CRITICAL,
}, {
	str:   "Alert",
	level: loggo.ALERT,
}, {
	str:   "EMERG",
	level: loggo.EMERGENCY,
}, {
	str:   "EMERGENCY",
	level: loggo.EMERGENCY,
}, {
	str:   "Emerg",
	level: loggo.EMERGENCY,
}, {
	str:   "not_specified",
	level: loggo.UNSPECIFIED,
	fail:  true,
}, {
	str:   "other",
	level: loggo.UNSPECIFIED,
	fail:  true,
}, {
	str:   "",
	level: loggo.UNSPECIFIED,
	fail:  true,
}}

func (s *LevelSuite) TestParseLevel(c *gc.C) {
	for _, test := range parseLevelTests {
		level, ok := loggo.ParseLevel(test.str)
		c.Logf("str=%s, level=%v, ok=%v", test.str, level, ok)
		c.Assert(level, gc.Equals, test.level)
		c.Assert(ok, gc.Equals, !test.fail)
	}
}

var levelStringValueTests = map[loggo.Level]string{
	loggo.UNSPECIFIED: "UNSPECIFIED",
	loggo.DEBUG:       "DEBUG",
	loggo.TRACE:       "TRACE",
	loggo.INFO:        "INFO",
	loggo.NOTICE:      "NOTICE",
	loggo.WARNING:     "WARNING",
	loggo.ERROR:       "ERROR",
	loggo.CRITICAL:    "CRITICAL",
	loggo.ALERT:       "ALERT",
	loggo.EMERGENCY:   "EMERGENCY",
	loggo.Level(42):   "UNKNOWN", // other values are unknown
}

func (s *LevelSuite) TestLevelStringValue(c *gc.C) {
	for level, str := range levelStringValueTests {
		c.Assert(level.String(), gc.Equals, str)
	}
}
