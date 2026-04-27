// Copyright 2024 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package slog

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/juju/loggo/v2"
	"github.com/juju/loggo/v2/attrs"
)

// mockHandler captures records written via Handle for inspection.
type mockHandler struct {
	records []slog.Record
}

func (h *mockHandler) Enabled(_ context.Context, _ slog.Level) bool { return true }
func (h *mockHandler) Handle(_ context.Context, r slog.Record) error {
	h.records = append(h.records, r)
	return nil
}
func (h *mockHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *mockHandler) WithGroup(_ string) slog.Handler      { return h }

func TestNewSlogWriter(t *testing.T) {
	handler := &mockHandler{}
	w := NewSlogWriter(handler)
	if w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestWriteBasicEntry(t *testing.T) {
	handler := &mockHandler{}
	w := NewSlogWriter(handler)

	now := time.Now()
	entry := loggo.Entry{
		Level:     loggo.INFO,
		Module:    "test.module",
		Filename:  "/path/to/file.go",
		Line:      42,
		Timestamp: now,
		Message:   "hello world",
	}

	_ = w.Write(context.Background(), entry)

	if len(handler.records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(handler.records))
	}

	rec := handler.records[0]
	if rec.Message != "hello world" {
		t.Errorf("expected message %q, got %q", "hello world", rec.Message)
	}
	if rec.Level != slog.LevelInfo {
		t.Errorf("expected level %v, got %v", slog.LevelInfo, rec.Level)
	}
	if !rec.Time.Equal(now) {
		t.Errorf("expected time %v, got %v", now, rec.Time)
	}

	// Check default attrs (module, filename, line).
	attrMap := recordAttrs(rec)
	if v := attrMap["module"]; v != "test.module" {
		t.Errorf("expected module %q, got %q", "test.module", v)
	}
	if v := attrMap["filename"]; v != "/path/to/file.go" {
		t.Errorf("expected filename %q, got %q", "/path/to/file.go", v)
	}
	if v := attrMap["line"]; v != int64(42) {
		t.Errorf("expected line %d, got %v", 42, v)
	}
}

func TestWriteWithLabels(t *testing.T) {
	handler := &mockHandler{}
	w := NewSlogWriter(handler)

	entry := loggo.Entry{
		Level:     loggo.DEBUG,
		Module:    "test",
		Filename:  "file.go",
		Line:      1,
		Timestamp: time.Now(),
		Message:   "labeled",
		Labels:    loggo.Labels{"env": "prod", "version": "1.0"},
	}

	_ = w.Write(context.Background(), entry)

	if len(handler.records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(handler.records))
	}

	attrMap := recordAttrs(handler.records[0])
	if v := attrMap["env"]; v != "prod" {
		t.Errorf("expected env %q, got %v", "prod", v)
	}
	if v := attrMap["version"]; v != "1.0" {
		t.Errorf("expected version %q, got %v", "1.0", v)
	}
}

func TestWriteWithAllAttrTypes(t *testing.T) {
	handler := &mockHandler{}
	w := NewSlogWriter(handler)

	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	dur := 5 * time.Second

	entry := loggo.Entry{
		Level:     loggo.WARNING,
		Module:    "test",
		Filename:  "file.go",
		Line:      1,
		Timestamp: time.Now(),
		Message:   "attrs test",
		Attrs: []any{
			attrs.String("s", "hello"),
			attrs.Int("i", 99),
			attrs.Int64("i64", 1234567890123),
			attrs.Uint64("u64", 9876543210),
			attrs.Float64("f64", 2.718),
			attrs.Bool("b", true),
			attrs.Time("ts", now),
			attrs.Duration("dur", dur),
			attrs.Any("obj", []int{1, 2, 3}),
		},
	}

	_ = w.Write(context.Background(), entry)

	if len(handler.records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(handler.records))
	}

	attrMap := recordAttrs(handler.records[0])
	if v := attrMap["s"]; v != "hello" {
		t.Errorf("expected s=%q, got %v", "hello", v)
	}
	if v := attrMap["i"]; v != int64(99) {
		t.Errorf("expected i=%d, got %v", 99, v)
	}
	if v := attrMap["i64"]; v != int64(1234567890123) {
		t.Errorf("expected i64=%d, got %v", int64(1234567890123), v)
	}
	if v := attrMap["u64"]; v != uint64(9876543210) {
		t.Errorf("expected u64=%d, got %v", uint64(9876543210), v)
	}
	if v := attrMap["f64"]; v != 2.718 {
		t.Errorf("expected f64=%f, got %v", 2.718, v)
	}
	if v := attrMap["b"]; v != true {
		t.Errorf("expected b=%t, got %v", true, v)
	}
	if v, ok := attrMap["ts"].(time.Time); !ok || !v.Equal(now) {
		t.Errorf("expected ts=%v, got %v", now, attrMap["ts"])
	}
	if v, ok := attrMap["dur"].(time.Duration); !ok || v != dur {
		t.Errorf("expected dur=%v, got %v", dur, attrMap["dur"])
	}
	if v := attrMap["obj"]; v == nil {
		t.Error("expected obj to be non-nil")
	}
}

