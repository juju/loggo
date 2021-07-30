// Copyright 2021 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggoemoji

import (
	"bytes"
	"fmt"

	"github.com/juju/loggo"
	gc "gopkg.in/check.v1"
)

type WriterSuite struct{}

var _ = gc.Suite(&WriterSuite{})

func (s *WriterSuite) TestWriteEmoji(c *gc.C) {
	tests := []struct {
		Level    loggo.Level
		Expected string
	}{{
		Level:    loggo.TRACE,
		Expected: "✏️",
	}, {
		Level:    loggo.DEBUG,
		Expected: "🐞",
	}, {
		Level:    loggo.INFO,
		Expected: "🧐",
	}, {
		Level:    loggo.WARNING,
		Expected: "⚠️ ",
	}, {
		Level:    loggo.ERROR,
		Expected: "😱",
	}, {
		Level:    loggo.CRITICAL,
		Expected: "💥",
	}}
	for i, test := range tests {
		c.Logf("test %d level %s", i, test.Level.Short())

		buf := new(bytes.Buffer)
		writer := NewWriter(buf)
		writer.Write(loggo.Entry{
			Level:   test.Level,
			Message: "Hello",
		})

		c.Assert(buf.String(), gc.Equals, fmt.Sprintf("00:00:00 %s %s  .:0 Hello\n", test.Level.Short(), test.Expected))
	}
}
