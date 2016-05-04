// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"fmt"
	"sort"
	"strings"
)

// loggerInfo returns information about the given loggers and their
// logging levels. The information is returned in the format expected
// by ConfigureLoggers. Loggers with UNSPECIFIED level will not
// be included.
func loggerInfo(modules map[string]*module) string {
	output := []string{}
	// output in alphabetical order.
	keys := []string{}
	for key := range modules {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, name := range keys {
		mod := modules[name]
		severity := mod.level.get()
		if severity == UNSPECIFIED {
			continue
		}
		output = append(output, fmt.Sprintf("%s=%s", mod.Name(), severity))
	}
	return strings.Join(output, ";")
}

// ParseConfigurationString parses a logger configuration string into a map of
// logger names and their associated log level. This method is provided to
// allow other programs to pre-validate a configuration string rather than
// just calling ConfigureLoggers.
//
// Loggers are colon- or semicolon-separated; each module is specified as
// <modulename>=<level>.  White space outside of module names and levels is
// ignored.  The root module is specified with the name "<root>".
//
// As a special case, a log level may be specified on its own.
// This is equivalent to specifying the level of the root module,
// so "DEBUG" is equivalent to `<root>=DEBUG`
//
// An example specification:
//	`<root>=ERROR; foo.bar=WARNING`
func ParseConfigurationString(specification string) (map[string]Level, error) {
	levels := make(map[string]Level)
	if level, ok := ParseLevel(specification); ok {
		levels[""] = level
		return levels, nil
	}
	values := strings.FieldsFunc(specification, func(r rune) bool { return r == ';' || r == ':' })
	for _, value := range values {
		s := strings.SplitN(value, "=", 2)
		if len(s) < 2 {
			return nil, fmt.Errorf("logger specification expected '=', found %q", value)
		}
		name := strings.TrimSpace(s[0])
		levelStr := strings.TrimSpace(s[1])
		if name == "" || levelStr == "" {
			return nil, fmt.Errorf("logger specification %q has blank name or level", value)
		}
		if name == "<root>" {
			name = ""
		}
		level, ok := ParseLevel(levelStr)
		if !ok {
			return nil, fmt.Errorf("unknown severity level %q", levelStr)
		}
		levels[name] = level
	}
	return levels, nil
}

type loggerGetter interface {
	Get(name string) Logger
}

func configureLoggers(specification string, loggers loggerGetter) error {
	if specification == "" {
		return nil
	}
	levels, err := ParseConfigurationString(specification)
	if err != nil {
		return err
	}
	for name, level := range levels {
		loggers.Get(name).SetLogLevel(level)
	}
	return nil
}
