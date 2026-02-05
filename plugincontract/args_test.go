package plugincontract

import (
	"testing"
)

func TestArgs_String(t *testing.T) {
	args := Args{
		"name":   "test",
		"number": 42,
		"empty":  "",
	}

	t.Run("returns string value", func(t *testing.T) {
		val, ok := args.String("name")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if val != "test" {
			t.Errorf("expected 'test', got %q", val)
		}
	})

	t.Run("returns empty string when value is empty", func(t *testing.T) {
		val, ok := args.String("empty")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if val != "" {
			t.Errorf("expected empty string, got %q", val)
		}
	})

	t.Run("returns false for missing key", func(t *testing.T) {
		_, ok := args.String("missing")
		if ok {
			t.Fatal("expected ok to be false for missing key")
		}
	})

	t.Run("returns false for non-string value", func(t *testing.T) {
		_, ok := args.String("number")
		if ok {
			t.Fatal("expected ok to be false for non-string value")
		}
	})
}

func TestArgs_StringOr(t *testing.T) {
	args := Args{
		"name":   "test",
		"number": 42,
	}

	t.Run("returns string value when present", func(t *testing.T) {
		val := args.StringOr("name", "default")
		if val != "test" {
			t.Errorf("expected 'test', got %q", val)
		}
	})

	t.Run("returns default for missing key", func(t *testing.T) {
		val := args.StringOr("missing", "default")
		if val != "default" {
			t.Errorf("expected 'default', got %q", val)
		}
	})

	t.Run("returns default for non-string value", func(t *testing.T) {
		val := args.StringOr("number", "default")
		if val != "default" {
			t.Errorf("expected 'default', got %q", val)
		}
	})
}

func TestArgs_Int(t *testing.T) {
	args := Args{
		"int":      int64(42),
		"float":    float64(99),
		"intval":   int(10),
		"string":   "not a number",
		"zero":     float64(0),
		"negfloat": float64(-5),
	}

	t.Run("returns int64 value", func(t *testing.T) {
		val, ok := args.Int("int")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if val != 42 {
			t.Errorf("expected 42, got %d", val)
		}
	})

	t.Run("converts float64 to int64 (JSON numbers)", func(t *testing.T) {
		val, ok := args.Int("float")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if val != 99 {
			t.Errorf("expected 99, got %d", val)
		}
	})

	t.Run("converts int to int64", func(t *testing.T) {
		val, ok := args.Int("intval")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if val != 10 {
			t.Errorf("expected 10, got %d", val)
		}
	})

	t.Run("handles zero", func(t *testing.T) {
		val, ok := args.Int("zero")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if val != 0 {
			t.Errorf("expected 0, got %d", val)
		}
	})

	t.Run("handles negative numbers", func(t *testing.T) {
		val, ok := args.Int("negfloat")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if val != -5 {
			t.Errorf("expected -5, got %d", val)
		}
	})

	t.Run("returns false for missing key", func(t *testing.T) {
		_, ok := args.Int("missing")
		if ok {
			t.Fatal("expected ok to be false for missing key")
		}
	})

	t.Run("returns false for non-numeric value", func(t *testing.T) {
		_, ok := args.Int("string")
		if ok {
			t.Fatal("expected ok to be false for non-numeric value")
		}
	})
}

func TestArgs_IntOr(t *testing.T) {
	args := Args{
		"count":  float64(5),
		"string": "not a number",
	}

	t.Run("returns int value when present", func(t *testing.T) {
		val := args.IntOr("count", 0)
		if val != 5 {
			t.Errorf("expected 5, got %d", val)
		}
	})

	t.Run("returns default for missing key", func(t *testing.T) {
		val := args.IntOr("missing", 10)
		if val != 10 {
			t.Errorf("expected 10, got %d", val)
		}
	})

	t.Run("returns default for non-numeric value", func(t *testing.T) {
		val := args.IntOr("string", 20)
		if val != 20 {
			t.Errorf("expected 20, got %d", val)
		}
	})
}

func TestArgs_Float(t *testing.T) {
	args := Args{
		"pi":     float64(3.14),
		"int":    int64(42),
		"intval": int(10),
		"string": "not a number",
	}

	t.Run("returns float64 value", func(t *testing.T) {
		val, ok := args.Float("pi")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if val != 3.14 {
			t.Errorf("expected 3.14, got %f", val)
		}
	})

	t.Run("converts int64 to float64", func(t *testing.T) {
		val, ok := args.Float("int")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if val != 42.0 {
			t.Errorf("expected 42.0, got %f", val)
		}
	})

	t.Run("converts int to float64", func(t *testing.T) {
		val, ok := args.Float("intval")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if val != 10.0 {
			t.Errorf("expected 10.0, got %f", val)
		}
	})

	t.Run("returns false for missing key", func(t *testing.T) {
		_, ok := args.Float("missing")
		if ok {
			t.Fatal("expected ok to be false for missing key")
		}
	})

	t.Run("returns false for non-numeric value", func(t *testing.T) {
		_, ok := args.Float("string")
		if ok {
			t.Fatal("expected ok to be false for non-numeric value")
		}
	})
}

