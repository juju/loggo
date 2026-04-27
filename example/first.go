package main

import (
	"context"

	"github.com/juju/loggo/v2"
	"github.com/juju/loggo/v2/attrs"
)

var first = loggo.GetLogger("first")

func FirstCritical(message string) {
	_ = first.Criticalf(context.Background(), message, attrs.String("baz", "boo"))
}

func FirstError(message string) {
	_ = first.Errorf(context.Background(), message)
}

func FirstWarning(message string) {
	_ = first.Warningf(context.Background(), message)
}

func FirstInfo(message string) {
	_ = first.Infof(context.Background(), message)
}

func FirstDebug(message string) {
	_ = first.Debugf(context.Background(), message)
}

func FirstTrace(message string) {
	_ = first.Tracef(context.Background(), message)
}
