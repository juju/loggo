package main

import (
	"github.com/juju/loggo/v2"
)

var second = loggo.GetLogger("second")

func SecondCritical(message string) {
	second.Criticalf(message)
}

func SecondError(message string) {
	second.Errorf(message)
}

func SecondWarning(message string) {
	second.Warningf(message)
}

func SecondInfo(message string) {
	second.Infof(message)
}

func SecondDebug(message string) {
	second.Debugf(message)
}
func SecondTrace(message string) {
	second.Tracef(message)
}
