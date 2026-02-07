package plugincontract_test

import (
	"fmt"

	"github.com/jarrod-lowe/jmap-service-libs/plugincontract"
)

func ExampleArgs_String() {
	args := plugincontract.Args{"name": "inbox"}
	val, ok := args.String("name")
	fmt.Println(val, ok)
	// Output: inbox true
}

func ExampleArgs_StringOr() {
	args := plugincontract.Args{}
	val := args.StringOr("name", "default")
	fmt.Println(val)
	// Output: default
}

func ExampleArgs_Int() {
	// JSON numbers unmarshal as float64; Int handles the conversion.
	args := plugincontract.Args{"count": float64(42)}
	val, ok := args.Int("count")
	fmt.Println(val, ok)
	// Output: 42 true
}

func ExampleArgs_Float() {
	args := plugincontract.Args{"ratio": float64(3.14)}
	val, ok := args.Float("ratio")
	fmt.Println(val, ok)
	// Output: 3.14 true
}

func ExampleArgs_Bool() {
	args := plugincontract.Args{"enabled": true}
	val, ok := args.Bool("enabled")
	fmt.Println(val, ok)
	// Output: true true
}

func ExampleArgs_StringSlice() {
	args := plugincontract.Args{"ids": []any{"a", "b", "c"}}
	val, ok := args.StringSlice("ids")
	fmt.Println(val, ok)
	// Output: [a b c] true
}

func ExampleArgs_Object() {
	args := plugincontract.Args{
		"filter": map[string]any{"field": "subject"},
	}
	obj, ok := args.Object("filter")
	if ok {
		field, _ := obj.String("field")
		fmt.Println(field)
	}
	// Output: subject
}

func ExampleArgs_Has() {
	args := plugincontract.Args{"key": nil}
	fmt.Println(args.Has("key"))
	fmt.Println(args.Has("missing"))
	// Output:
	// true
	// false
}
