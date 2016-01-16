// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"time"

	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
)

type colorFormatterSuite struct{}

var _ = gc.Suite(&colorFormatterSuite{})

func (*colorFormatterSuite) TestColorFormat(c *gc.C) {
	location, err := time.LoadLocation("UTC")
	c.Assert(err, gc.IsNil)
	testTime := time.Date(2013, 5, 3, 10, 53, 24, 123456, location)
	colorFormatter := &loggo.ColorFormatter{}
	colorFormatted := colorFormatter.Format(loggo.WARNING, "test.module", "some/deep/filename", 42, testTime, "hello world!")
	c.Assert(colorFormatted, gc.Equals, "10:53:24 \x1b[33mwarn \x1b[0m test.module filename:42 hello world!")
}
