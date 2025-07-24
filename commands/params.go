package commands

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/strcase"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/pkg/errors"
	"github.com/orange-cloudavenue/common-go/validators"
)

// validArgNameRegex regex to check that args words are lower-case or digit starting and ending with a letter.
var validArgNameRegex = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

const (
	sliceSchema = "{index}"
	mapSchema   = "{key}"
)

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

// This function take a go struct and a name that comply with ParamSpec name notation (e.g "friends.{index}.name")
func getParamType(paramType reflect.Type, name string) (reflect.Type, error) {
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
func getValuesForFieldByName(
	value reflect.Value,
	parts []string,
) (values []reflect.Value, err error) {
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
