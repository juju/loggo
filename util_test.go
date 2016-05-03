// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	gc "gopkg.in/check.v1"

	"github.com/juju/loggo"
)

func init() {
	setLocationsForTags("logger_test.go")
	setLocationsForTags("writer_test.go")
}

func newTraceLogger(name string) (loggo.Logger, *loggo.TestWriter) {
	writer := &loggo.TestWriter{}
	loggo.ReplaceDefaultWriter(writer)
	logger := loggo.GetLogger(name)
	// Make it so the logger itself writes all messages.
	logger.SetLogLevel(loggo.TRACE)
	return logger, writer
}

func checkLastMessage(c *gc.C, writer *loggo.TestWriter, expected string) {
	log := writer.Log()
	writer.Clear()
	obtained := log[len(log)-1].Message
	c.Check(obtained, gc.Equals, expected)
}

func setupTempFileWriter(c *gc.C) (logFile *os.File, cleanup func()) {
	loggo.RemoveWriter("test")
	logFile, err := ioutil.TempFile("", "loggo-test")
	c.Assert(err, gc.IsNil)
	cleanup = func() {
		logFile.Close()
		os.Remove(logFile.Name())
	}
	writer := loggo.NewSimpleWriter(logFile, &loggo.DefaultFormatter{})
	err = loggo.RegisterWriter("testfile", writer, loggo.TRACE)
	c.Assert(err, gc.IsNil)
	return
}

func assertLocation(c *gc.C, msg loggo.TestLogValues, tag string) {
	loc := location(tag)
	c.Assert(msg.Filename, gc.Equals, loc.file)
	c.Assert(msg.Line, gc.Equals, loc.line)
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
				panic(fmt.Errorf("tag %q already processed previously"))
			}
			tagToLocation[tag] = Location{file: filename, line: i + 1}
		}
	}
}

func Between(start, end time.Time) gc.Checker {
	if end.Before(start) {
		return &betweenChecker{end, start}
	}
	return &betweenChecker{start, end}
}

type betweenChecker struct {
	start, end time.Time
}

func (checker *betweenChecker) Info() *gc.CheckerInfo {
	info := gc.CheckerInfo{
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
