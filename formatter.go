// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

// Formatter defines the single method Format, which takes the logging
// record and converts it to a string.
type Formatter interface {
	Format(Record) string
}
