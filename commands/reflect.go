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
	"sort"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/strcase"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
)

const (
	sliceSchema = "{index}"
	mapSchema   = "{key}"
)

// This function take a go struct and a name that comply with ParamSpec name notation (e.g "friends.{index}.name")
func GetParamType(paramType reflect.Type, name string) (reflect.Type, error) {
	var recursiveFunc func(paramType reflect.Type, parts []string) (reflect.Type, error)
	recursiveFunc = func(paramType reflect.Type, parts []string) (reflect.Type, error) {
		switch {
		case paramType.Kind() == reflect.Ptr:
			return recursiveFunc(paramType.Elem(), parts)
		case len(parts) == 0:
			return paramType, nil
		case parts[0] == sliceSchema:
			return recursiveFunc(paramType.Elem(), parts[1:])
		case parts[0] == mapSchema:
			return recursiveFunc(paramType.Elem(), parts[1:])
		default:
			// We cannot rely on dest.GetFieldByName() as reflect library is doing deep traversing when using anonymous field.
			// Because of that we should rely on our own logic
			//
			// - First we try to find a field with the correct name in the current struct
			// - If it does not exist we try to find it in all nested anonymous fields
			//   Anonymous fields are traversed from last to first as the last one in the struct declaration should take precedence

			// We construct two caches:
			anonymousFieldIndexes := []int(nil)
			fieldIndexByName := map[string]int{}
			for i := range paramType.NumField() {
				field := paramType.Field(i)
				if field.Anonymous {
					anonymousFieldIndexes = append(anonymousFieldIndexes, i)
				} else {
					fieldIndexByName[field.Name] = i
				}
			}

			// Try to find the correct field in the current struct.
			fieldName := strcase.ToPublicGoName(parts[0])
			if fieldIndex, exist := fieldIndexByName[fieldName]; exist {
				return recursiveFunc(paramType.Field(fieldIndex).Type, parts[1:])
			}

			// If it does not exist we try to find it in nested anonymous field
			for i := len(anonymousFieldIndexes) - 1; i >= 0; i-- {
				paramType, err := recursiveFunc(paramType.Field(anonymousFieldIndexes[i]).Type, parts)
				if err == nil {
					return paramType, nil
				}
			}
		}

		return nil, fmt.Errorf("count not find %s", name)
	}

	return recursiveFunc(paramType, strings.Split(name, "."))
}

// getValuesForFieldByName recursively search for fields in a cmdArgs' value and returns its values if they exist.
// The search is based on the name of the field.
func getValuesForFieldByName(value reflect.Value, parts []string) (values []reflect.Value, err error) {
	if len(parts) == 0 {
		return []reflect.Value{value}, nil
	}
	switch value.Kind() {
	case reflect.Ptr:
		return getValuesForFieldByName(value.Elem(), parts)

	case reflect.Slice:
		values := []reflect.Value(nil)
		errs := []error(nil)

		for i := range value.Len() {
			newValues, err := getValuesForFieldByName(value.Index(i), parts[1:])
			if err != nil {
				errs = append(errs, err)
			} else {
				values = append(values, newValues...)
			}
		}

		if len(values) == 0 && len(errs) != 0 {
			return nil, errors.Join(errs...)
		}

		return values, nil

	case reflect.Map:
		if value.IsNil() {
			return nil, nil
		}

		values := []reflect.Value(nil)

		mapKeys := value.MapKeys()
		sort.Slice(mapKeys, func(i, j int) bool {
			return mapKeys[i].String() < mapKeys[j].String()
		})

		for _, mapKey := range mapKeys {
			mapValue := value.MapIndex(mapKey)
			newValues, err := getValuesForFieldByName(mapValue, parts[1:])
			if err != nil {
				return nil, err
			}
			values = append(values, newValues...)
		}

		return values, nil

	case reflect.Struct:
		anonymousFieldIndexes := []int(nil)
		fieldIndexByName := map[string]int{}

		for i := range value.NumField() {
			field := value.Type().Field(i)
			if field.Anonymous {
				anonymousFieldIndexes = append(anonymousFieldIndexes, i)
			} else {
				fieldIndexByName[field.Name] = i
			}
		}

		fieldName := strcase.ToPublicGoName(parts[0])
		if fieldIndex, exist := fieldIndexByName[fieldName]; exist {
			return getValuesForFieldByName(value.Field(fieldIndex), parts[1:])
		}

		// If it does not exist we try to find it in nested anonymous field
		for _, fieldIndex := range anonymousFieldIndexes {
			newValues, err := getValuesForFieldByName(value.Field(fieldIndex), parts)
			if err == nil {
				return newValues, nil
			}
		}

		return nil, fmt.Errorf("field %v does not exist for %v", fieldName, value.Type().Name())
	}

	return nil, errors.New("case is not handled")
}

type DocModel struct {
	Object        string `json:"object"`
	Type          string `json:"type"`
	Documentation string `json:"documentation"`
}

// GetModelType reflect the model type of a command and returns it at ModelSpec notation (e.g. "friends.{index}.name").
// E.g
//
//	xx := struct {
//	  Friends []struct {
//	    ID string `json:"id"`
//	    Name string `json:"name"`
//	  } `json:"friends"`
//	}{}
//
// GetModelType(reflect.TypeOf(xx)) will return []string{"friends.{index}.id", "friends.{index}.name"}
func GetModelTypes(modelType reflect.Type) ([]DocModel, error) {
	var result []DocModel

	var walk func(t reflect.Type, prefix, doc string) error
	walk = func(t reflect.Type, prefix, doc string) error {
		// Dereference pointer
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		switch t.Kind() {
		case reflect.Struct:
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				// Skip unexported fields
				if field.PkgPath != "" {
					continue
				}
				// Get json tag or fallback to field name
				fieldName := strcase.ToSnake(field.Tag.Get("json"))
				if fieldName == "" {
					fieldName = strcase.ToSnake(field.Name)
				} else {
					// Remove omitempty or other tag options
					fieldName = strings.Split(fieldName, ",")[0]
				}
				if fieldName == "-" {
					continue
				}

				var path string
				if prefix != "" {
					path = prefix + "." + fieldName
				} else {
					path = fieldName
				}
				// Recurse into anonymous fields as if they were part of the parent
				if field.Anonymous {
					if err := walk(field.Type, prefix, ""); err != nil {
						return err
					}
					continue
				}

				if err := walk(field.Type, path, field.Tag.Get("documentation")); err != nil {
					return err
				}
			}
		case reflect.Slice:
			// Add .{index} to path
			if prefix == "" {
				return nil // can't have slice at root
			}
			// If the slice is not a struct we cannot recurse into it
			if t.Elem().Kind() != reflect.Struct {
				result = append(result, DocModel{
					Object:        prefix + "." + sliceSchema,
					Type:          t.Elem().String(),
					Documentation: doc,
				})
				return nil
			}

			return walk(t.Elem(), prefix+"."+sliceSchema, "")
		case reflect.Map:
			// Add .{key} to path
			if prefix == "" {
				return nil // can't have map at root
			}
			return walk(t.Elem(), prefix+"."+mapSchema, "")
		default:
			// Leaf field
			result = append(result, DocModel{
				Object:        prefix,
				Type:          t.String(),
				Documentation: doc,
			})
		}
		return nil
	}

	if err := walk(modelType, "", ""); err != nil {
		return nil, err
	}
	return result, nil
}
