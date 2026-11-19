# go-optional

A minimal, dependency-free `Option[T]` type for Go — a generic container that either holds a value
(`Some`) or is empty (`None`), with JSON support for nullable fields.

## Install

```sh
go get github.com/heilmeiera/go-optional
```

## Usage

```go
package main

import (
	"fmt"

	optional "github.com/heilmeiera/go-optional"
)

func main() {
	a := optional.Some(42)
	b := optional.None[int]()

	if v, ok := a.Get(); ok {
		fmt.Println(v) // 42
	}
	fmt.Println(b.UnwrapOr(-1)) // -1
	fmt.Println(b)              // None[]
}
```

`Get` is the comma-ok form. `Unwrap` returns the value and panics on a `None` option; `UnwrapOr`
returns a fallback instead.

## JSON

`Option[T]` implements `json.Marshaler` and `json.Unmarshaler`, which makes it a natural type for
nullable or optional fields in configuration and parameter files:

```go
type Config struct {
	Retries optional.Option[int] `json:"retries"`          // absent or null: None
	Timeout optional.Option[int] `json:"timeout,omitzero"` // None is left out when marshalling
}
```

- A `None` option marshals to `null`, and `null` unmarshals to `None` — also into an option that held
  a value, unlike `encoding/json`'s convention of leaving a value unchanged on `null`.
- A field missing from the JSON leaves its option as it was: `None` in a freshly declared struct.
- To leave `None` fields out when marshalling, tag them `omitzero` (Go 1.24 or newer). `omitempty`
  does not apply to structs and writes `null`.
- A `Some` option whose value marshals to `null` itself — a nil pointer, slice or map — comes back as
  `None`.

## Design

`Option[T]` uses the comma-ok representation (`value T` + `ok bool`) rather than an internal
pointer, so the zero value is a valid `None` and an option — including any struct embedding one —
is copied deeply by an ordinary assignment.

## License

MIT — see [LICENSE](LICENSE).
