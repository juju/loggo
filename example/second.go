package main

import (
	"context"

	"github.com/juju/loggo/v3"
)

var second = loggo.GetLogger("second")

func SecondCritical(message string) {
	_ = second.Criticalf(context.Background(), message)
}

func SecondError(message string) {
	_ = second.Errorf(context.Background(), message)
}

func SecondWarning(message string) {
	_ = second.Warningf(context.Background(), message)
}

func SecondInfo(message string) {
	_ = second.Infof(context.Background(), message)
}

func SecondDebug(message string) {
	_ = second.Debugf(context.Background(), message)
}
func SecondTrace(message string) {
	_ = second.Tracef(context.Background(), message)
}
