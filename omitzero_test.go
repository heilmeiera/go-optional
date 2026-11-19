//go:build go1.24

package optional

import (
	"encoding/json"
	"testing"
)

// omitzero (Go 1.24) leaves a None field out; omitempty does not apply to structs.
func TestOptionOmitzero(t *testing.T) {
	type holder struct {
		A Option[int] `json:"a,omitzero"`
		B Option[int] `json:"b,omitzero"`
		C Option[int] `json:"c,omitempty"`
	}
	encoded, err := json.Marshal(holder{A: Some(0)})
	if err != nil {
		t.Fatalf("marshalling failed: %v", err)
	}
	if string(encoded) != `{"a":0,"c":null}` {
		t.Errorf(`marshalled holder = %s, want {"a":0,"c":null}`, encoded)
	}
}
