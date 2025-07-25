package commands

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"regexp"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/strcase"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/common-go/validators"
)

// validArgNameRegex regex to check that args words are lower-case or digit starting and ending with a letter.
var validArgNameRegex = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

func (c *Command) validate(ctx context.Context, params any) error {
	// 1. Reflect params.to extract fields and their types.
	paramType := reflect.TypeOf(params)
	if paramType.Kind() == reflect.Ptr {
		paramType = paramType.Elem()
	}

	sfs := []reflect.StructField{}

	// 2. Recreate the struct type with all fields and set the tags.
	for _, p := range c.ParamsSpecs {
		tag := ""
		for _, validator := range p.Validators {
			if tag != "" {
				tag += ","
			}
			tag += validator.GetKey()
		}

		f, ok := paramType.FieldByName(strcase.ToPublicGoName(p.Name))
		if !ok {
			continue
		}
		sf := reflect.StructField{
			Name: f.Name,
			Type: f.Type,
			Tag:  reflect.StructTag(fmt.Sprintf(`validate:"%s"`, tag)),
		}

		sfs = append(sfs, sf)
	}

	// 3. Create a new struct type with the fields and tags.
	newType := reflect.StructOf(sfs)
	// 4. Create a new value of the new type.
	paramsValue := reflect.New(newType)

	// 5. Set the values of the new value.
	for _, p := range c.ParamsSpecs {
		fieldName := strcase.ToPublicGoName(p.Name)
		fieldValues, err := getValuesForFieldByName(
			reflect.ValueOf(params),
			strings.Split(fieldName, "."),
		)
		if err != nil {
			slog.Error("could not validate arg value",
				slog.String("arg", p.Name),
				slog.String("fieldName", fieldName),
				slog.Any("error", err),
			)
			// logger.Infof(
			// 	"could not validate arg value for '%v': invalid fieldName: %v: %v",
			// 	p.Name,
			// 	fieldName,
			// 	err.Error(),
			// )
			continue
		}

		// 6. inject the value contained in fieldValues into paramsValue.
		for _, v := range fieldValues {
			if !v.IsValid() {
				continue
			}

			// If the value is a pointer, we need to set the value of the pointer.
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}

			// If the value is a slice or a map, we need to set the value of each element.
			if v.Kind() == reflect.Slice || v.Kind() == reflect.Map {
				for _, key := range v.MapKeys() {
					paramsValue.Elem().SetMapIndex(key, v.MapIndex(key))
				}
			} else {
				paramsValue.Elem().FieldByName(fieldName).Set(v)
			}
		}
	}

	// 7. Validate the new value.
	if err := validators.New().Struct(paramsValue.Interface()); err != nil {
		return errors.New(fmt.Sprintf("invalid params: %v", err.Error()))
	}

	return nil
}
