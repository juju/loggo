// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"testing"

	"github.com/juju/tc"
)

type ConfigSuite struct{}

func TestConfigSuite(t *testing.T) {
	tc.Run(t, &ConfigSuite{})
}

func (*ConfigSuite) TestParseConfigValue(c *tc.C) {
	for i, test := range []struct {
		value  string
		module string
		level  Level
		err    string
	}{{
		err: `config value expected '=', found ""`,
	}, {
		value: "WARNING",
		err:   `config value expected '=', found "WARNING"`,
	}, {
		value: "=WARNING",
		err:   `config value "=WARNING" has missing module name`,
	}, {
		value:  " name = WARNING ",
		module: "name",
		level:  WARNING,
	}, {
		value: "name = foo",
		err:   `unknown severity level "foo"`,
	}, {
		value: "name=DEBUG=INFO",
		err:   `unknown severity level "DEBUG=INFO"`,
	}, {
		value:  "<root> = info",
		module: "",
		level:  INFO,
	}, {
		value:  "#tag = info",
		module: "#tag",
		level:  INFO,
	}, {
		value:  "#TAG = info",
		module: "#tag",
		level:  INFO,
	}, {
		value: "#tag.1 = info",
		err:   `config tag should not contain '.', found "#tag.1"`,
	}} {
		c.Logf("%d: %s", i, test.value)
		module, level, err := parseConfigValue(test.value)
		if test.err == "" {
			c.Check(err, tc.IsNil)
			c.Check(module, tc.Equals, test.module)
			c.Check(level, tc.Equals, test.level)
		} else {
			c.Check(module, tc.Equals, "")
			c.Check(level, tc.Equals, UNSPECIFIED)
			c.Check(err.Error(), tc.Equals, test.err)
		}
	}
}

func (*ConfigSuite) TestParseConfigurationString(c *tc.C) {
	for i, test := range []struct {
		configuration string
		expected      Config
		err           string
	}{{
		configuration: "",
		// nil Config, no error
	}, {
		configuration: "INFO",
		expected:      Config{"": INFO},
	}, {
		configuration: "=INFO",
		err:           `config value "=INFO" has missing module name`,
	}, {
		configuration: "<root>=UNSPECIFIED",
		expected:      Config{"": UNSPECIFIED},
	}, {
		configuration: "<root>=DEBUG",
		expected:      Config{"": DEBUG},
	}, {
		configuration: "#label=DEBUG",
		expected:      Config{"#label": DEBUG},
	}, {
		configuration: "test.module=debug",
		expected:      Config{"test.module": DEBUG},
	}, {
		configuration: "module=info; sub.module=debug; other.module=warning",
		expected: Config{
			"module":       INFO,
			"sub.module":   DEBUG,
			"other.module": WARNING,
		},
	}, {
		// colons not semicolons
		configuration: "module=info: sub.module=debug: other.module=warning",
		expected: Config{
			"module":       INFO,
			"sub.module":   DEBUG,
			"other.module": WARNING,
		},
	}, {
		configuration: "  foo.bar \t\r\n= \t\r\nCRITICAL \t\r\n; \t\r\nfoo \r\t\n = DEBUG",
		expected: Config{
			"foo":     DEBUG,
			"foo.bar": CRITICAL,
		},
	}, {
		configuration: "foo;bar",
		err:           `config value expected '=', found "foo"`,
	}, {
		configuration: "foo=",
		err:           `unknown severity level ""`,
	}, {
		configuration: "foo=unknown",
		err:           `unknown severity level "unknown"`,
	}} {
		c.Logf("%d: %q", i, test.configuration)
		config, err := ParseConfigString(test.configuration)
		if test.err == "" {
			c.Check(err, tc.IsNil)
			c.Check(config, tc.DeepEquals, test.expected)
		} else {
			c.Check(config, tc.IsNil)
			c.Check(err.Error(), tc.Equals, test.err)
		}
	}
}

func (*ConfigSuite) TestConfigString(c *tc.C) {
	for i, test := range []struct {
		config   Config
		expected string
	}{{
		config:   nil,
		expected: "",
	}, {
		config:   Config{"": INFO},
		expected: "<root>=INFO",
	}, {
		config:   Config{"": UNSPECIFIED},
		expected: "<root>=UNSPECIFIED",
	}, {
		config:   Config{"": DEBUG},
		expected: "<root>=DEBUG",
	}, {
		config:   Config{"test.module": DEBUG},
		expected: "test.module=DEBUG",
	}, {
		config: Config{
			"":             WARNING,
			"module":       INFO,
			"sub.module":   DEBUG,
			"other.module": WARNING,
		},
		expected: "<root>=WARNING;module=INFO;other.module=WARNING;sub.module=DEBUG",
	}} {
		c.Logf("%d: %q", i, test.expected)
		c.Check(test.config.String(), tc.Equals, test.expected)
	}
}
