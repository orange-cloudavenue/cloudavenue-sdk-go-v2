/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package commands

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/orange-cloudavenue/common-go/strcase"
)

// GetValueAtPath retrieves the value from a Go struct or nested data structure
// using a fully-resolved ParamSpecNameNotation path. This function supports
// accessing fields within structs, elements within slices/arrays, and values
// within maps, traversing deeply nested structures as needed.
//
// See params_path_utils.go for shared helpers.
func GetValueAtPath(params interface{}, path string) (interface{}, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}
	if path == "" {
		return params, nil
	}

	parts := strings.Split(path, ".")
	val := reflect.ValueOf(params)
	val = derefValue(val)

	for i, part := range parts {
		if !val.IsValid() {
			return nil, fmt.Errorf("invalid path at '%s' (nil or non-existent)", part)
		}

		switch val.Kind() {
		case reflect.Struct:
			field := fieldByLowerName(val, part)
			if !field.IsValid() {
				return nil, fmt.Errorf("field '%s' not found in struct at '%s'", part, strings.Join(parts[:i], "."))
			}
			val = field
		case reflect.Ptr:
			if val.IsNil() {
				return nil, fmt.Errorf("nil pointer at '%s'", strings.Join(parts[:i], "."))
			}
			val = val.Elem()
			i-- // retry this part with dereferenced value
			continue
		case reflect.Slice, reflect.Array:
			index, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("expected integer index for slice/array at '%s', got '%s'", strings.Join(parts[:i], "."), part)
			}
			if index < 0 || index >= val.Len() {
				return nil, fmt.Errorf("index %d out of bounds at '%s'", index, strings.Join(parts[:i], "."))
			}
			val = val.Index(index)
		case reflect.Map:
			mapKey := reflect.ValueOf(part)
			if val.Type().Key().Kind() != reflect.String {
				mapKeyConv, err := convertStringToType(part, val.Type().Key())
				if err != nil {
					return nil, fmt.Errorf("cannot convert key '%s' to %v at '%s': %w", part, val.Type().Key(), strings.Join(parts[:i], "."), err)
				}
				mapKey = reflect.ValueOf(mapKeyConv)
			}
			elem := val.MapIndex(mapKey)
			if !elem.IsValid() {
				return nil, fmt.Errorf("map key '%s' not found at '%s'", part, strings.Join(parts[:i], "."))
			}
			val = elem
		default:
			return nil, fmt.Errorf("cannot traverse kind %v at '%s'", val.Kind(), strings.Join(parts[:i], "."))
		}

		// Dereference pointer for intermediate values, except for the last element
		if i < len(parts)-1 {
			val = derefValue(val)
		}
	}

	// Dereference if final value is a pointer
	val = derefValue(val)
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return nil, errors.New("getvalueatpath: final value is nil pointer")
	}

	return val.Interface(), nil
}

// StoreValueAtPath sets the value at the specified path in the params struct.
func StoreValueAtPath(params interface{}, path, value string) error {
	if params == nil {
		return errors.New("params is nil")
	}
	if path == "" {
		return errors.New("path is empty")
	}

	parts := strings.Split(path, ".")
	val := reflect.ValueOf(params)
	val = derefValue(val)

	for i, part := range parts {
		if !val.IsValid() {
			return fmt.Errorf("invalid path at '%s' (nil or non-existent)", part)
		}

		switch val.Kind() {
		case reflect.Struct:
			field := fieldByLowerName(val, part)
			if !field.IsValid() {
				return fmt.Errorf("field '%s' not found in struct at '%s'", part, strings.Join(parts[:i], "."))
			}
			val = field
		case reflect.Ptr:
			if val.IsNil() {
				return fmt.Errorf("nil pointer at '%s'", strings.Join(parts[:i], "."))
			}
			val = val.Elem()
			i-- // retry this part with dereferenced value
			continue
		case reflect.Slice, reflect.Array:
			index, err := strconv.Atoi(part)
			if err != nil {
				return fmt.Errorf("expected integer index for slice/array at '%s', got '%s'", strings.Join(parts[:i], "."), part)
			}
			if index < 0 || index >= val.Len() {
				// Add init the slice
				newElem := reflect.New(val.Type().Elem()).Elem()
				val.Set(reflect.Append(val, newElem))
			}
			val = val.Index(index)
		case reflect.Map:
			keyVal := reflect.ValueOf(part)
			if val.Type().Key().Kind() != reflect.String {
				keyConv, err := convertStringToType(part, val.Type().Key())
				if err != nil {
					return fmt.Errorf("cannot convert key '%s' to %v at '%s': %w", part, val.Type().Key(), strings.Join(parts[:i], "."), err)
				}
				keyVal = reflect.ValueOf(keyConv)
			}
			elem := val.MapIndex(keyVal)
			if !elem.IsValid() {
				return fmt.Errorf("map key '%v' not found at '%s'", part, strings.Join(parts[:i], "."))
			}
			// If this is the last part, or next is a struct field, handle set
			if i == len(parts)-2 && elem.Kind() == reflect.Struct {
				// Make addressable copy of struct
				addr := reflect.New(elem.Type())
				addr.Elem().Set(elem)
				field := addr.Elem().FieldByName(strcase.ToPublicGoName(parts[i+1]))
				if !field.IsValid() {
					return fmt.Errorf("field '%s' not found in struct at '%s'", parts[i+1], strings.Join(parts[:i+1], "."))
				}
				if !field.CanSet() {
					return fmt.Errorf("cannot set field '%s' at '%s'", parts[i+1], strings.Join(parts[:i+1], "."))
				}
				valToSet := reflect.ValueOf(value)
				if valToSet.Type().AssignableTo(field.Type()) {
					field.Set(valToSet)
				} else if valToSet.Type().ConvertibleTo(field.Type()) {
					field.Set(valToSet.Convert(field.Type()))
				} else {
					return fmt.Errorf("cannot assign value of type %T to field '%s' of type %s", value, parts[i+1], field.Type())
				}
				// Write back updated struct to map
				val.SetMapIndex(keyVal, addr.Elem())
				return nil
			}
			// For intermediate struct map values, make addressable copy for further traversal
			if elem.Kind() == reflect.Struct {
				elemAddr := reflect.New(elem.Type())
				elemAddr.Elem().Set(elem)
				val = elemAddr.Elem()
			} else {
				val = elem
			}
		default:
			return fmt.Errorf("cannot traverse kind %v at '%s'", val.Kind(), strings.Join(parts[:i], "."))
		}

		// Dereference pointer for intermediate values, except for the last element
		if i < len(parts)-1 {
			val = derefValue(val)
		}
	}

	// Dereference if final value is a pointer and not nil
	val = derefValue(val)

	x, err := convertStringToType(value, val.Type())
	if err != nil {
		return fmt.Errorf("cannot convert value '%s' to type %s: %w", value, val.Type(), err)
	}
	if val.Kind() == reflect.Ptr {
		// Pour les pointeurs, créer un nouveau pointer vers la valeur convertible
		ptr := reflect.New(val.Type().Elem())
		ptr.Elem().Set(reflect.ValueOf(x))
		val.Set(ptr)
	} else {
		val.Set(reflect.ValueOf(x))
	}

	return nil
}