func TestLevelMapping(t *testing.T) {
	tests := []struct {
		input    loggo.Level
		expected slog.Level
	}{
		{input: loggo.TRACE, expected: slog.LevelDebug - 1},
		{input: loggo.DEBUG, expected: slog.LevelDebug},
		{input: loggo.INFO, expected: slog.LevelInfo},
		{input: loggo.WARNING, expected: slog.LevelWarn},
		{input: loggo.ERROR, expected: slog.LevelError},
		{input: loggo.CRITICAL, expected: slog.LevelError + 1},
	}

	for _, tc := range tests {
		t.Run(tc.input.String(), func(t *testing.T) {
			got := Level(tc.input)
			if got != tc.expected {
				t.Errorf("Level(%v) = %v, want %v", tc.input, got, tc.expected)
			}
		})
	}
}

func TestLevelPanicsOnUnknown(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for unknown level")
		}
	}()
	Level(loggo.UNSPECIFIED)
}

func TestDefaultLevelSingleEntry(t *testing.T) {
	config := loggo.Config{
		"root": loggo.WARNING,
	}
	got := DefaultLevel(config)
	if got != Level(loggo.WARNING) {
		t.Errorf("expected %v, got %v", Level(loggo.WARNING), got)
	}
}

func TestDefaultLevelMultipleEntries(t *testing.T) {
	config := loggo.Config{
		"root":   loggo.ERROR,
		"module": loggo.DEBUG,
		"other":  loggo.WARNING,
	}
	got := DefaultLevel(config)
	// DEBUG is the lowest level configured.
	if got != Level(loggo.DEBUG) {
		t.Errorf("expected %v, got %v", Level(loggo.DEBUG), got)
	}
}

func TestDefaultLevelAllSame(t *testing.T) {
	config := loggo.Config{
		"a": loggo.INFO,
		"b": loggo.INFO,
		"c": loggo.INFO,
	}
	got := DefaultLevel(config)
	if got != Level(loggo.INFO) {
		t.Errorf("expected %v, got %v", Level(loggo.INFO), got)
	}
}

func TestDefaultLevelWithTrace(t *testing.T) {
	config := loggo.Config{
		"root":  loggo.ERROR,
		"debug": loggo.TRACE,
	}
	got := DefaultLevel(config)
	// TRACE is the lowest possible level.
	if got != Level(loggo.TRACE) {
		t.Errorf("expected %v, got %v", Level(loggo.TRACE), got)
	}
}

func TestWriteWithNoAttrsOrLabels(t *testing.T) {
	handler := &mockHandler{}
	w := NewSlogWriter(handler)

	entry := loggo.Entry{
		Level:     loggo.ERROR,
		Module:    "minimal",
		Filename:  "min.go",
		Line:      10,
		Timestamp: time.Now(),
		Message:   "simple error",
	}

	_ = w.Write(context.Background(), entry)

	if len(handler.records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(handler.records))
	}

	rec := handler.records[0]
	if rec.Level != slog.LevelError {
		t.Errorf("expected level %v, got %v", slog.LevelError, rec.Level)
	}
}

// recordAttrs extracts all attributes from a slog.Record into a map.
func recordAttrs(r slog.Record) map[string]any {
	result := make(map[string]any)
	r.Attrs(func(a slog.Attr) bool {
		result[a.Key] = a.Value.Any()
		return true
	})
	return result
}
