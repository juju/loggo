// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/juju/loggo/attrs"
)

// DefaultFormatter returns the parameters separated by spaces except for
// filename and line which are separated by a colon.  The timestamp is shown
// to second resolution in UTC. For example:
//
//	2016-07-02 15:04:05
func DefaultFormatter(entry Entry) string {
	ts := entry.Timestamp.In(time.UTC).Format("2006-01-02 15:04:05")
	// Just get the basename from the filename
	filename := filepath.Base(entry.Filename)

	var (
		format string
		values []any
	)
	for _, attr := range entry.Attrs {
		switch a := attr.(type) {
		case attrs.AttrValue[string]:
			format += " %s=%s"
			values = append(values, a.Key(), a.Value())
		case attrs.AttrValue[int]:
			format += " %s=%d"
			values = append(values, a.Key(), a.Value())
		case attrs.AttrValue[int64]:
			format += " %s=%d"
			values = append(values, a.Key(), a.Value())
		case attrs.AttrValue[uint64]:
			format += " %s=%d"
			values = append(values, a.Key(), a.Value())
		case attrs.AttrValue[float64]:
			format += " %s=%f"
			values = append(values, a.Key(), a.Value())
		case attrs.AttrValue[bool]:
			format += " %s=%t"
			values = append(values, a.Key(), a.Value())
		case attrs.AttrValue[time.Time]:
			format += " %s=%v"
			values = append(values, a.Key(), a.Value())
		case attrs.AttrValue[time.Duration]:
			format += " %s=%v"
			values = append(values, a.Key(), a.Value())
		case attrs.AttrValue[any]:
			format += " %s=%v"
			values = append(values, a.Key(), a.Value())
		}
	}

	args := []any{ts, entry.Level, entry.Module, filename, entry.Line, entry.Message}
	args = append(args, values...)

	return fmt.Sprintf("%s %s %s %s:%d %s"+format, args...)
}

// TimeFormat is the time format used for the default writer.
// This can be set with the environment variable LOGGO_TIME_FORMAT.
var TimeFormat = initTimeFormat()

func initTimeFormat() string {
	format := os.Getenv("LOGGO_TIME_FORMAT")
	if format != "" {
		return format
	}
	return "15:04:05"
}
