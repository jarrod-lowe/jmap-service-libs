package plugincontract

// Args represents method arguments or response data in JMAP plugin communication.
// It provides type-safe accessor methods for retrieving values.
type Args map[string]any

// String returns the string value for the given key.
// Returns false if the key doesn't exist or the value is not a string.
func (a Args) String(key string) (string, bool) {
	if a == nil {
		return "", false
	}
	v, exists := a[key]
	if !exists {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// StringOr returns the string value for the given key, or the default value
// if the key doesn't exist or the value is not a string.
func (a Args) StringOr(key, defaultVal string) string {
	if v, ok := a.String(key); ok {
		return v
	}
	return defaultVal
}

// Int returns the int64 value for the given key.
// Handles JSON numbers (which unmarshal as float64) by converting to int64.
// Returns false if the key doesn't exist or the value is not numeric.
func (a Args) Int(key string) (int64, bool) {
	if a == nil {
		return 0, false
	}
	v, exists := a[key]
	if !exists {
		return 0, false
	}
	switch n := v.(type) {
	case int64:
		return n, true
	case float64:
		return int64(n), true
	case int:
		return int64(n), true
	default:
		return 0, false
	}
}

// IntOr returns the int64 value for the given key, or the default value
// if the key doesn't exist or the value is not numeric.
func (a Args) IntOr(key string, defaultVal int64) int64 {
	if v, ok := a.Int(key); ok {
		return v
	}
	return defaultVal
}

// Float returns the float64 value for the given key.
// Returns false if the key doesn't exist or the value is not numeric.
func (a Args) Float(key string) (float64, bool) {
	if a == nil {
		return 0, false
	}
	v, exists := a[key]
	if !exists {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return n, true
	case int64:
		return float64(n), true
	case int:
		return float64(n), true
	default:
		return 0, false
	}
}

// Bool returns the bool value for the given key.
// Returns false if the key doesn't exist or the value is not a bool.
func (a Args) Bool(key string) (bool, bool) {
	if a == nil {
		return false, false
	}
	v, exists := a[key]
	if !exists {
		return false, false
	}
	b, ok := v.(bool)
	return b, ok
}

// BoolOr returns the bool value for the given key, or the default value
// if the key doesn't exist or the value is not a bool.
func (a Args) BoolOr(key string, defaultVal bool) bool {
	if v, ok := a.Bool(key); ok {
		return v
	}
	return defaultVal
}

// StringSlice returns the []string value for the given key.
// Returns false if the key doesn't exist, the value is not a slice,
// or any element in the slice is not a string.
func (a Args) StringSlice(key string) ([]string, bool) {
	if a == nil {
		return nil, false
	}
	v, exists := a[key]
	if !exists {
		return nil, false
	}
	slice, ok := v.([]any)
	if !ok {
		return nil, false
	}
	result := make([]string, len(slice))
	for i, elem := range slice {
		s, ok := elem.(string)
		if !ok {
			return nil, false
		}
		result[i] = s
	}
	return result, true
}

// Object returns the nested Args value for the given key.
// Returns false if the key doesn't exist or the value is not a map[string]any.
func (a Args) Object(key string) (Args, bool) {
	if a == nil {
		return nil, false
	}
	v, exists := a[key]
	if !exists {
		return nil, false
	}
	m, ok := v.(map[string]any)
	if !ok {
		return nil, false
	}
	return Args(m), true
}

// Has returns true if the key exists in the Args map.
// Returns true even if the value is nil.
func (a Args) Has(key string) bool {
	if a == nil {
		return false
	}
	_, exists := a[key]
	return exists
}
