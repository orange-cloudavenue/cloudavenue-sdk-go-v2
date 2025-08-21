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
	"embed"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/kr/pretty"
)

type cmds []*cmd

type cmd struct {
	Package                    string
	Namespace, Resource, Verb  string
	AutoGenerateCustomFuncName string
	ParamsType, ModelType      string
	AutoGenerate               bool
	LongDocumentation          string
	CommandName                string
}

//go:embed generator.tmpl
var tmplFile embed.FS

func main() {

	var (
		flagPath  = flag.String("path", "", "The path to the file to generate commands from")
		flagDebug = flag.Bool("debug", false, "Enable debug mode")
	)

	flag.Parse()

	if *flagPath == "" {
		flag.Usage()
		return
	}

	// print pwd
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	commands := make(cmds, 0)

	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, pwd+"/"+*flagPath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	if *flagDebug {
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

					cmd := &cmd{
						Package: fileAst.Name.Name,
					}

					for _, arg := range c.Args {
						compLit, ok := arg.(*ast.CompositeLit)
						if !ok {
							continue
						}

						cmdValue := reflect.ValueOf(cmd).Elem()

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

							fieldValue := cmdValue.FieldByName(key.Name)
							if fieldValue.IsValid() && fieldValue.CanSet() {
								v := findValue(kv)
								if v == "cav.Job" {
									v = ""
								}
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
								}
							}

						}
					}

					if cmd.AutoGenerate {
						if cmd.Resource == "" && cmd.Verb == "" {
							// Command is a top-level command used to generate documentation
							// Ignore it
							continue
						}
						if cmd.AutoGenerateCustomFuncName != "" {
							cmd.CommandName = cmd.AutoGenerateCustomFuncName
						} else if cmd.Resource != "" && strings.EqualFold(cmd.Namespace, cmd.Package) {
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

	if *flagDebug {
		pretty.Print("Commands", commands)
	}

	var commandTmpl = struct {
		PackageName string
		Commands    cmds
	}{
		PackageName: fileAst.Name.Name,
		Commands:    commands,
	}

	tmpl, err := template.ParseFS(tmplFile, "generator.tmpl")
	if err != nil {
		panic(err)
	}

	// flagPath == /api/edgegateway/v1/edgegateway_commands.go
	// We want to generate the file in /api/edgegateway/v1/zz_<namespace>.go
	split := strings.Split(pwd+"/"+*flagPath, "/")
	nameExtracted := split[len(split)-1]
	nameExtracted = nameExtracted[:len(nameExtracted)-len("_commands.go")]

	outputPath := strings.Join(split[:len(split)-1], "/") + "/zz_" + nameExtracted + ".go"

	log.Default().Print("Path to output file: ", outputPath)

	// Create io.Writer to write the output to file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, commandTmpl)
	if err != nil {
		panic(err)
	}

}

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
		switch v.Type.(type) {
		case *ast.Ident:
			return clean(v.Type.(*ast.Ident).Name)
		case *ast.SelectorExpr:
			return fmt.Sprintf("%s.%s", v.Type.(*ast.SelectorExpr).X, v.Type.(*ast.SelectorExpr).Sel.Name)
		default:
			fmt.Println("Could not find type for composite literal:", v.Type)
			return ""
		}
	default:
		return ""
	}
}
