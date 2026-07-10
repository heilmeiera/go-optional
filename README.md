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

	fmt.Println(a.IsSome())     // true
	fmt.Println(a.Unwrap())     // 42
	fmt.Println(b.UnwrapOr(-1)) // -1
	fmt.Println(b)              // None[]
}
```

`Unwrap` panics on a `None` option; use `UnwrapOr` (or check `IsSome`/`IsNone` first) when a
fallback is preferable to a panic.

## JSON

`Option[T]` implements `json.Marshaler` and `json.Unmarshaler`: a `None` option marshals to `null`,
and `null` (or a missing field) unmarshals to `None` — which makes it a natural type for nullable
or optional fields in configuration and parameter files:

```go
type Config struct {
	Retries optional.Option[int] `json:"retries"` // absent or null decodes to None
}
```

## Design

`Option[T]` uses the comma-ok representation (`value T` + `ok bool`) rather than an internal
pointer, so the zero value is a valid `None` and an option — including any struct embedding one —
is copied deeply by an ordinary assignment.

## License

MIT — see [LICENSE](LICENSE).
