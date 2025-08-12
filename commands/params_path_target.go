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
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// GetAllValuesAtTarget retrieves all values from a nested data structure (structs, slices, arrays, or maps)
// according to a dot-separated pattern string. The pattern can include field names, numeric indices for slices/arrays,
// "{index}" to iterate all elements of a slice/array, and "{key}" to iterate all keys of a map.
// For example, pattern "Users.{index}.Name" will return the Name field of all elements in the Users slice.
// Returns a slice of interface{} containing all matched values, or an error if the pattern is invalid or traversal fails.
func GetAllValuesAtTarget(params interface{}, target string) ([]interface{}, error) {
	if params == nil {
		return nil, fmt.Errorf("params is nil")
	}
	if target == "" {
		return []interface{}{params}, nil
	}
	parts := strings.Split(target, ".")
	return getAllValuesRecursive(reflect.ValueOf(params), parts)
}

// getAllValuesRecursive is a recursive helper that traverses the given reflect.Value according to the provided pattern parts.
// It supports struct fields (case-insensitive), slice/array indices or iteration, and map keys or iteration.
// Returns a slice of interface{} with all values matching the pattern, or an error if traversal fails.
func getAllValuesRecursive(val reflect.Value, parts []string) ([]interface{}, error) {
	val = derefValue(val)
	if len(parts) == 0 {
		return []interface{}{val.Interface()}, nil
	}
	part := parts[0]
	rest := parts[1:]

	switch val.Kind() {
	case reflect.Struct:
		field := fieldByLowerName(val, part)
		if !field.IsValid() {
			return nil, fmt.Errorf("field '%s' not found in struct", part)
		}
		return getAllValuesRecursive(field, rest)
	case reflect.Slice, reflect.Array:
		if part == "{index}" {
			var results []interface{}
			for i := 0; i < val.Len(); i++ {
				sub, err := getAllValuesRecursive(val.Index(i), rest)
				if err != nil {
					return nil, err
				}
				results = append(results, sub...)
			}
			return results, nil
		}
		idx, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("expected index or '{index}' for slice/array, got '%s'", part)
		}
		if idx < 0 || idx >= val.Len() {
			return nil, fmt.Errorf("index %d out of bounds", idx)
		}
		return getAllValuesRecursive(val.Index(idx), rest)
	case reflect.Map:
		keyKind := val.Type().Key().Kind()
		var mapKey reflect.Value
		var restKey []string

		if part == "{key}" {
			var results []interface{}
			for _, key := range val.MapKeys() {
				elem := val.MapIndex(key)
				sub, err := getAllValuesRecursive(elem, rest)
				if err != nil {
					return nil, err
				}
				results = append(results, sub...)
			}
			return results, nil
		}

		if keyKind == reflect.String {
			if strings.Contains(part, ".") {
				return nil, fmt.Errorf("map key '%s' contains a dot, which is not allowed for string map keys", part)
			}
			mapKey = reflect.ValueOf(part)
			restKey = rest
		} else {
			// Pour les clés non-string, recompose la clé sur plusieurs parts
			var lastConvErr error
			found := false
			for n := 1; n <= len(parts); n++ {
				tryKey := strings.Join(parts[:n], ".")
				mapKeyConv, err := convertStringToType(tryKey, val.Type().Key())
				if err == nil {
					mapKey = reflect.ValueOf(mapKeyConv)
					restKey = parts[n:]
					found = true
					break
				} else {
					lastConvErr = err
				}
			}
			if !found {
				// Retourne l'erreur de conversion la plus récente (plus précise)
				return nil, lastConvErr
			}
		}
		elem := val.MapIndex(mapKey)
		if !elem.IsValid() {
			return nil, fmt.Errorf("map key '%v' not found", mapKey.Interface())
		}
		return getAllValuesRecursive(elem, restKey)
	default:
		return nil, fmt.Errorf("cannot traverse kind %v", val.Kind())
	}
}
