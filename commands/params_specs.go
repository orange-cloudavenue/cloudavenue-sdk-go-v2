package commands

import (
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-sdk-go/strcase"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/common-go/validators"
)

func (p *ParamsSpecs) validate(params any) error {
	// Decode struct
	val := reflect.ValueOf(params)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		return errors.New("params must be a struct or pointer to struct")
	}

	var fields []reflect.StructField

	// For each field in the struct
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// Find the associated ParamSpec
		paramSpec := p.findParamSpecByGoField(field.Name)
		// Build the dynamic tag
		tag := buildTagFromParamSpec(paramSpec)

		// Copy the field and replace/add tag
		newField := field
		// Add or replace the "validate" tag
		if tag != "" {
			if newField.Tag == "" {
				newField.Tag = reflect.StructTag(fmt.Sprintf(`validate:"%s"`, tag))
			} else {
				newField.Tag += reflect.StructTag(fmt.Sprintf(` validate:"%s"`, tag))
			}
		}
		fields = append(fields, newField)
	}

	// Create the new dynamic structure
	dynType := reflect.StructOf(fields)
	dynValue := reflect.New(dynType).Elem()

	// Fill dynValue with the values from params (recursively)
	copyValuesRecursive(dynValue, val)

	// Validation
	validate := validators.New()
	if err := validate.Struct(dynValue.Addr().Interface()); err != nil {
		return errors.New(fmt.Sprintf("invalid params: %v", err.Error()))
	}

	return nil
}

// Find the ParamSpec from the field name (conversion Go -> ParamSpecNameNotation)
func (p ParamsSpecs) findParamSpecByGoField(goFieldName string) *ParamsSpec {
	for _, ps := range p {
		if strcase.ToPublicGoName(ps.Name) == goFieldName {
			return &ps
		}
	}
	return nil
}
