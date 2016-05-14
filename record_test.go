// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"time"

	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
)

type RecordSuite struct{}

var _ = gc.Suite(&RecordSuite{})

func (*RecordSuite) TestString(c *gc.C) {
	location, err := time.LoadLocation("UTC")
	c.Assert(err, gc.IsNil)
	testTime := time.Date(2013, 5, 3, 10, 53, 24, 123456, location)
	rec := loggo.Record{
		Level:      loggo.WARNING,
		LoggerName: "test.module",
		Filename:   "some/deep/filename",
		Line:       42,
		Timestamp:  testTime,
		Message:    "hello world!",
	}

	recStr := rec.String()

	c.Check(recStr, gc.Equals, "2013-05-03 10:53:24 WARNING test.module filename:42 hello world!")
}
