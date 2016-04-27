// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
)

type ConfigSuite struct{}

var _ = gc.Suite(&ConfigSuite{})

var configureLoggersTests = []struct {
	spec string
	info string
	err  string
}{{
	spec: "",
	info: "<root>=WARNING",
}, {
	spec: "<root>=UNSPECIFIED",
	info: "<root>=WARNING",
}, {
	spec: "<root>=DEBUG",
	info: "<root>=DEBUG",
}, {
	spec: "TRACE",
	info: "<root>=TRACE",
}, {
	spec: "test.module=debug",
	info: "<root>=WARNING;test.module=DEBUG",
}, {
	spec: "module=info; sub.module=debug; other.module=warning",
	info: "<root>=WARNING;module=INFO;other.module=WARNING;sub.module=DEBUG",
}, {
	spec: "  foo.bar \t\r\n= \t\r\nCRITICAL \t\r\n; \t\r\nfoo \r\t\n = DEBUG",
	info: "<root>=WARNING;foo=DEBUG;foo.bar=CRITICAL",
}, {
	spec: "foo;bar",
	info: "<root>=WARNING",
	err:  `logger specification expected '=', found "foo"`,
}, {
	spec: "=foo",
	info: "<root>=WARNING",
	err:  `logger specification "=foo" has blank name or level`,
}, {
	spec: "foo=",
	info: "<root>=WARNING",
	err:  `logger specification "foo=" has blank name or level`,
}, {
	spec: "=",
	info: "<root>=WARNING",
	err:  `logger specification "=" has blank name or level`,
}, {
	spec: "foo=unknown",
	info: "<root>=WARNING",
	err:  `unknown severity level "unknown"`,
}, {
	// Test that nothing is changed even when the
	// first part of the specification parses ok.
	spec: "module=info; foo=unknown",
	info: "<root>=WARNING",
	err:  `unknown severity level "unknown"`,
}}

func (s *ConfigSuite) TestConfigureLoggers(c *gc.C) {
	for i, test := range configureLoggersTests {
		c.Logf("test %d: %q", i, test.spec)
		loggo.ResetLoggers()
		err := loggo.ConfigureLoggers(test.spec)
		c.Check(loggo.LoggerInfo(), gc.Equals, test.info)
		if test.err != "" {
			c.Assert(err, gc.ErrorMatches, test.err)
			continue
		}
		c.Assert(err, gc.IsNil)

		// Test that it's idempotent.
		err = loggo.ConfigureLoggers(test.spec)
		c.Assert(err, gc.IsNil)
		c.Assert(loggo.LoggerInfo(), gc.Equals, test.info)

		// Test that calling ConfigureLoggers with the
		// output of LoggerInfo works too.
		err = loggo.ConfigureLoggers(test.info)
		c.Assert(err, gc.IsNil)
		c.Assert(loggo.LoggerInfo(), gc.Equals, test.info)
	}
}
