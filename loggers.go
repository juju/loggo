// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

// Loggers produces loggers for a hierarchy of modules. All the
// loggers will share the same set of log writers.
type Loggers struct {
	m *modules
	w *Writers
}

// NewLoggers returns a new Loggers that uses the provided writers.
// If the root level is UNSPECIFIED, WARNING is used.
func NewLoggers(rootLevel Level, writers *Writers) *Loggers {
	return &Loggers{
		m: newModules(rootLevel),
		w: writers,
	}
}

// LoggersFromConfig creates a new Loggers using the provided writers
// and configures loggers according to the given spec.
func LoggersFromConfig(spec string, writers *Writers) (*Loggers, error) {
	loggers := NewLoggers(UNSPECIFIED, writers)

	configs, err := ParseLoggersConfig(spec)
	if err != nil {
		return nil, err
	}
	loggers.ApplyConfig(configs)

	loggers.m.rootLevel = loggers.Root().MinLogLevel()
	return loggers, nil
}

// Root returns the root logger.
func (ls *Loggers) Root() SubLogger {
	return ls.Get(rootName)
}

// Get returns a logger for the given module name, creating it and
// its parents if necessary.
func (ls *Loggers) Get(name string) SubLogger {
	return newSubLogger(ls.m.get(name), ls.w)
}

// Config returns the current configuration of the Loggers. Loggers
// with UNSPECIFIED level will not be included.
func (ls *Loggers) Config() LoggersConfig {
	return ls.m.config()
}

// ApplyConfig configures the loggers according to the provided configs.
func (ls *Loggers) ApplyConfig(configs LoggersConfig) {
	for name, cfg := range configs {
		logger := ls.Get(name)
		logger.applyConfig(cfg)
	}
}

// resetLevels iterates through the known loggers and sets the levels
// of all to UNSPECIFIED, except for <root> which is set to WARNING.
func (ls *Loggers) resetLevels() {
	ls.m.resetLevels()
}
