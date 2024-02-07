// Copyright 2024 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

const (
	// LoggerTags is the name of the label used to record the
	// logger tags for a log entry.
	LoggerTags = "logger-tags"
)

// Labels represents key values which are assigned to a log entry.
type Labels map[string]string
