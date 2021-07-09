package loggostructured

import (
	"encoding/json"
	"io"
	"path/filepath"

	"github.com/juju/loggo"
)

type structuredWriter struct {
	encoder    *json.Encoder
	timeFormat string
}

// NewWriter will write out structured logs to the writer.
// We fallback to JSON Codec for encoding all loggo.Entries.
func NewWriter(writer io.Writer, timeFormat string) loggo.Writer {
	return &structuredWriter{
		encoder:    json.NewEncoder(writer),
		timeFormat: timeFormat,
	}
}

// Write implements Writer.
func (w *structuredWriter) Write(entry loggo.Entry) {
	ts := entry.Timestamp.Format(w.timeFormat)
	// Just get the basename from the filename
	filename := filepath.Base(entry.Filename)

	_ = w.encoder.Encode(line{
		Timestamp: ts,
		Filename:  filename,
		Line:      entry.Line,
		Level:     entry.Level.Short(),
		Module:    entry.Module,
		Message:   entry.Message,
		Labels:    entry.Labels,
	})
}

type line struct {
	Timestamp string   `json:"ts"`
	Filename  string   `json:"path"`
	Line      int      `json:"line"`
	Level     string   `json:"level"`
	Module    string   `json:"module"`
	Message   string   `json:"message"`
	Labels    []string `json:"labels,omitempty"`
}
