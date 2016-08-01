// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import gc "gopkg.in/check.v1"

type ConfigSuite struct{}

var _ = gc.Suite(&ConfigSuite{})

func (*ConfigSuite) TestParseConfigValue(c *gc.C) {
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
	}} {
		c.Logf("%d: %s", i, test.value)
		module, level, err := parseConfigValue(test.value)
		if test.err == "" {
			c.Check(err, gc.IsNil)
			c.Check(module, gc.Equals, test.module)
			c.Check(level, gc.Equals, test.level)
		} else {
			c.Check(module, gc.Equals, "")
			c.Check(level, gc.Equals, UNSPECIFIED)
			c.Check(err.Error(), gc.Equals, test.err)
		}
	}
}

func (*ConfigSuite) TestPaarseConfigurationString(c *gc.C) {
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
			c.Check(err, gc.IsNil)
			c.Check(config, gc.DeepEquals, test.expected)
		} else {
			c.Check(config, gc.IsNil)
			c.Check(err.Error(), gc.Equals, test.err)
		}
	}
}

func (*ConfigSuite) TestConfigString(c *gc.C) {
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
		c.Check(test.config.String(), gc.Equals, test.expected)
	}
}
