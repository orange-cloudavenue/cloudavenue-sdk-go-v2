/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package main

import (
	"fmt"
	"go/ast"
	"reflect"
	"strconv"
	"strings"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
)

func clean(s string) string {
	deniedStrings := []string{"nil"}

	for _, denied := range deniedStrings {
		if strings.Contains(s, denied) {
			return ""
		}
	}

	if len(s) < 2 {
		return s
	}
	if s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

func findValue(kv *ast.KeyValueExpr) string {
	if kv == nil {
		return ""
	}

	switch v := kv.Value.(type) {
	case *ast.BasicLit:
		return clean(v.Value)
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", v.X, v.Sel.Name)
	case *ast.Ident:
		return clean(v.Name)
	case *ast.CompositeLit:
		kvc, ok := v.Type.(*ast.Ident)
		if !ok {
			fmt.Println("Could not find type for composite literal:", v.Type)
			return ""
		}
		return clean(kvc.Name)
	default:
		return ""
	}
}

func decodeStruct(structDest reflect.Value, args []ast.Expr) {
	for _, arg := range args {
		compLit, ok := arg.(*ast.CompositeLit)
		if !ok {
			continue
		}

		structValue := structDest.Elem()

		// Retrieve the Namespace, Resource and Verb from the fields
		for _, field := range compLit.Elts {
			kv, ok := field.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			key, ok := kv.Key.(*ast.Ident)
			if !ok {
				continue
			}

			fieldValue := structValue.FieldByName(key.Name)
			if fieldValue.IsValid() && fieldValue.CanSet() {
				v := findValue(kv)
				value := reflect.ValueOf(v)
				switch fieldValue.Kind() {
				case reflect.String:
					if value.Kind() == reflect.String {
						fieldValue.SetString(value.String())
					}
				case reflect.Bool:
					boolValue := false
					switch value.Kind() {
					case reflect.Bool:
						boolValue = value.Bool()
					case reflect.String:
						// convert string to bool
						bv, err := strconv.ParseBool(v)
						if err != nil {
							continue
						}

						boolValue = bv
					}

					fieldValue.SetBool(boolValue)
				case reflect.Int, reflect.Int64:
					if value.Kind() == reflect.Int || value.Kind() == reflect.Int64 {
						intValue := value.Int()
						fieldValue.SetInt(intValue)
					}
				case reflect.Slice:
					switch key.Name {
					case "ParamsSpecs":
						// Special case for ParamsSpecs, which is a slice of commands.ParamsSpec
						var specs []commands.ParamsSpec
						for _, elem := range kv.Value.(*ast.CompositeLit).Elts {
							spec := commands.ParamsSpec{}
							decodeStruct(reflect.ValueOf(&spec), []ast.Expr{elem})
							specs = append(specs, spec)
						}
						fieldValue.Set(reflect.ValueOf(specs))
						// default:
						// 	// For other slices, we can set them directly
						// 	slice := reflect.MakeSlice(fieldValue.Type(), len(v.Elts), len(v.Elts))
						// 	for i, elem := range v.Elts {
						// 		elemValue := reflect.ValueOf(findValue(elem.(*ast.KeyValueExpr)))
						// 		if elemValue.Kind() == fieldValue.Type().Elem().Kind() {
						// 			slice.Index(i).Set(elemValue)
						// 		}
						// 	}
						// 	fieldValue.Set(slice)
					}

				}
			}

		}
	}
}
