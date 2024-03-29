package main

import (
	"github.com/juju/loggo/v2"
)

var first = loggo.GetLogger("first")

func FirstCritical(message string) {
	first.Criticalf(message)
}

func FirstError(message string) {
	first.Errorf(message)
}

func FirstWarning(message string) {
	first.Warningf(message)
}

func FirstInfo(message string) {
	first.Infof(message)
}

func FirstDebug(message string) {
	first.Debugf(message)
}

func FirstTrace(message string) {
	first.Tracef(message)
}
