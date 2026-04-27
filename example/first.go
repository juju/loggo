package main

import (
	"context"

	"github.com/juju/loggo/v2"
	"github.com/juju/loggo/v2/attrs"
)

var first = loggo.GetLogger("first")

func FirstCritical(message string) {
	first.Criticalf(context.Background(), message, attrs.String("baz", "boo"))
}

func FirstError(message string) {
	first.Errorf(context.Background(), message)
}

func FirstWarning(message string) {
	first.Warningf(context.Background(), message)
}

func FirstInfo(message string) {
	first.Infof(context.Background(), message)
}

func FirstDebug(message string) {
	first.Debugf(context.Background(), message)
}

func FirstTrace(message string) {
	first.Tracef(context.Background(), message)
}
