package attrs

import (
	"fmt"
	"time"
)

type AttrValue[T any] interface {
	Key() string
	Value() T
}

type attr[T any] struct {
	key   string
	value T
}

func (s attr[T]) Key() string {
	return s.key
}

func (s attr[T]) Value() T {
	return s.value
}

func String(k, v string) AttrValue[string] {
	return attr[string]{key: k, value: v}
}

func Int(k string, v int) AttrValue[int] {
	return attr[int]{key: k, value: v}
}

func Int64(k string, v int64) AttrValue[int64] {
	return attr[int64]{key: k, value: v}
}

func Uint64(k string, v uint64) AttrValue[uint64] {
	return attr[uint64]{key: k, value: v}
}

func Float64(k string, v float64) AttrValue[float64] {
	return attr[float64]{key: k, value: v}
}

func Bool(k string, v bool) AttrValue[bool] {
	return attr[bool]{key: k, value: v}
}

func Time(k string, v time.Time) AttrValue[time.Time] {
	return attr[time.Time]{key: k, value: v}
}

func Duration(k string, v time.Duration) AttrValue[time.Duration] {
	return attr[time.Duration]{key: k, value: v}
}

func Any(k string, v any) AttrValue[any] {
	return attr[any]{key: k, value: v}
}

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
