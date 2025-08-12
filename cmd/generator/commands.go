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
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"

	"github.com/kr/pretty"
	"github.com/spf13/cobra"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/commands"
)

func init() {
	rootCmd.AddCommand(commandCmd)
}

type cavCMDs []*cavCMD

type cavCMD struct {
	// Package is the package name where the command is defined
	Package string
	// Fields that are defined in the command definition
	Namespace, Resource, Verb                                    string
	ParamsType, ModelType                                        string
	ParamsSpecs                                                  commands.ParamsSpecs
	AutoGenerate                                                 bool
	ShortDocumentation, LongDocumentation, MarkdownDocumentation string

	// Computed fields
	CommandName string
}

var commandCmd = &cobra.Command{
	Use:   "commands",
	Short: "Generator is a command to generate commands code from the definitions.",
	Long:  `Generator is a command to generate commands code from the definitions.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		pwd, err := os.Getwd()
		if err != nil {
			return err
		}

		var (
			flagPath, _  = cmd.Flags().GetString("path")
			flagDebug, _ = cmd.Flags().GetBool("debug")
			commands     = make(cavCMDs, 0)
		)

		fset := token.NewFileSet()
		fileAst, err := parser.ParseFile(fset, pwd+"/"+flagPath, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		if flagDebug {
			ast.Print(fset, fileAst)
		}

		for _, decl := range fileAst.Decls {
			f, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			// Check if the function is init
			if f.Name.Name == "init" {
				// Generate commands for the init function

				for _, stmt := range f.Body.List {
					e, ok := stmt.(*ast.ExprStmt)
					if !ok {
						continue
					}

					c, ok := e.X.(*ast.CallExpr)
					if !ok {
						continue
					}

					// Check if the call is to Cmds.Register
					fun, ok := c.Fun.(*ast.SelectorExpr)
					if !ok || fun.Sel.Name != "Register" {
						continue
					}

					// Check if the first argument is a composite literal
					if len(c.Args) == 1 {

						cmd := &cavCMD{
							Package: fileAst.Name.Name,
						}

						decodeStruct(reflect.ValueOf(cmd), c.Args)

						// for _, arg := range c.Args {
						// 	compLit, ok := arg.(*ast.CompositeLit)
						// 	if !ok {
						// 		continue
						// 	}

						// 	cmdValue := reflect.ValueOf(cmd).Elem()

						// 	// Retrieve the Namespace, Resource and Verb from the fields
						// 	for _, field := range compLit.Elts {
						// 		kv, ok := field.(*ast.KeyValueExpr)
						// 		if !ok {
						// 			continue
						// 		}

						// 		key, ok := kv.Key.(*ast.Ident)
						// 		if !ok {
						// 			continue
						// 		}

						// 		fieldValue := cmdValue.FieldByName(key.Name)
						// 		if fieldValue.IsValid() && fieldValue.CanSet() {
						// 			v := findValue(kv)
						// 			value := reflect.ValueOf(v)
						// 			switch fieldValue.Kind() {
						// 			case reflect.String:
						// 				if value.Kind() == reflect.String {
						// 					fieldValue.SetString(value.String())
						// 				}
						// 			case reflect.Bool:
						// 				boolValue := false
						// 				switch value.Kind() {
						// 				case reflect.Bool:
						// 					boolValue = value.Bool()
						// 				case reflect.String:
						// 					// convert string to bool
						// 					bv, err := strconv.ParseBool(v)
						// 					if err != nil {
						// 						continue
						// 					}

						// 					boolValue = bv
						// 				}

						// 				fieldValue.SetBool(boolValue)
						// 			case reflect.Int, reflect.Int64:
						// 				if value.Kind() == reflect.Int || value.Kind() == reflect.Int64 {
						// 					intValue := value.Int()
						// 					fieldValue.SetInt(intValue)
						// 				}
						// 			case reflect.Slice:

						// 			}
						// 		}

						// 	}
						// }

						if cmd.AutoGenerate {
							if cmd.Resource != "" && strings.EqualFold(cmd.Namespace, cmd.Package) {
								cmd.CommandName = fmt.Sprintf("%s%s", cmd.Verb, cmd.Resource)
							} else {
								cmd.CommandName = fmt.Sprintf("%s%s%s", cmd.Verb, cmd.Namespace, cmd.Resource)
							}

							commands = append(commands, cmd)
						}
					}
				}

			}
		}

		pretty.Print("Commands", commands)

		return nil
	},
}
