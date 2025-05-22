// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/juju/loggo/v2"

	gc "gopkg.in/check.v1"
)

func init() {
	setLocationsForTags("logging_test.go")
	setLocationsForTags("writer_test.go")
}

func assertLocation(c *gc.C, msg loggo.Entry, tag string) {
	loc := location(tag)
	c.Check(fmt.Sprintf("%s:%d", msg.Filename, msg.Line), gc.Equals,
		fmt.Sprintf("%s:%d", loc.file, loc.line),
		gc.Commentf("tag=%s", tag))
}

// All this location stuff is to avoid having hard coded line numbers
// in the tests.  Any line where as a test writer you want to capture the
// file and line number, add a comment that has `//tag name` as the end of
// the line.  The name must be unique across all the tests, and the test
// will panic if it is not.  This name is then used to read the actual
// file and line numbers.

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
				panic(fmt.Errorf("tag %q already processed previously", tag))
			}
			tagToLocation[tag] = Location{file: filename, line: i + 1}
		}
	}
}
