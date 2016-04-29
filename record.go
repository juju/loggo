// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"fmt"
	"runtime"
	"time"
)

// Record holds the information for a single log record.
type Record struct {
	// Level is the log level of the log record.
	Level Level

	// LoggerName is the name of the logger for which the record
	// was generated.
	LoggerName string

	// Filename is the path to the code that triggered the record.
	Filename string

	// Line is the line number of the code that triggered the record.
	Line int

	// Timestamp is the time that the record was created.
	Timestamp time.Time

	// Message is the requested log message.
	Message string
}

// NewRecord creates a new log record for the given log level, logger name,
// and message. It uses the identified entry the call stack to determine
// the filename and line number. The current time is used for the
// timestamp.
func NewRecord(calldepth int, level Level, loggerName, message string) Record {
	// Gather time, filename, and line number.
	now := time.Now() // get this early.
	// Param to Caller is the call depth.  Since this method is called from
	// the Logger methods, we want the place that those were called from.
	_, file, line, ok := runtime.Caller(calldepth + 1)
	if !ok {
		file = "???"
		line = 0
	}

	// Trim newline off format string, following usual
	// Go logging conventions.
	if len(message) > 0 && message[len(message)-1] == '\n' {
		message = message[0 : len(message)-1]
	}

	return Record{
		Level:      level,
		LoggerName: loggerName,
		Filename:   file,
		Line:       line,
		Timestamp:  now,
		Message:    message,
	}
}

// NewRecordf creates a new log record for the given info. The only
// difference from NewRecord() is that the provided args are applied// to the message using fmt.Sprintf().
func NewRecordf(calldepth int, level Level, loggerName, message string, args ...interface{}) Record {
	rec := NewRecord(calldepth+1, level, loggerName, message)
	if len(args) == 0 {
		return rec
	}

	// To avoid having a proliferation of Info/Infof methods,
	// only use Sprintf if there are any args, and rely on the
	// `go vet` tool for the obvious cases where someone has forgotten
	// to provide an arg.
	rec.Message = fmt.Sprintf(rec.Message, args...)

	return rec
}
