// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

// TODO(ericsnow) Drop these helpers as soon as the actual functions exist.
func NewLogger(name string, writer Writer) Logger {
	if writer != nil {
		ReplaceDefaultWriter(writer)
	}
	return GetLogger(name)
}

func NewLoggerWithParent(name string, parent Logger, writer Writer) Logger {
	// We assume the parent is already added and can ignore it.
	return NewLogger(name, writer)
}
