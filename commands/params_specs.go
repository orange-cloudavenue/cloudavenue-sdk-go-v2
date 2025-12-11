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
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/kr/pretty"
	dynamicstruct "github.com/ompluscator/dynamic-struct"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands/pspecs"
	"github.com/orange-cloudavenue/common-go/strcase"
	"github.com/orange-cloudavenue/common-go/validators"
)

// BuildAndValidateDynamicStruct dynamically builds a struct from paramsSpecs and params,
// then validates this struct with go-playground/validator.
// It handles nested structs, slices, and maps, and applies the validation tags defined in ParamsSpecs.
func buildAndValidateDynamicStruct(paramsDef pspecs.Params, params any) error {
	if params == nil {
		return fmt.Errorf("params is nil")
	}

	// Create a dynamic struct builder. The objective is to build a struct with the same shape as params, but with validation tags.
	// We will then copy the values from params to this new struct and validate it.
	builder := dynamicstruct.NewStruct()
	buildedStruct, err := buildDynamicStruct(paramsDef)
	if err != nil {
		return err
	}

	builder.Build()

	data, err := json.Marshal(params)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &buildedStruct)
	if err != nil {
		return err
	}

	pretty.Print(buildedStruct)
	// Validate the dynamic struct
	validate := validators.New()
	if err := validate.Struct(buildedStruct); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			validationErr := errors.New("validation error")
			for _, fe := range errs {
				validationErr = fmt.Errorf("%w: param '%s'(%s) failed on the '%s' property. Allowed values '%v' got '%v'. (Actual Tag: '%s')",
					validationErr,
					fe.Namespace(),
					fe.Kind(),
					fe.Tag(),
					fe.Param(),
					fe.Value(),
					fe.ActualTag(),
				)
			}
			return validationErr
		}
		return fmt.Errorf("invalid params: %v", err.Error())
	}

	return nil
}

func buildDynamicStruct(paramsDef []pspecs.ParamSpec) (buildedStruct any, err error) {
	// Create a dynamic struct builder.
	builder := dynamicstruct.NewStruct()

	// Iterate over the fields of paramsDef to define the fields of the dynamic struct.
	for _, paramSpec := range paramsDef {
		switch x := paramSpec.(type) {
		case pspecs.ParamSpecNested:
			// Recursively build the nested struct
			nestedStruct, err := buildDynamicStruct(x.GetItemsSpec())
			if err != nil {
				return nil, err
			}
			builder.AddField(
				strcase.ToPublicGoName(x.GetName()),
				reflect.MakeSlice(
					reflect.SliceOf(
						reflect.TypeOf(
							nestedStruct,
						),
					), 0, 0).Interface(),
				buildValidatorTag(paramSpec),
			)
		default:
			// Define the field in the dynamic struct.
			builder.AddField(
				strcase.ToPublicGoName(paramSpec.GetName()),
				paramSpec.GetType().Interface(),
				buildValidatorTag(paramSpec),
			)
		}
	}

	return builder.Build().New(), nil
}

func buildValidatorTag(pS pspecs.ParamSpec) string {
	var tags []string
	if pS.IsRequired() {
		tags = append(tags, "required")
	}

	if pS.GetType().Kind() == reflect.Slice {
		tags = append(tags, "dive")
	}

	for _, v := range pS.GetValidators() {
		if v.GetKey() != "" {
			tags = append(tags, v.GetKey())
		}
	}

	if len(tags) == 0 {
		return ""
	}

	return fmt.Sprintf(`validate:"%s"`, strings.Join(tags, ","))
}
