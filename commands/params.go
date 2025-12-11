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
	"reflect"
	"unsafe"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/cav"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands/pspecs"
)

// Boilerplate to generate a tag from a paramSpec
func buildTagFromParamSpec(spec pspecs.ParamSpec) string {
	if spec == nil || len(spec.GetValidators()) == 0 {
		return ""
	}
	tag := ""
	for _, validator := range spec.GetValidators() {
		if tag != "" {
			tag += ","
		}
		tag += validator.GetKey()
	}
	return tag
}

// getCavClientFromInterface attempts to extract a field named 'c' of type cav.Client from any struct.
// This function uses reflection and unsafe to access both exported and unexported ('private') fields.
// It works for both pointer and non-pointer structs, as long as the struct contains a field named 'c'.
// Returns the cav.Client value and true if found and type assertion is successful, otherwise returns nil and false.
//
// Use cases:
//   - When you receive an interface{} and do not know the concrete struct type, but you know it embeds or contains a cav.Client as the field 'c'.
//   - To provide a generic access mechanism for packages that wrap or embed cav.Client.
//
// Limitations:
//   - Accessing unexported fields is only possible from within the same package as the struct definition (enforced by Go).
//   - Using unsafe may break in future Go versions or with changes in internal struct layout.
//   - If the struct does not have a field named 'c', or if 'c' is not of type cav.Client, (nil, false) is returned.
//
// Example usage:
//
//	var myAny interface{} = &MyStruct{c: myClient}
//	client, ok := getCavClientFromInterface(myAny)
//	if ok {
//	    // Use client (type cav.Client)
//	}
//
// Safety note:
//   - Prefer to expose a getter method for cav.Client if possible, for type safety and API clarity.
//
// Parameters:
//   - obj: An interface{} which should point to a struct or be a struct containing a field 'c'.
//
// Returns:
//   - cav.Client: The underlying cav.Client field if found and of correct type, else nil.
//   - bool: True if extraction and type assertion succeed, false otherwise.
func getCavClientFromInterface(obj interface{}) (cav.Client, bool) {
	v := reflect.ValueOf(obj)
	if !v.IsValid() {
		return nil, false
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, false
	}
	f := v.FieldByName("c")
	if !f.IsValid() {
		return nil, false
	}
	// If the field is unexported, use unsafe to allow Interface(), only from the same package
	if !f.CanInterface() {
		f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
	}
	client, ok := f.Interface().(cav.Client)
	return client, ok
}

// func (c *Command) validate(client, params interface{}) error {

// 	val := reflect.ValueOf(params)
// 	if val.Kind() == reflect.Ptr {
// 		val = val.Elem()
// 	}
// 	typ := val.Type()
// 	if typ.Kind() != reflect.Struct {
// 		return errors.New("params must be a struct or pointer to struct")
// 	}

// 	var fields []reflect.StructField

// 	for i := 0; i < typ.NumField(); i++ {
// 		field := typ.Field(i)
// 		// Search for the associated ParamSpec
// 		paramSpec := c.findParamSpecByGoField(field.Name)
// 		// Build the dynamic tag
// 		tag := buildTagFromParamSpec(paramSpec)

// 		// Copy the field and replace/add tag
// 		newField := field
// 		// Add or replace the "validate" tag
// 		if tag != "" {
// 			if newField.Tag == "" {
// 				newField.Tag = reflect.StructTag(fmt.Sprintf(`validate:"%s"`, tag))
// 			} else {
// 				newField.Tag += reflect.StructTag(fmt.Sprintf(` validate:"%s"`, tag))
// 			}
// 		}
// 		fields = append(fields, newField)
// 	}

// 	// Create the new dynamic structure
// 	dynType := reflect.StructOf(fields)
// 	dynValue := reflect.New(dynType).Elem()

// 	// Fill dynValue with the values from params (recursively)
// 	copyValuesRecursive(dynValue, val)

// 	// Validation
// 	validate := validators.New()
// 	if err := validate.Struct(dynValue.Addr().Interface()); err != nil {
// 		return errors.New(fmt.Sprintf("invalid params: %v", err.Error()))
// 	}

// 	return nil
// }

// Boilerplate: Recursive copy of values val -> dynVal
// (to be completed to handle Map/Slice/Struct/Embedded)
// func copyValuesRecursive(dynVal, val reflect.Value) {
// 	if !val.IsValid() || !dynVal.CanSet() {
// 		return
// 	}

// 	switch dynVal.Kind() {
// 	case reflect.Struct:
// 		for i := 0; i < dynVal.NumField(); i++ {
// 			origField := val.Field(i)
// 			newField := dynVal.Field(i)
// 			// Handle embedded fields (anonymous)
// 			if dynVal.Type().Field(i).Anonymous {
// 				copyValuesRecursive(newField, origField)
// 				continue
// 			}
// 			switch origField.Kind() {
// 			case reflect.Struct:
// 				copyValuesRecursive(newField, origField)
// 			case reflect.Slice:
// 				copyValuesRecursive(newField, origField)
// 			case reflect.Map:
// 				copyValuesRecursive(newField, origField)
// 			default:
// 				if newField.CanSet() && origField.IsValid() {
// 					newField.Set(origField)
// 				}
// 			}
// 		}
// 	case reflect.Slice:
// 		if val.Kind() != reflect.Slice {
// 			return
// 		}
// 		slice := reflect.MakeSlice(dynVal.Type(), val.Len(), val.Len())
// 		for i := 0; i < val.Len(); i++ {
// 			elemDst := reflect.New(dynVal.Type().Elem()).Elem()
// 			copyValuesRecursive(elemDst, val.Index(i))
// 			slice.Index(i).Set(elemDst)
// 		}
// 		dynVal.Set(slice)
// 	case reflect.Map:
// 		if val.Kind() != reflect.Map {
// 			return
// 		}
// 		mapType := dynVal.Type()
// 		newMap := reflect.MakeMapWithSize(mapType, val.Len())
// 		for _, key := range val.MapKeys() {
// 			valElem := val.MapIndex(key)
// 			newElem := reflect.New(mapType.Elem()).Elem()
// 			copyValuesRecursive(newElem, valElem)
// 			newMap.SetMapIndex(key, newElem)
// 		}
// 		dynVal.Set(newMap)
// 	default:
// 		if dynVal.CanSet() && val.IsValid() {
// 			dynVal.Set(val)
// 		}
// 	}
// }
