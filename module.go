// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

// Do not change rootName: modules.resolve() will misbehave if it isn't "".
const (
	rootString = "<root>"
)

type module struct {
	name    string
	level   Level
	parent  *module
	context *Context

	tags       []string
	tagsLookup map[string]struct{}

	labels Labels
}

// Name returns the module's name.
func (m *module) Name() string {
	if m.name == "" {
		return rootString
	}
	return m.name
}

func (m *module) willWrite(level Level) bool {
	if level < TRACE || level > CRITICAL {
		return false
	}
	return level >= m.getEffectiveLogLevel()
}

func (m *module) getEffectiveLogLevel() Level {
	// Note: the root module is guaranteed to have a
	// specified logging level, so acts as a suitable sentinel
	// for this loop.
	for {
		if level := m.level.get(); level != UNSPECIFIED {
			return level
		}
		m = m.parent
	}
}

// setLevel sets the severity level of the given module.
// The root module cannot be set to UNSPECIFIED level.
func (m *module) setLevel(level Level) {
	// The root module can't be unspecified.
	if m.name == "" && level == UNSPECIFIED {
		level = WARNING
	}
	m.level.set(level)
}

func (m *module) write(entry Entry) {
	entry.Module = m.name
	m.context.write(entry)
}

// hasLabelIntersection returns true if the module has the given labels, which
// may be a subset of the module's labels. The label key and values must match
// and it must match all of the labels in the argument.
func (m *module) hasLabelIntersection(labels Labels) bool {
	for labelKey, labelValue := range labels {
		if moduleValue, ok := m.labels[labelKey]; !ok || moduleValue != labelValue {
			return false
		}
	}
	return true
}
