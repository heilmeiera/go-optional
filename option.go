// Package optional provides Option[T], a generic container that either holds a value (Some) or is
// empty (None), with JSON support: a None option marshals to null, and null (or a missing field)
// unmarshals to None.
package optional

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Option[T] holds an optional value using the comma-ok representation: a None option is the zero
// value (ok == false). Keeping it a plain value type (no internal pointer) means an Option, and any
// struct that embeds one, is copied deeply by an ordinary assignment.
type Option[T any] struct {
	value T
	ok    bool
}

// Some returns an option holding v.
func Some[T any](v T) Option[T] {
	return Option[T]{value: v, ok: true}
}

// None returns an empty option.
func None[T any]() Option[T] {
	return Option[T]{}
}

// IsSome reports whether the option holds a value.
func (o Option[T]) IsSome() bool {
	return o.ok
}

// IsNone reports whether the option is empty.
func (o Option[T]) IsNone() bool {
	return !o.ok
}

// Unwrap returns the contained value; it panics if the option is None. Use UnwrapOr when a
// fallback value is preferable to a panic.
func (o Option[T]) Unwrap() T {
	if o.IsNone() {
		panic("optional: Unwrap called on a None option")
	}
	return o.value
}

// UnwrapOr returns the contained value, or fallbackValue if the option is None.
func (o Option[T]) UnwrapOr(fallbackValue T) T {
	if o.IsNone() {
		return fallbackValue
	}
	return o.value
}

// String renders the option as "Some[value]" or "None[]". A contained value implementing
// fmt.Stringer is rendered via its String method.
func (o Option[T]) String() string {
	if o.IsNone() {
		return "None[]"
	}

	if stringer, ok := any(o.value).(fmt.Stringer); ok {
		return fmt.Sprintf("Some[%s]", stringer)
	}
	return fmt.Sprintf("Some[%v]", o.value)
}

var jsonNull = []byte("null")

// MarshalJSON encodes the contained value, or null for a None option.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsNone() {
		return jsonNull, nil
	}
	return json.Marshal(o.value)
}

// UnmarshalJSON decodes null (or an empty input) into a None option, and any other JSON value into
// a Some option via the regular unmarshalling of T.
func (o *Option[T]) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || bytes.Equal(b, jsonNull) {
		*o = None[T]()
		return nil
	}

	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*o = Some(v)
	return nil
}
