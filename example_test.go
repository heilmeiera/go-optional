package optional_test

import (
	"encoding/json"
	"fmt"

	optional "github.com/heilmeiera/go-optional"
)

func Example() {
	a := optional.Some(42)
	b := optional.None[int]()

	fmt.Println(a.IsSome(), b.IsSome())
	fmt.Println(a.Unwrap(), b.UnwrapOr(-1))
	fmt.Println(a, b)
	// Output:
	// true false
	// 42 -1
	// Some[42] None[]
}

func ExampleOption_Get() {
	retries := optional.Some(5)
	if n, ok := retries.Get(); ok && n > 3 {
		fmt.Println("more than three retries:", n)
	}
	// Output: more than three retries: 5
}

func Example_json() {
	type config struct {
		Retries optional.Option[int] `json:"retries"`
		Timeout optional.Option[int] `json:"timeout"`
	}

	var cfg config
	if err := json.Unmarshal([]byte(`{"retries": null}`), &cfg); err != nil {
		panic(err)
	}
	fmt.Println(cfg.Retries, cfg.Timeout) // null and absent are both None here

	cfg.Timeout = optional.Some(30)
	out, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	// Output:
	// None[] None[]
	// {"retries":null,"timeout":30}
}
