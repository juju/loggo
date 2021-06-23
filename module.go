// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

// Do not change rootName: modules.resolve() will misbehave if it isn't "".
const (
	rootString       = "<root>"
	defaultRootLevel = WARNING
	defaultLevel     = UNSPECIFIED
)

type module interface {
	Name() string
	FullName() string
	WillWrite(Level) bool
	Level() *Level
	GetEffectiveLogLevel() Level
	SetLevel(Level)
	Write(Entry)
	SetParent(module)
	Parent() module
	Context() *Context
}

type single struct {
	name    string
	level   *Level
	parent  module
	context *Context
}

// Name returns the module's name.
func (m *single) Name() string {
	return m.name
}

func (m *single) FullName() string {
	if m.name == "" {
		return rootString
	}
	return m.name
}

func (m *single) WillWrite(level Level) bool {
	if level < TRACE || level > CRITICAL {
		return false
	}
	return level >= m.GetEffectiveLogLevel()
}

func (m *single) GetEffectiveLogLevel() Level {
	// Note: the root module is guaranteed to have a
	// specified logging level, so acts as a suitable sentinel
	// for this loop.
	var p module = m
	for {
		if level := p.Level().get(); level != UNSPECIFIED {
			return level
		}
		p = p.Parent()
	}
}

func (m *single) Level() *Level {
	return m.level
}

// SetLevel sets the severity level of the given module.
// The root module cannot be set to UNSPECIFIED level.
func (m *single) SetLevel(level Level) {
	// The root module can't be unspecified.
	if m.name == "" && level == UNSPECIFIED {
		level = WARNING
	}
	m.level.set(level)
}

func (m *single) Write(entry Entry) {
	entry.Module = m.name
	m.context.write(entry)
}

func (m *single) Parent() module {
	return m.parent
}

func (m *single) SetParent(p module) {
	m.parent = p
}

func (m *single) Context() *Context {
	return m.context
}

type multiple struct {
	name    string
	modules []module
}

func (m *multiple) Name() string {
	return m.name
}

func (m *multiple) FullName() string {
	if m.name == "" {
		return rootString
	}
	return m.name
}

func (m *multiple) WillWrite(level Level) bool {
	for _, module := range m.modules {
		if !module.WillWrite(level) {
			return false
		}
	}
	return true
}

func (m *multiple) Level() *Level {
	level := m.GetEffectiveLogLevel()
	return &level
}

func (m *multiple) GetEffectiveLogLevel() Level {
	// Get the highest level...
	level := UNSPECIFIED
	for _, module := range m.modules {
		if m := *module.Level(); m > level {
			level = m
		}
	}
	if level == UNSPECIFIED {
		return WARNING
	}
	return level
}

func (m *multiple) SetLevel(level Level) {
	for _, module := range m.modules {
		module.SetLevel(level)
	}
}

func (m *multiple) Write(entry Entry) {
	for _, module := range m.modules {
		module.Write(entry)
	}
}

func (m *multiple) Parent() module {
	return nil
}

func (m *multiple) SetParent(p module) {
	// no-op
}

func (m *multiple) Context() *Context {
	return nil
}
