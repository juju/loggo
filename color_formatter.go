// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

// ColorFormatter is a compact formatter with colored output, ideal for output expected to be
// viewed in a terminal
type ColorFormatter struct{}

var levelStrs = map[Level]string{
	TRACE:    color.New(color.FgWhite).SprintFunc()("trace"),
	DEBUG:    color.New(color.FgGreen).SprintFunc()("debug"),
	INFO:     color.New(color.FgBlue).SprintFunc()("info "),
	WARNING:  color.New(color.FgYellow).SprintFunc()("warn "),
	ERROR:    color.New(color.FgRed).SprintFunc()("error"),
	CRITICAL: color.New(color.BgRed).SprintFunc()("critc"),
}

// Format returns the parameters separated by spaces except for filename and
// line which are separated by a colon.  Only the time is shown to second resolution
// to make the output compact.
func (*ColorFormatter) Format(level Level, module string, filename string, line int, timestamp time.Time, message string) string {
	ts := timestamp.In(time.UTC).Format("15:04:05")
	filename = filepath.Base(filename)
	return fmt.Sprintf("%s %s %s %s:%d %s", ts, levelStrs[level], module, filename, line, message)
}
