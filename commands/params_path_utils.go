package commands

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// derefValue returns the value after recursively dereferencing pointers (if not nil).
func derefValue(val reflect.Value) reflect.Value {
	for val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}
	return val
}

// fieldByLowerName returns the struct field value by matching the field name (case-insensitive, always lowercase ParamSpecNameNotation).
func fieldByLowerName(val reflect.Value, name string) reflect.Value {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// Support: CamelCase, lower, snake_case
		goName := field.Name
		if strings.EqualFold(goName, name) ||
			strings.EqualFold(toSnakeCase(goName), name) ||
			strings.EqualFold(strings.ToLower(goName), name) {
			return val.Field(i)
		}
	}
	return reflect.Value{}
}

// convertStringToType attempts to convert a string to the given reflect.Type (for map keys).
func convertStringToType(s string, t reflect.Type) (interface{}, error) {
	switch t.Kind() {
	case reflect.String:
		return s, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(i).Convert(t).Interface(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(u).Convert(t).Interface(), nil
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(f).Convert(t).Interface(), nil
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}
		return b, nil
	default:
		return nil, fmt.Errorf("unsupported map key type: %v", t)
	}
}
