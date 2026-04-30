// Package attrs provides typed key-value attributes for structured logging.
// Attributes are used to attach additional metadata to log entries in a
// type-safe manner, supporting common Go types such as string, int, float64,
// bool, time.Time, and time.Duration.
package attrs

import (
	"fmt"
	"time"
)

// AttrValue represents a typed key-value attribute. The type parameter T
// determines the value type stored by the attribute.
type AttrValue[T any] interface {
	// Key returns the attribute's key name.
	Key() string
	// Value returns the attribute's typed value.
	Value() T
}

type attr[T any] struct {
	key   string
	value T
}

// Key returns the attribute's key name.
func (s attr[T]) Key() string {
	return s.key
}

// Value returns the attribute's typed value.
func (s attr[T]) Value() T {
	return s.value
}

// String creates a string-typed attribute with the given key and value.
func String(k, v string) AttrValue[string] {
	return attr[string]{key: k, value: v}
}

// Int creates an int-typed attribute with the given key and value.
func Int(k string, v int) AttrValue[int] {
	return attr[int]{key: k, value: v}
}

// Int64 creates an int64-typed attribute with the given key and value.
func Int64(k string, v int64) AttrValue[int64] {
	return attr[int64]{key: k, value: v}
}

// Uint64 creates a uint64-typed attribute with the given key and value.
func Uint64(k string, v uint64) AttrValue[uint64] {
	return attr[uint64]{key: k, value: v}
}

// Float64 creates a float64-typed attribute with the given key and value.
func Float64(k string, v float64) AttrValue[float64] {
	return attr[float64]{key: k, value: v}
}

// Bool creates a bool-typed attribute with the given key and value.
func Bool(k string, v bool) AttrValue[bool] {
	return attr[bool]{key: k, value: v}
}

// Time creates a time.Time-typed attribute with the given key and value.
func Time(k string, v time.Time) AttrValue[time.Time] {
	return attr[time.Time]{key: k, value: v}
}

// Duration creates a time.Duration-typed attribute with the given key and value.
func Duration(k string, v time.Duration) AttrValue[time.Duration] {
	return attr[time.Duration]{key: k, value: v}
}

// Any creates an attribute with an arbitrary value type. Use this for values
// that do not fit one of the other typed constructors.
func Any(k string, v any) AttrValue[any] {
	return attr[any]{key: k, value: v}
}

// Valid checks that all attributes are of a valid type.
func Valid(attrs []any) error {
	for _, attr := range attrs {
		switch a := attr.(type) {
		case AttrValue[string]:
		case AttrValue[int]:
		case AttrValue[int64]:
		case AttrValue[uint64]:
		case AttrValue[float64]:
		case AttrValue[bool]:
		case AttrValue[time.Time]:
		case AttrValue[time.Duration]:
		case AttrValue[any]:
		default:
			return fmt.Errorf("invalid attribute type %T", a)
		}
	}
	return nil
}

// Attrs separates attributes into generic and typed attributes, but keeps
// them in the same order as provided.
func Attrs(attrs ...any) ([]any, []any) {
	var (
		a []any
		b []any
	)
	for _, attr := range attrs {
		switch attr.(type) {
		case AttrValue[string]:
			b = append(b, attr)
		case AttrValue[int]:
			b = append(b, attr)
		case AttrValue[int64]:
			b = append(b, attr)
		case AttrValue[uint64]:
			b = append(b, attr)
		case AttrValue[float64]:
			b = append(b, attr)
		case AttrValue[bool]:
			b = append(b, attr)
		case AttrValue[time.Time]:
			b = append(b, attr)
		case AttrValue[time.Duration]:
			b = append(b, attr)
		case AttrValue[any]:
			b = append(b, attr)
		default:
			a = append(a, attr)
		}
	}
	return a, b
}
