// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

var root = &module{level: WARNING}

type module struct {
	name   string
	level  Level
	parent *module
}

// Name returns the module's name.
func (module *module) Name() string {
	if module.name == "" {
		return "<root>"
	}
	return module.name
}

func (module *module) getEffectiveLogLevel() Level {
	// Note: the root module is guaranteed to have a
	// specified logging level, so acts as a suitable sentinel
	// for this loop.
	for {
		if level := module.level.get(); level != UNSPECIFIED {
			return level
		}
		module = module.parent
	}
	panic("unreachable")
}
