// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo_test

import (
	"io"
	"os"
	"testing"

	"github.com/juju/loggo/v3"
	"github.com/juju/tc"
)

func BenchmarkLoggingNoWriters(b *testing.B) {
	c := &tc.TBC{TB: b}
	logger, _ := setupTest(c)

	// No writers
	_, _ = loggo.RemoveWriter("test")
	for i := 0; i < b.N; i++ {
		_ = logger.Warningf(b.Context(), "just a simple warning for %d", i)
	}
}

func BenchmarkLoggingNoWritersNoFormat(b *testing.B) {
	c := &tc.TBC{TB: b}
	logger, _ := setupTest(c)

	// No writers
	_, _ = loggo.RemoveWriter("test")
	for i := 0; i < b.N; i++ {
		_ = logger.Warningf(b.Context(), "just a simple warning")
	}
}

func BenchmarkLoggingTestWriters(b *testing.B) {
	c := &tc.TBC{TB: b}
	logger, writer := setupTest(c)

	for i := 0; i < b.N; i++ {
		_ = logger.Warningf(b.Context(), "just a simple warning for %d", i)
	}
	c.Assert(writer.Log(), tc.HasLen, b.N)
}

func BenchmarkLoggingDiskWriter(b *testing.B) {
	c := &tc.TBC{TB: b}
	logger, _ := setupTest(c)

	logFile := setupTempFileWriter(c)
	defer func() { _ = logFile.Close() }()
	msg := "just a simple warning for %d"
	for i := 0; i < b.N; i++ {
		_ = logger.Warningf(b.Context(), msg, i)
	}
	offset, err := logFile.Seek(0, io.SeekCurrent)
	c.Assert(err, tc.IsNil)
	c.Assert((offset > int64(len(msg))*int64(b.N)), tc.Equals, true,
		tc.Commentf("Not enough data was written to the log file."))
}

func BenchmarkLoggingDiskWriterNoMessages(b *testing.B) {
	c := &tc.TBC{TB: b}
	logger, _ := setupTest(c)

	logFile := setupTempFileWriter(c)
	defer func() { _ = logFile.Close() }()
	// Change the log level
	writer, err := loggo.RemoveWriter("testfile")
	c.Assert(err, tc.IsNil)

	err = loggo.RegisterWriter("testfile", loggo.NewMinimumLevelWriter(writer, loggo.WARNING))
	c.Assert(err, tc.IsNil)

	msg := "just a simple warning for %d"
	for i := 0; i < b.N; i++ {
		_ = logger.Debugf(b.Context(), msg, i)
	}
	offset, err := logFile.Seek(0, io.SeekCurrent)
	c.Assert(err, tc.IsNil)
	c.Assert(offset, tc.Equals, int64(0),
		tc.Commentf("Data was written to the log file."))
}

func BenchmarkLoggingDiskWriterNoMessagesLogLevel(b *testing.B) {
	c := &tc.TBC{TB: b}
	logger, _ := setupTest(c)

	logFile := setupTempFileWriter(c)
	defer func() { _ = logFile.Close() }()
	// Change the log level
	logger.SetLogLevel(loggo.WARNING)

	msg := "just a simple warning for %d"
	for i := 0; i < b.N; i++ {
		_ = logger.Debugf(b.Context(), msg, i)
	}
	offset, err := logFile.Seek(0, io.SeekCurrent)
	c.Assert(err, tc.IsNil)
	c.Assert(offset, tc.Equals, int64(0),
		tc.Commentf("Data was written to the log file."))
}

func setupTest(c *tc.TBC) (loggo.Logger, *writer) {
	loggo.ResetLogging()
	logger := loggo.GetLogger("test.writer")
	writer := &writer{}
	err := loggo.RegisterWriter("test", writer)
	c.Assert(err, tc.IsNil)

	return logger, writer
}

func setupTempFileWriter(c *tc.TBC) *os.File {
	_, _ = loggo.RemoveWriter("test")
	logFile, err := os.CreateTemp(c.MkDir(), "loggo-test")
	c.Assert(err, tc.IsNil)

	writer := loggo.NewSimpleWriter(logFile, loggo.DefaultFormatter)
	err = loggo.RegisterWriter("testfile", writer)
	c.Assert(err, tc.IsNil)
	return logFile
}
