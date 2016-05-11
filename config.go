// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"fmt"
	"sort"
	"strings"
)

// LoggerConfig holds the configuration for a single logger.
type LoggerConfig struct {
	// Name is the logger's name.
	Name string

	// Level is the log level that should be used by the logger.
	Level Level
}

// ParseLoggerConfig parses a logger configuration string into the
// configuration for a single logger. Whitespace around the spec is
// ignored.
func ParseLoggerConfig(spec string) (LoggerConfig, error) {
	var cfg LoggerConfig

	spec = strings.TrimSpace(spec)
	if spec == "" {
		return cfg, fmt.Errorf("logger config is blank")
	}

	// TODO(ericsnow) Get the name.  We need to sort out backward
	// compability first.

	levelStr := spec // For now level is the only thing in the spec.
	level, ok := ParseLevel(levelStr)
	if !ok {
		return cfg, fmt.Errorf("unknown log level %q", levelStr)
	}
	cfg.Level = level

	return cfg, nil
}

// String returns a logger configuration string that may be parsed
// using ParseLoggerConfig or ParseLoggersConfig().
func (cfg LoggerConfig) String() string {
	// TODO(ericsnow) Include the name.  We need to sort out backward
	// compability first.
	return fmt.Sprintf("%s", cfg.Level)
}

// LoggersConfig is a mapping of logger module names to logger configs.
type LoggersConfig map[string]LoggerConfig

// String returns a logger configuration string that may be parsed
// using ParseLoggersConfig.
func (configs LoggersConfig) String() string {
	// output in alphabetical order.
	names := []string{}
	for name := range configs {
		if name == rootModuleName {
			// This could potentially result in a duplicate entry...
			name = rootName
		}
		names = append(names, name)
	}
	sort.Strings(names)

	var entries []string
	for _, name := range names {
		cfg := configs[name]
		entry := fmt.Sprintf("%s=%s", name, cfg)
		entries = append(entries, entry)
	}
	return strings.Join(entries, ";")
}

// ParseLoggersConfig parses a logger configuration string into a set
// of named logger configs. This method is provided to allow other
// programs to pre-validate a configuration string rather than just
// calling ConfigureLoggers.
//
// Loggers are colon- or semicolon-separated; each module is formatted
// as:
//
//   <modulename>=<config>, where <config> consists of <level>
//
// White space outside of module names and config is ignored. The root
// module is specified with the name "<root>".
//
// As a special case, a config may be specified on its own, without
// a module name. This is equivalent to specifying the configuration
// of the root module, so "DEBUG" is equivalent to `<root>=DEBUG`
//
// An example specification:
//	`<root>=ERROR; foo.bar=WARNING`
func ParseLoggersConfig(spec string) (LoggersConfig, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil, nil
	}

	entries := strings.FieldsFunc(spec, func(r rune) bool { return r == ';' || r == ':' })
	if len(entries) == 1 && !strings.Contains(spec, "=") {
		cfg, err := ParseLoggerConfig(spec)
		if err != nil {
			return nil, err
		}
		return LoggersConfig{rootModuleName: cfg}, nil
	}

	configs := make(LoggersConfig)
	for _, entry := range entries {
		name, cfg, err := parseConfigEntry(entry)
		if err != nil {
			return nil, err
		}
		// last entry wins, for a given name
		configs[name] = cfg
	}
	return configs, nil
}

func parseConfigEntry(entry string) (string, LoggerConfig, error) {
	var cfg LoggerConfig
	pair := strings.SplitN(entry, "=", 2)
	if len(pair) < 2 {
		return "", cfg, fmt.Errorf("logger entry expected '=', found %q", entry)
	}
	name, spec := rootModuleName, entry
	if len(pair) == 2 {
		name, spec = strings.TrimSpace(pair[0]), strings.TrimSpace(pair[1])
		if name == "" {
			return "", cfg, fmt.Errorf("logger entry %q has blank name", entry)
		}
		if name == rootName {
			name = rootModuleName
		}
	}
	if spec == "" {
		return "", cfg, fmt.Errorf("logger entry %q has blank config", entry)
	}
	cfg, err := ParseLoggerConfig(spec)
	return name, cfg, err
}