func TestArgs_Bool(t *testing.T) {
	args := Args{
		"enabled":  true,
		"disabled": false,
		"string":   "true",
	}

	t.Run("returns true value", func(t *testing.T) {
		val, ok := args.Bool("enabled")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if !val {
			t.Error("expected true")
		}
	})

	t.Run("returns false value", func(t *testing.T) {
		val, ok := args.Bool("disabled")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if val {
			t.Error("expected false")
		}
	})

	t.Run("returns false for missing key", func(t *testing.T) {
		_, ok := args.Bool("missing")
		if ok {
			t.Fatal("expected ok to be false for missing key")
		}
	})

	t.Run("returns false for non-bool value", func(t *testing.T) {
		_, ok := args.Bool("string")
		if ok {
			t.Fatal("expected ok to be false for non-bool value")
		}
	})
}

func TestArgs_BoolOr(t *testing.T) {
	args := Args{
		"enabled": true,
		"string":  "true",
	}

	t.Run("returns bool value when present", func(t *testing.T) {
		val := args.BoolOr("enabled", false)
		if !val {
			t.Error("expected true")
		}
	})

	t.Run("returns default for missing key", func(t *testing.T) {
		val := args.BoolOr("missing", true)
		if !val {
			t.Error("expected true (default)")
		}
	})

	t.Run("returns default for non-bool value", func(t *testing.T) {
		val := args.BoolOr("string", true)
		if !val {
			t.Error("expected true (default)")
		}
	})
}

func TestArgs_StringSlice(t *testing.T) {
	args := Args{
		"ids":    []any{"a", "b", "c"},
		"empty":  []any{},
		"mixed":  []any{"a", 1, "b"},
		"single": []any{"only"},
		"string": "not a slice",
	}

	t.Run("returns string slice", func(t *testing.T) {
		val, ok := args.StringSlice("ids")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if len(val) != 3 {
			t.Fatalf("expected 3 elements, got %d", len(val))
		}
		if val[0] != "a" || val[1] != "b" || val[2] != "c" {
			t.Errorf("expected [a, b, c], got %v", val)
		}
	})

	t.Run("returns empty slice", func(t *testing.T) {
		val, ok := args.StringSlice("empty")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if len(val) != 0 {
			t.Errorf("expected empty slice, got %v", val)
		}
	})

	t.Run("returns single element slice", func(t *testing.T) {
		val, ok := args.StringSlice("single")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		if len(val) != 1 || val[0] != "only" {
			t.Errorf("expected [only], got %v", val)
		}
	})

	t.Run("returns false for missing key", func(t *testing.T) {
		_, ok := args.StringSlice("missing")
		if ok {
			t.Fatal("expected ok to be false for missing key")
		}
	})

	t.Run("returns false for non-slice value", func(t *testing.T) {
		_, ok := args.StringSlice("string")
		if ok {
			t.Fatal("expected ok to be false for non-slice value")
		}
	})

	t.Run("returns false for slice with non-string elements", func(t *testing.T) {
		_, ok := args.StringSlice("mixed")
		if ok {
			t.Fatal("expected ok to be false for mixed slice")
		}
	})
}

func TestArgs_Object(t *testing.T) {
	args := Args{
		"nested": map[string]any{
			"key": "value",
		},
		"string": "not an object",
	}

	t.Run("returns nested Args", func(t *testing.T) {
		val, ok := args.Object("nested")
		if !ok {
			t.Fatal("expected ok to be true")
		}
		str, strOk := val.String("key")
		if !strOk {
			t.Fatal("expected nested key to exist")
		}
		if str != "value" {
			t.Errorf("expected 'value', got %q", str)
		}
	})

	t.Run("returns false for missing key", func(t *testing.T) {
		_, ok := args.Object("missing")
		if ok {
			t.Fatal("expected ok to be false for missing key")
		}
	})

	t.Run("returns false for non-object value", func(t *testing.T) {
		_, ok := args.Object("string")
		if ok {
			t.Fatal("expected ok to be false for non-object value")
		}
	})
}

func TestArgs_Has(t *testing.T) {
	args := Args{
		"exists": "value",
		"null":   nil,
	}

	t.Run("returns true for existing key", func(t *testing.T) {
		if !args.Has("exists") {
			t.Error("expected Has to return true for existing key")
		}
	})

	t.Run("returns true for nil value", func(t *testing.T) {
		if !args.Has("null") {
			t.Error("expected Has to return true for nil value")
		}
	})

	t.Run("returns false for missing key", func(t *testing.T) {
		if args.Has("missing") {
			t.Error("expected Has to return false for missing key")
		}
	})
}

func TestArgs_NilArgs(t *testing.T) {
	var args Args

	t.Run("String returns false on nil Args", func(t *testing.T) {
		_, ok := args.String("key")
		if ok {
			t.Error("expected ok to be false on nil Args")
		}
	})

	t.Run("Has returns false on nil Args", func(t *testing.T) {
		if args.Has("key") {
			t.Error("expected Has to return false on nil Args")
		}
	})
}
