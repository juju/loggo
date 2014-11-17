// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

func assertLocation(c *gc.C, msg loggo.TestLogValues, tag string) {
	loc := location(tag)
	c.Assert(msg.Filename, gc.Equals, loc.file)
	c.Assert(msg.Line, gc.Equals, loc.line)
}

func location(tag string) Location {
	loc, ok := tagToLocation[tag]
	if !ok {
		panic(fmt.Errorf("tag %q not found", tag))
	}
	return loc
}

type Location struct {
	file string
	line int
}

func (loc Location) String() string {
	return fmt.Sprintf("%s:%d", loc.file, loc.line)
}

var tagToLocation = make(map[string]Location)

func setLocationsForTags(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if j := strings.Index(line, "//tag "); j >= 0 {
			tag := line[j+len("//tag "):]
			if _, found := tagToLocation[tag]; found {
				panic(fmt.Errorf("tag %q already processed previously"))
			}
			tagToLocation[tag] = Location{file: filename, line: i + 1}
		}
	}
}

func init() {
	setLocationsForTags("logger_test.go")
	setLocationsForTags("writer_test.go")
}
