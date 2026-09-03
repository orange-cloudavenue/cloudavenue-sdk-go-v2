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
)

func scanEndpoints(path string, debug bool) (*ast.File, endpoints, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	fset := token.NewFileSet()
	fileAST, err := parser.ParseFile(fset, pwd+"/"+path, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	if debug {
		ast.Print(fset, fileAST)
	}

	return fileAST, extractEndpoints(fileAST), nil
}

func extractEndpoints(fileAST *ast.File) endpoints {
	endpts := make(endpoints, 0)

	for _, decl := range fileAST.Decls {
		f, ok := decl.(*ast.FuncDecl)
		if !ok || f.Name.Name != "init" {
			continue
		}

		for _, stmt := range f.Body.List {
			endpoint, ok := extractRegisteredEndpoint(stmt)
			if ok {
				endpts = append(endpts, endpoint)
			}
		}
	}

	return endpts
}

func extractRegisteredEndpoint(stmt ast.Stmt) (endpoint, bool) {
	e, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return endpoint{}, false
	}

	c, ok := e.X.(*ast.CallExpr)
	if !ok {
		return endpoint{}, false
	}

	fun, ok := c.Fun.(*ast.SelectorExpr)
	if !ok || fun.Sel.Name != "Register" {
		return endpoint{}, false
	}

	compositeLit, ok := fun.X.(*ast.CompositeLit)
	if !ok {
		return endpoint{}, false
	}

	endpt := endpoint{}
	value := reflect.ValueOf(&endpt).Elem()

	for _, arg := range compositeLit.Elts {
		kv, ok := arg.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}

		fieldValue := value.FieldByName(key.Name)
		if fieldValue.IsValid() && fieldValue.CanSet() {
			fieldValue.SetString(findValue(kv))
		}
	}

	return endpt, true
}

func cleanQuote(s string) string {
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
		return cleanQuote(v.Value)
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", v.X, v.Sel.Name)
	case *ast.CompositeLit:
		kvc, ok := v.Type.(*ast.Ident)
		if !ok {
			return ""
		}
		return kvc.Name
	default:
		return ""
	}
}
