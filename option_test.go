package optional

import (
	"encoding/json"
	"testing"
)

func TestOptionIsSome(t *testing.T) {
	if None[int]().IsSome() {
		t.Error("None().IsSome() = true, want false")
	}
	if !Some(123).IsSome() {
		t.Error("Some(123).IsSome() = false, want true")
	}
}

func TestOptionIsNone(t *testing.T) {
	if !None[int]().IsNone() {
		t.Error("None().IsNone() = false, want true")
	}
	if Some(123).IsNone() {
		t.Error("Some(123).IsNone() = true, want false")
	}
}

func TestOptionUnwrap(t *testing.T) {
	if got := Some("foo").Unwrap(); got != "foo" {
		t.Errorf(`Some("foo").Unwrap() = %q, want "foo"`, got)
	}
}

func TestOptionUnwrapPanicsOnNone(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("None().Unwrap() did not panic")
		}
	}()
	None[int]().Unwrap()
}

func TestOptionUnwrapOr(t *testing.T) {
	if got := Some("foo").UnwrapOr("bar"); got != "foo" {
		t.Errorf(`Some("foo").UnwrapOr("bar") = %q, want "foo"`, got)
	}
	if got := None[string]().UnwrapOr("bar"); got != "bar" {
		t.Errorf(`None().UnwrapOr("bar") = %q, want "bar"`, got)
	}
}

func TestOptionString(t *testing.T) {
	if got := None[int]().String(); got != "None[]" {
		t.Errorf(`None().String() = %q, want "None[]"`, got)
	}
	if got := Some(7).String(); got != "Some[7]" {
		t.Errorf(`Some(7).String() = %q, want "Some[7]"`, got)
	}
	if got := Some(1.5).String(); got != "Some[1.5]" {
		t.Errorf(`Some(1.5).String() = %q, want "Some[1.5]"`, got)
	}
}

func TestOptionMarshalJSON(t *testing.T) {
	noneJSON, err := json.Marshal(None[int]())
	if err != nil {
		t.Fatalf("marshalling None failed: %v", err)
	}
	if string(noneJSON) != "null" {
		t.Errorf("marshalled None = %s, want null", noneJSON)
	}

	someJSON, err := json.Marshal(Some(42))
	if err != nil {
		t.Fatalf("marshalling Some(42) failed: %v", err)
	}
	if string(someJSON) != "42" {
		t.Errorf("marshalled Some(42) = %s, want 42", someJSON)
	}
}

func TestOptionUnmarshalJSON(t *testing.T) {
	var fromNull Option[int]
	if err := json.Unmarshal([]byte("null"), &fromNull); err != nil {
		t.Fatalf("unmarshalling null failed: %v", err)
	}
	if !fromNull.IsNone() {
		t.Error("unmarshalling null did not yield a None option")
	}

	var fromValue Option[int]
	if err := json.Unmarshal([]byte("42"), &fromValue); err != nil {
		t.Fatalf("unmarshalling 42 failed: %v", err)
	}
	if fromValue.IsNone() || fromValue.Unwrap() != 42 {
		t.Errorf("unmarshalling 42 yielded %s, want Some[42]", fromValue)
	}

	var fromBad Option[int]
	if err := json.Unmarshal([]byte(`"not-an-int"`), &fromBad); err == nil {
		t.Error("unmarshalling a string into Option[int] did not fail")
	}
}

func TestOptionJSONRoundTrip(t *testing.T) {
	type holder struct {
		A Option[string] `json:"a"`
		B Option[string] `json:"b"`
	}
	original := holder{A: Some("x"), B: None[string]()}

	encoded, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshalling failed: %v", err)
	}
	if string(encoded) != `{"a":"x","b":null}` {
		t.Errorf(`marshalled holder = %s, want {"a":"x","b":null}`, encoded)
	}

	var decoded holder
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshalling failed: %v", err)
	}
	if decoded != original {
		t.Errorf("round trip yielded %+v, want %+v", decoded, original)
	}
}
