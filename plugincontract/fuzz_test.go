package plugincontract

import (
	"testing"
)

// FuzzArgsString verifies that Args.String never panics on arbitrary input.
func FuzzArgsString(f *testing.F) {
	f.Add("key", "value")
	f.Add("", "")
	f.Add("key", "123")

	f.Fuzz(func(t *testing.T, key, value string) {
		args := Args{key: value}
		args.String(key)
		args.StringOr(key, "default")

		// Also test with nil value
		args2 := Args{key: nil}
		args2.String(key)
		args2.StringOr(key, "default")
	})
}

// FuzzArgsInt verifies that Args.Int never panics on arbitrary numeric-like input.
func FuzzArgsInt(f *testing.F) {
	f.Add("key", float64(42))
	f.Add("key", float64(-1))
	f.Add("key", float64(0))
	f.Add("", float64(1e18))

	f.Fuzz(func(t *testing.T, key string, value float64) {
		args := Args{key: value}
		args.Int(key)
		args.IntOr(key, 0)

		// Also test with non-numeric values
		args2 := Args{key: "not a number"}
		args2.Int(key)
		args2.IntOr(key, 0)

		// Test with nil
		args3 := Args{key: nil}
		args3.Int(key)
		args3.IntOr(key, 0)
	})
}

// FuzzArgsFloat verifies that Args.Float never panics on arbitrary input.
func FuzzArgsFloat(f *testing.F) {
	f.Add("key", float64(3.14))
	f.Add("key", float64(-1e308))
	f.Add("key", float64(0))

	f.Fuzz(func(t *testing.T, key string, value float64) {
		args := Args{key: value}
		args.Float(key)

		// Also test with non-float values
		args2 := Args{key: "not a float"}
		args2.Float(key)

		// Test with nil
		args3 := Args{key: nil}
		args3.Float(key)
	})
}

// FuzzArgsBool verifies that Args.Bool never panics on arbitrary input.
func FuzzArgsBool(f *testing.F) {
	f.Add("key", true)
	f.Add("key", false)
	f.Add("", true)

	f.Fuzz(func(t *testing.T, key string, value bool) {
		args := Args{key: value}
		args.Bool(key)
		args.BoolOr(key, false)

		// Also test with non-bool values
		args2 := Args{key: "true"}
		args2.Bool(key)
		args2.BoolOr(key, false)

		// Test with nil
		args3 := Args{key: nil}
		args3.Bool(key)
		args3.BoolOr(key, false)
	})
}

// FuzzArgsStringSlice verifies that Args.StringSlice never panics on arbitrary input.
func FuzzArgsStringSlice(f *testing.F) {
	f.Add("key", "elem1")
	f.Add("key", "")
	f.Add("", "test")

	f.Fuzz(func(t *testing.T, key, elem string) {
		// Test with a valid string slice
		args := Args{key: []any{elem}}
		args.StringSlice(key)

		// Test with mixed-type slice
		args2 := Args{key: []any{elem, 42, nil, true}}
		args2.StringSlice(key)

		// Test with non-slice value
		args3 := Args{key: elem}
		args3.StringSlice(key)

		// Test with nil value
		args4 := Args{key: nil}
		args4.StringSlice(key)

		// Test with empty slice
		args5 := Args{key: []any{}}
		args5.StringSlice(key)
	})
}

// FuzzArgsObject verifies that Args.Object never panics on arbitrary input.
func FuzzArgsObject(f *testing.F) {
	f.Add("key", "innerKey", "innerValue")
	f.Add("", "", "")

	f.Fuzz(func(t *testing.T, key, innerKey, innerValue string) {
		// Test with a valid nested object
		args := Args{key: map[string]any{innerKey: innerValue}}
		obj, ok := args.Object(key)
		if ok {
			obj.String(innerKey)
			obj.Has(innerKey)
		}

		// Test with non-object value
		args2 := Args{key: innerValue}
		args2.Object(key)

		// Test with nil value
		args3 := Args{key: nil}
		args3.Object(key)
	})
}

// FuzzArgsHas verifies that Args.Has never panics on arbitrary input.
func FuzzArgsHas(f *testing.F) {
	f.Add("key")
	f.Add("")

	f.Fuzz(func(t *testing.T, key string) {
		args := Args{key: "value"}
		args.Has(key)
		args.Has(key + "_missing")

		var nilArgs Args
		nilArgs.Has(key)
	})
}
