// Copyright 2024 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package attrs

import (
	"testing"
	"time"
)

func TestStringAttr(t *testing.T) {
	a := String("key", "value")
	if a.Key() != "key" {
		t.Errorf("expected key %q, got %q", "key", a.Key())
	}
	if a.Value() != "value" {
		t.Errorf("expected value %q, got %q", "value", a.Value())
	}
}

func TestIntAttr(t *testing.T) {
	a := Int("count", 42)
	if a.Key() != "count" {
		t.Errorf("expected key %q, got %q", "count", a.Key())
	}
	if a.Value() != 42 {
		t.Errorf("expected value %d, got %d", 42, a.Value())
	}
}

func TestInt64Attr(t *testing.T) {
	a := Int64("big", 1234567890123)
	if a.Key() != "big" {
		t.Errorf("expected key %q, got %q", "big", a.Key())
	}
	if a.Value() != 1234567890123 {
		t.Errorf("expected value %d, got %d", int64(1234567890123), a.Value())
	}
}

func TestUint64Attr(t *testing.T) {
	a := Uint64("ubig", 9876543210)
	if a.Key() != "ubig" {
		t.Errorf("expected key %q, got %q", "ubig", a.Key())
	}
	if a.Value() != 9876543210 {
		t.Errorf("expected value %d, got %d", uint64(9876543210), a.Value())
	}
}

func TestFloat64Attr(t *testing.T) {
	a := Float64("pi", 3.14159)
	if a.Key() != "pi" {
		t.Errorf("expected key %q, got %q", "pi", a.Key())
	}
	if a.Value() != 3.14159 {
		t.Errorf("expected value %f, got %f", 3.14159, a.Value())
	}
}

func TestBoolAttr(t *testing.T) {
	a := Bool("flag", true)
	if a.Key() != "flag" {
		t.Errorf("expected key %q, got %q", "flag", a.Key())
	}
	if a.Value() != true {
		t.Errorf("expected value %t, got %t", true, a.Value())
	}
}

func TestTimeAttr(t *testing.T) {
	now := time.Now()
	a := Time("ts", now)
	if a.Key() != "ts" {
		t.Errorf("expected key %q, got %q", "ts", a.Key())
	}
	if !a.Value().Equal(now) {
		t.Errorf("expected value %v, got %v", now, a.Value())
	}
}

func TestDurationAttr(t *testing.T) {
	d := 5 * time.Second
	a := Duration("elapsed", d)
	if a.Key() != "elapsed" {
		t.Errorf("expected key %q, got %q", "elapsed", a.Key())
	}
	if a.Value() != d {
		t.Errorf("expected value %v, got %v", d, a.Value())
	}
}

func TestAnyAttr(t *testing.T) {
	type custom struct{ X int }
	val := custom{X: 7}
	a := Any("obj", val)
	if a.Key() != "obj" {
		t.Errorf("expected key %q, got %q", "obj", a.Key())
	}
	if a.Value() != (any)(val) {
		t.Errorf("expected value %v, got %v", val, a.Value())
	}
}

func TestValidWithValidAttrs(t *testing.T) {
	now := time.Now()
	validAttrs := []any{
		String("s", "val"),
		Int("i", 1),
		Int64("i64", 2),
		Uint64("u64", 3),
		Float64("f", 1.0),
		Bool("b", true),
		Time("t", now),
		Duration("d", time.Second),
		Any("a", "anything"),
	}
	if err := Valid(validAttrs); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidWithEmptySlice(t *testing.T) {
	if err := Valid(nil); err != nil {
		t.Errorf("expected no error for nil, got %v", err)
	}
	if err := Valid([]any{}); err != nil {
		t.Errorf("expected no error for empty slice, got %v", err)
	}
}

func TestValidWithInvalidAttr(t *testing.T) {
	invalidAttrs := []any{
		String("s", "val"),
		"not an attr",
	}
	err := Valid(invalidAttrs)
	if err == nil {
		t.Fatal("expected error for invalid attribute, got nil")
	}
	expected := "invalid attribute type string"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestValidWithMultipleInvalidAttrs(t *testing.T) {
	invalidAttrs := []any{
		123,
		String("s", "val"),
	}
	err := Valid(invalidAttrs)
	if err == nil {
		t.Fatal("expected error for invalid attribute, got nil")
	}
	expected := "invalid attribute type int"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestAttrsAllTyped(t *testing.T) {
	now := time.Now()
	input := []any{
		String("s", "val"),
		Int("i", 1),
		Int64("i64", 2),
		Uint64("u64", 3),
		Float64("f", 1.0),
		Bool("b", true),
		Time("t", now),
		Duration("d", time.Second),
		Any("a", "anything"),
	}
	generic, typed := Attrs(input...)
	if len(generic) != 0 {
		t.Errorf("expected no generic attrs, got %d", len(generic))
	}
	if len(typed) != len(input) {
		t.Errorf("expected %d typed attrs, got %d", len(input), len(typed))
	}
}

func TestAttrsAllGeneric(t *testing.T) {
	input := []any{"foo", 42, 3.14, true}
	generic, typed := Attrs(input...)
	if len(generic) != 4 {
		t.Errorf("expected 4 generic attrs, got %d", len(generic))
	}
	if len(typed) != 0 {
		t.Errorf("expected no typed attrs, got %d", len(typed))
	}
}

func TestAttrsMixed(t *testing.T) {
	input := []any{
		"plain string",
		String("key", "value"),
		42,
		Int("num", 7),
		true,
		Bool("flag", false),
	}
	generic, typed := Attrs(input...)
	if len(generic) != 3 {
		t.Errorf("expected 3 generic attrs, got %d", len(generic))
	}
	if len(typed) != 3 {
		t.Errorf("expected 3 typed attrs, got %d", len(typed))
	}
	// Verify order is preserved within each group.
	if generic[0] != "plain string" {
		t.Errorf("expected first generic to be %q, got %v", "plain string", generic[0])
	}
	if generic[1] != 42 {
		t.Errorf("expected second generic to be %d, got %v", 42, generic[1])
	}
	if generic[2] != true {
		t.Errorf("expected third generic to be %t, got %v", true, generic[2])
	}
}

func TestAttrsEmpty(t *testing.T) {
	generic, typed := Attrs()
	if generic != nil {
		t.Errorf("expected nil generic, got %v", generic)
	}
	if typed != nil {
		t.Errorf("expected nil typed, got %v", typed)
	}
}
