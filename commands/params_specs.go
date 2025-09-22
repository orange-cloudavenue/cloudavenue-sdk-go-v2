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
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/common-go/strcase"
	"github.com/orange-cloudavenue/common-go/validators"
)

// BuildAndValidateDynamicStruct dynamically builds a struct from paramsSpecs and params,
// then validates this struct with go-playground/validator.
// It handles nested structs, slices, and maps, and applies the validation tags defined in ParamsSpecs.
func (p *ParamsSpecs) buildAndValidateDynamicStruct(params any) error {
	if params == nil {
		return fmt.Errorf("params is nil")
	}

	type fieldWithValue struct {
		Field reflect.StructField
		Value interface{}
	}

	var fields []fieldWithValue

	for _, spec := range *p {
		// Get all matching values for the spec.Name path
		values, err := GetAllValuesAtTarget(params, spec.Name)
		if err != nil {
			return fmt.Errorf("error getting values for %s: %w", spec.Name, err)
		}
		for i, val := range values {
			// Unique name for each field (ex: Users_0_Name)
			fieldName := buildFieldNameFromPath(spec.Name, i)
			tag := buildTagFromParamSpec(&spec)

			// If the field is a slice or map with {index} or {key}, add "dive" to the tag
			if strings.Contains(spec.Name, "{index}") || strings.Contains(spec.Name, "{key}") {
				tag = fmt.Sprintf("dive,%s", tag)
			}
			sf := reflect.StructField{
				Name: fieldName,
				Type: reflect.TypeOf(val),
				Tag: func() reflect.StructTag {
					if tag == "" {
						return ""
					}
					return reflect.StructTag(fmt.Sprintf(`validate:"%s"`, tag))
				}(),
			}
			fields = append(fields, fieldWithValue{Field: sf, Value: val})
		}
	}

	// Construct the dynamic struct type
	var structFields []reflect.StructField
	for _, fv := range fields {
		structFields = append(structFields, fv.Field)
	}
	dynType := reflect.StructOf(structFields)
	dynValue := reflect.New(dynType).Elem()

	for i, fv := range fields {
		dynValue.Field(i).Set(reflect.ValueOf(fv.Value))
	}

	// Validate the dynamic struct
	validate := validators.New()
	if err := validate.Struct(dynValue.Addr().Interface()); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			validationErr := errors.New("validation error")
			for _, fe := range errs {
				validationErr = fmt.Errorf("%w: param '%s' failed on the '%s' property. Allowed values %v got %v.",
					validationErr,
					fe.Namespace(),
					fe.Tag(),
					fe.Param(),
					fe.Value(),
				)
			}
			return validationErr
		}
		return errors.New(fmt.Sprintf("invalid params: %v", err.Error()))
	}

	return nil
}

// buildFieldNameFromPath generates a valid Go field name from a path (e.g., "Users.{index}.Name", 0 => "Users_0_Name")
func buildFieldNameFromPath(path string, idx int) string {
	parts := strings.Split(path, ".")
	var out []string
	for _, p := range parts {
		if p == "{index}" || p == "{key}" {
			out = append(out, fmt.Sprintf("%d", idx))
		} else {
			name := strcase.ToPublicGoName(p)
			// Force first letter uppercase (exported)
			if len(name) > 0 {
				name = strings.ToUpper(name[:1]) + name[1:]
			}
			out = append(out, name)
		}
	}
	return strings.Join(out, "_")
}

// func (p *ParamsSpecs) validate(params any) error {
// 	rV, rT, err := p.decode(params)
// 	if err != nil {
// 		return err
// 	}

// 	fields := p.buildFields(rT, "")

// 	// Create the new dynamic structure
// 	dynType := reflect.StructOf(fields)
// 	dynValue := reflect.New(dynType).Elem()

// 	// Fill dynValue with the values from params (recursively)
// 	copyValuesRecursive(dynValue, rV)

// 	// Validation
// 	validate := validators.New()
// 	if err := validate.Struct(dynValue.Addr().Interface()); err != nil {
// 		var errs validator.ValidationErrors
// 		if errors.As(err, &errs) {
// 			validationErr := errors.New("validation error")

// 			for _, fe := range errs {
// 				validationErr = fmt.Errorf("%w: param '%s' failed on the '%s' property. Allowed values %v got %v.",
// 					validationErr,
// 					fe.Namespace(), // Field with struct name
// 					fe.Tag(),
// 					fe.Param(),
// 					fe.Value(),
// 				)
// 			}

// 			return validationErr
// 		}

// 		return errors.New(fmt.Sprintf("invalid params: %v", err.Error()))
// 	}

// 	return nil
// }

// func (p *ParamsSpecs) buildFields(rT reflect.Type, prefix string) []reflect.StructField {
// 	var fields []reflect.StructField
// 	for i := 0; i < rT.NumField(); i++ {
// 		field := rT.Field(i)
// 		fieldName := field.Name
// 		paramSpec := p.findParamSpecByGoField(prefix + fieldName)
// 		tag := buildTagFromParamSpec(paramSpec)

// 		newField := field

// 		// If the field is a struct or a slice of struct, we need to build its fields recursively
// 		if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
// 			subFields := p.buildFields(field.Type.Elem(), prefix+fieldName+".{index}.")
// 			tag = "dive"
// 			newField.Type = reflect.SliceOf(reflect.StructOf(subFields))
// 		} else if field.Type.Kind() == reflect.Struct && field.Anonymous == false {
// 			subFields := p.buildFields(field.Type, prefix+fieldName+".")
// 			newField.Type = reflect.StructOf(subFields)
// 		}

// 		if tag != "" {
// 			if strings.Contains(tag, "{index}") || strings.Contains(tag, "{key}") {
// 				tag = fmt.Sprintf("dive,%s", tag)
// 			}
// 			if newField.Tag == "" {
// 				newField.Tag = reflect.StructTag(fmt.Sprintf(`validate:"%s"`, tag))
// 			} else {
// 				newField.Tag += reflect.StructTag(fmt.Sprintf(` validate:"%s"`, tag))
// 			}
// 		}
// 		fields = append(fields, newField)
// 	}

// 	return fields
// }

// // Find the ParamSpec from the field name (conversion Go -> ParamSpecNameNotation)
// func (p ParamsSpecs) findParamSpecByGoField(goFieldName string) *ParamsSpec {
// 	for _, ps := range p {
// 		if strcase.ToPublicGoName(ps.Name) == strcase.ToPublicGoName(goFieldName) {
// 			return &ps
// 		}
// 	}
// 	return nil
// }

// func (p *ParamsSpecs) decode(params any) (reflect.Value, reflect.Type, error) {
// 	// Decode struct
// 	val := reflect.ValueOf(params)
// 	if val.Kind() == reflect.Ptr {
// 		val = val.Elem()
// 	}
// 	typ := val.Type()
// 	if typ.Kind() != reflect.Struct {
// 		return reflect.Value{}, nil, errors.New("params must be a struct or pointer to struct")
// 	}
// 	return val, typ, nil
// }
