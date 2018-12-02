// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"fmt"
	"strings"
	"sync/atomic"
)

// The severity levels. Higher values are more considered more
// important.
const (
	UNSPECIFIED Level = 0000
	EMERGENCY   Level = 0001 // Syslog level 0
	ALERT       Level = 1000 // Syslog level 1
	CRITICAL    Level = 2000 // Syslog level 2
	ERROR       Level = 3000 // Syslog level 3
	WARNING     Level = 4000 // Syslog level 4
	NOTICE      Level = 5000 // Syslog level 5
	INFO        Level = 6000 // Syslog level 6
	DEBUG       Level = 7000 // Syslog level 7
	TRACE       Level = 8000
)

var Levels = map[Level]LevelType{
	UNSPECIFIED: {"", UNSPECIFIED, "UNSPECIFIED"},
	TRACE:       {"TRACE", TRACE, "TRACE"},
	DEBUG:       {"DEBUG", DEBUG, "DEBUG"},         // Syslog 7
	INFO:        {"INFO", INFO, "INFO"},            // Syslog 6
	NOTICE:      {"NOTE", NOTICE, "NOTICE"},        // Syslog 5
	WARNING:     {"WARN", WARNING, "WARNING"},      // Syslog 4
	ERROR:       {"ERROR", ERROR, "ERROR"},         // Syslog 3
	CRITICAL:    {"CRIT", CRITICAL, "CRITICAL"},    // Syslog 2
	ALERT:       {"ALERT", ALERT, "ALERT"},         // Syslog 1
	EMERGENCY:   {"EMERG", EMERGENCY, "EMERGENCY"}, // Syslog 0
}

// Level holds a severity level.
type Level uint32

type LevelType struct {
	ShortName string
	Value     Level
	LongName  string
}

// ParseLevel converts a string representation of a logging level to a
// Level. It returns the level and whether it was valid or not.
func ParseLevel(level string) (Level, bool) {
	if len(level) == 0 {
		return UNSPECIFIED, false
	}

	level = strings.ToUpper(level)
	for _, l := range Levels {
		if level == l.LongName || level == l.ShortName {
			return l.Value, true
		}
	}
	return UNSPECIFIED, false
}

// String implements Stringer.
func (level Level) String() string {
	if l, ok := Levels[level]; ok {
		return l.LongName
	} else {
		return "UNKNOWN"
	}
}

// Short returns a five character string to use in
// aligned logging output.
func (level Level) Short() string {
	if l, ok := Levels[level]; ok {
		return fmt.Sprintf("%5s", l.ShortName)
	} else {
		return "UNKNOWN"
	}
	return "UNKN "

}

// get atomically gets the value of the given level.
func (level *Level) get() Level {
	return Level(atomic.LoadUint32((*uint32)(level)))
}

// set atomically sets the value of the receiver
// to the given level.
func (level *Level) set(newLevel Level) {
	atomic.StoreUint32((*uint32)(level), uint32(newLevel))
}
