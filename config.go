// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"fmt"
	"sort"
	"strings"
)

// Config is a mapping of logger module names to logging severity levels.
type Config map[string]Level

// String returns a logger configuration string that may be parsed
// using ParseConfigurationString.
func (c Config) String() string {
	if c == nil {
		return ""
	}
	// output in alphabetical order.
	var names []string
	for name := range c {
		names = append(names, name)
	}
	sort.Strings(names)

	var entries []string
	for _, name := range names {
		level := c[name]
		if name == "" {
			name = rootString
		}
		entry := fmt.Sprintf("%s=%s", name, level)
		entries = append(entries, entry)
	}
	return strings.Join(entries, ";")
}

func parseConfigValue(value string) (string, Level, error) {
	pair := strings.SplitN(value, "=", 2)
	if len(pair) < 2 {
		return "", UNSPECIFIED, fmt.Errorf("config value expected '=', found %q", value)
	}
	name := strings.TrimSpace(pair[0])
	if name == "" {
		return "", UNSPECIFIED, fmt.Errorf("config value %q has missing module name", value)
	}

	if label := extractConfigLabel(name); label != "" {
		if strings.Contains(label, ".") {
			// Show the original name and not text potentially extracted config
			// label.
			return "", UNSPECIFIED, fmt.Errorf("config label should not contain '.', found %q", name)
		}
		// Ensure once the normalised extraction has happened, we put the prefix
		// back on, so that we don't loose the fact that the config is a label.
		//
		// Ideally we would change Config from map[string]Level to
		// map[string]ConfigEntry and then we wouldn't need this step, but that
		// causes lots of issues in Juju directly.
		name = fmt.Sprintf("#%s", label)
	}

	levelStr := strings.TrimSpace(pair[1])
	level, ok := ParseLevel(levelStr)
	if !ok {
		return "", UNSPECIFIED, fmt.Errorf("unknown severity level %q", levelStr)
	}
	if name == rootString {
		name = ""
	}
	return name, level, nil
}

// ParseConfigString parses a logger configuration string into a map of logger
// names and their associated log level. This method is provided to allow
// other programs to pre-validate a configuration string rather than just
// calling ConfigureLoggers.
//
// Logging modules are colon- or semicolon-separated; each module is specified
// as <modulename>=<level>.  White space outside of module names and levels is
// ignored.  The root module is specified with the name "<root>".
//
// As a special case, a log level may be specified on its own.
// This is equivalent to specifying the level of the root module,
// so "DEBUG" is equivalent to `<root>=DEBUG`
//
// An example specification:
//	`<root>=ERROR; foo.bar=WARNING`
//	`[LABEL]=ERROR`
func ParseConfigString(specification string) (Config, error) {
	specification = strings.TrimSpace(specification)
	if specification == "" {
		return nil, nil
	}
	cfg := make(Config)
	if level, ok := ParseLevel(specification); ok {
		cfg[""] = level
		return cfg, nil
	}

	values := strings.FieldsFunc(specification, func(r rune) bool { return r == ';' || r == ':' })
	for _, value := range values {
		name, level, err := parseConfigValue(value)
		if err != nil {
			return nil, err
		}
		cfg[name] = level
	}
	return cfg, nil
}

func extractConfigLabel(s string) string {
	name := strings.TrimSpace(s)
	if len(s) < 2 {
		return ""
	}
	if name[0] == '#' {
		return strings.ToLower(name[1:])
	}
	return ""
}
