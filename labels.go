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

// mergeLabels merges multiple sets of labels into a single set.
// Later sets of labels take precedence over earlier sets.
func mergeLabels(labels []Labels) Labels {
	result := make(Labels)
	for _, l := range labels {
		for k, v := range l {
			result[k] = v
		}
	}
	return result
}
