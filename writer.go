// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package loggo

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// defaultWriterName is the name of the writer default writer.
const defaultWriterName = "default"

// Writer is implemented by any recipient of log messages.
type Writer interface {
	// Write writes a message to the Writer with the given
	// level and module name. The filename and line hold
	// the file name and line number of the code that is
	// generating the log message; the time stamp holds
	// the time the log message was generated, and
	// message holds the log message itself.
	Write(level Level, name, filename string, line int, timestamp time.Time, message string)
}

type simpleWriter struct {
	writer    io.Writer
	formatter Formatter
}

// NewSimpleWriter returns a new writer that writes
// log messages to the given io.Writer formatting the
// messages with the given formatter.
func NewSimpleWriter(writer io.Writer, formatter Formatter) Writer {
	return &simpleWriter{writer, formatter}
}

func (simple *simpleWriter) Write(level Level, module, filename string, line int, timestamp time.Time, message string) {
	logLine := simple.formatter.Format(level, module, filename, line, timestamp, message)
	fmt.Fprintln(simple.writer, logLine)
}

type minLevelWriter struct {
	writer Writer
	level  Level
}

// Writers holds a set of Writers and provides operations for
// acting on that set. It also acts as a single Writer.
type Writers struct {
	mu               sync.Mutex
	combinedMinLevel Level
	all              map[string]*minLevelWriter
	init             func() map[string]*minLevelWriter
}

// NewWriters creates a new set of Writers using the provided
// details.
func NewWriters(initial map[string]*minLevelWriter) *Writers {
	ws := &Writers{}
	ws.reset(initial)
	return ws
}

// reset puts the list of Writers back into the initial state.
func (ws *Writers) reset(initial map[string]*minLevelWriter) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	ws.all = make(map[string]*minLevelWriter)
	for name, writer := range initial {
		ws.addUnlocked(name, writer)
	}
	ws.resetMinLevel()
}

// AddWithLevel adds the writer to the list of Writers that get notified
// when Write() is called. When adding, the caller specifies the
// minimum logging level that will be written and a name for the
// writer. The name is used to identify the writer later (e.g. when
// removing it).
//
// If there is already a writer with that name, an error is returned.
func (ws *Writers) AddWithLevel(name string, writer Writer, minLevel Level) error {
	if writer == nil {
		return fmt.Errorf("Writer cannot be nil")
	}
	return ws.add(name, &minLevelWriter{
		writer: writer,
		level:  minLevel,
	})
}

// add adds the writer to the list of Writers that get notified when
// Write() is called. When adding, the caller specifies the name for
// the writer. The name is used to identify the writer later (e.g. when
// removing it).
//
// If there is already a writer with that name, an error is returned.
func (ws *Writers) add(name string, writer *minLevelWriter) error {
	if writer == nil {
		return fmt.Errorf("Writer cannot be nil")
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if err := ws.addUnlocked(name, writer); err != nil {
		return err
	}

	ws.resetMinLevel()
	return nil
}

func (ws *Writers) addUnlocked(name string, writer *minLevelWriter) error {
	if _, found := ws.all[name]; found {
		return fmt.Errorf("there is already a Writer with the name %q", name)
	}
	ws.all[name] = writer
	return nil
}

// remove drops the Writer identified by 'name' from the set and
// returns that writer. If the Writer is not found, an error is
// returned.
func (ws *Writers) remove(name string) (*minLevelWriter, error) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	writer, err := ws.removeUnlocked(name)
	if err != nil {
		return nil, err
	}

	ws.resetMinLevel()
	return writer, nil
}

func (ws *Writers) removeUnlocked(name string) (*minLevelWriter, error) {
	writer, found := ws.all[name]
	if !found {
		return nil, fmt.Errorf("Writer %q is not recognized", name)
	}
	delete(ws.all, name)
	return writer, nil
}

// replace is a convenience method that does the atomic equivalent of
// calling remove() and then add(). The previous writer, which
// must exist, is returned.
func (ws *Writers) replace(name string, newWriter Writer) (Writer, error) {
	if newWriter == nil {
		return nil, fmt.Errorf("Writer cannot be nil")
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()

	mlw, err := ws.removeUnlocked(name)
	if err != nil {
		return nil, err
	}
	oldWriter := mlw.writer
	mlw.writer = newWriter // keep the level
	if err := ws.addUnlocked(name, mlw); err != nil {
		return nil, err
	}
	ws.resetMinLevel()
	return oldWriter, nil
}

// Write implements Writer, sending the message to each known writer.
func (ws *Writers) Write(level Level, module, filename string, line int, timestamp time.Time, message string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	for _, mlw := range ws.all {
		if level >= mlw.level {
			mlw.writer.Write(level, module, filename, line, timestamp, message)
		}
	}
}

// WillWrite returns whether there are any Writers
// at or above the given severity level. If it returns
// false, any log message at the given level will be discarded.
func (ws *Writers) WillWrite(level Level) bool {
	return level >= ws.combinedMinLevel.get()
}

func (ws *Writers) resetMinLevel() {
	// We assume the lock is already held
	minLevel := CRITICAL
	for _, writer := range ws.all {
		if writer.level < minLevel {
			minLevel = writer.level
		}
	}
	ws.combinedMinLevel.set(minLevel)
}
