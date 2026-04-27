// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"testing"
	"time"

	"github.com/juju/loggo/v3"
	"github.com/juju/tc"
)

type formatterSuite struct{}

func TestFormatterSuite(t *testing.T) {
	tc.Run(t, &formatterSuite{})
}

func (*formatterSuite) TestDefaultFormat(c *tc.C) {
	location, err := time.LoadLocation("UTC")
	testTime := time.Date(2013, 5, 3, 10, 53, 24, 123456, location)
	c.Assert(err, tc.IsNil)
	entry := loggo.Entry{
		Level:     loggo.WARNING,
		Module:    "test.module",
		Filename:  "some/deep/filename",
		Line:      42,
		Timestamp: testTime,
		Message:   "hello world!",
	}
	formatted := loggo.DefaultFormatter(entry)
	c.Assert(formatted, tc.Equals, "2013-05-03 10:53:24 WARNING test.module filename:42 hello world!")
}
