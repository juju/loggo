// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"fmt"
	"time"

	"github.com/juju/tc"
)

func Between(start, end time.Time) tc.Checker {
	if end.Before(start) {
		return &betweenChecker{start: end, end: start}
	}
	return &betweenChecker{start: start, end: end}
}

type betweenChecker struct {
	start, end time.Time
}

func (checker *betweenChecker) Info() *tc.CheckerInfo {
	info := tc.CheckerInfo{
		Name:   "Between",
		Params: []string{"obtained"},
	}
	return &info
}

func (checker *betweenChecker) Check(params []interface{}, names []string) (result bool, error string) {
	when, ok := params[0].(time.Time)
	if !ok {
		return false, "obtained value type must be time.Time"
	}
	if when.Before(checker.start) {
		return false, fmt.Sprintf("obtained time %q is before start time %q", when, checker.start)
	}
	if when.After(checker.end) {
		return false, fmt.Sprintf("obtained time %q is after end time %q", when, checker.end)
	}
	return true, ""
}
