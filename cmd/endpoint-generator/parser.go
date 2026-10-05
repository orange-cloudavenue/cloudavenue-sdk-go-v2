/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
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
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

// constTable resolves identifiers used as composite literal values to their
// string content. Constants declared by the scanned file shadow identically
// named ones coming from its sibling files.
type constTable struct {
	fset   *token.FileSet
	values map[string]string
}

func scanEndpoints(path string, debug bool) (*ast.File, endpoints, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	fullPath := filepath.Join(pwd, path)

	table, err := newConstTable(token.NewFileSet(), fullPath)
	if err != nil {
		return nil, nil, err
	}

	fileAST, err := parser.ParseFile(table.fset, fullPath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	if debug {
		ast.Print(table.fset, fileAST)
	}

	endpts, err := extractEndpoints(fileAST, table)
	if err != nil {
		return nil, nil, err
	}

	return fileAST, endpts, nil
}

// newConstTable indexes every string constant reachable from the file at path.
// The generator runs once per file while shared constants are declared in a
// sibling file of the same package, so the whole directory is scanned.
func newConstTable(fset *token.FileSet, path string) (*constTable, error) {
	table := &constTable{fset: fset, values: make(map[string]string)}

	if err := table.addFile(path); err != nil {
		return nil, err
	}

	siblings, err := filepath.Glob(filepath.Join(filepath.Dir(path), "*.go"))
	if err != nil {
		return nil, err
	}

	for _, sibling := range siblings {
		if sibling == path || strings.HasSuffix(sibling, "_test.go") {
			continue
		}

		if err := table.addFile(sibling); err != nil {
			return nil, err
		}
	}

	return table, nil
}

func (t *constTable) addFile(path string) error {
	fileAST, err := parser.ParseFile(t.fset, path, nil, parser.SkipObjectResolution)
	if err != nil {
		return err
	}

	// Constants may sit at package level or inside an init body, both are
	// collected so file scoped declarations are found wherever they live.
	ast.Inspect(fileAST, func(node ast.Node) bool {
		if decl, ok := node.(*ast.GenDecl); ok && decl.Tok == token.CONST {
			t.addDecl(decl)
		}

		return true
	})

	return nil
}

func (t *constTable) addDecl(decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		for i, name := range valueSpec.Names {
			if _, exists := t.values[name.Name]; exists {
				continue
			}

			// A string constant resolves to its content. Every other constant
			// resolves to its own name, mirroring how a qualified reference
			// renders as pkg.Name. That covers typed constants such as
			// BackendInfrapi, iota continuations that carry no value of their
			// own, and aliases, none of which carry renderable text.
			if i < len(valueSpec.Values) {
				if lit, ok := valueSpec.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
					t.values[name.Name] = cleanQuote(lit.Value)
					continue
				}
			}

			t.values[name.Name] = name.Name
		}
	}
}

// lookup returns the string content behind name. Identifiers without a string
// constant report false so callers never render a silently empty value.
func (t *constTable) lookup(name string) (string, bool) {
	value, ok := t.values[name]
	return value, ok
}

func (t *constTable) position(pos token.Pos) string {
	return t.fset.Position(pos).String()
}

func extractEndpoints(fileAST *ast.File, table *constTable) (endpoints, error) {
	endpts := make(endpoints, 0)

	for _, decl := range fileAST.Decls {
		f, ok := decl.(*ast.FuncDecl)
		if !ok || f.Name.Name != "init" {
			continue
		}

		for _, stmt := range f.Body.List {
			endpt, ok, err := extractRegisteredEndpoint(stmt, table)
			if err != nil {
				return nil, err
			}

			if ok {
				endpts = append(endpts, endpt)
			}
		}
	}

	return endpts, nil
}

func extractRegisteredEndpoint(stmt ast.Stmt, table *constTable) (endpoint, bool, error) {
	e, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return endpoint{}, false, nil
	}

	c, ok := e.X.(*ast.CallExpr)
	if !ok {
		return endpoint{}, false, nil
	}

	fun, ok := c.Fun.(*ast.SelectorExpr)
	if !ok || fun.Sel.Name != "Register" {
		return endpoint{}, false, nil
	}

	compositeLit, ok := fun.X.(*ast.CompositeLit)
	if !ok {
		return endpoint{}, false, nil
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
		if !fieldValue.IsValid() || !fieldValue.CanSet() {
			continue
		}

		resolved, err := findValue(kv, table)
		if err != nil {
			return endpoint{}, false, err
		}

		fieldValue.SetString(resolved)
	}

	return endpt, true, nil
}

func cleanQuote(s string) string {
	if unquoted, err := strconv.Unquote(s); err == nil {
		return unquoted
	}

	return s
}

// findValue renders a composite literal value as a plain string, resolving
// identifiers through table so a shared constant yields its literal content.
func findValue(kv *ast.KeyValueExpr, table *constTable) (string, error) {
	if kv == nil {
		return "", nil
	}

	return resolveExpr(kv.Value, table)
}

func resolveExpr(expr ast.Expr, table *constTable) (string, error) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		return cleanQuote(v.Value), nil
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", v.X, v.Sel.Name), nil
	case *ast.Ident:
		return resolveIdent(v, table)
	case *ast.BinaryExpr:
		return resolveConcat(v, table)
	case *ast.CompositeLit:
		return typeName(v.Type), nil
	default:
		return "", fmt.Errorf("%s: unsupported endpoint value expression %T", table.position(expr.Pos()), expr)
	}
}

// resolveIdent reads a bare identifier as a constant, falling back to an empty
// value for the predeclared literals that carry no text.
func resolveIdent(expr *ast.Ident, table *constTable) (string, error) {
	switch expr.Name {
	case "nil", "true", "false", "iota":
		return "", nil
	default:
	}

	value, ok := table.lookup(expr.Name)
	if !ok {
		return "", fmt.Errorf("%s: unresolved identifier %q used as an endpoint value", table.position(expr.Pos()), expr.Name)
	}

	return value, nil
}

// resolveConcat folds a constant string concatenation such as pathBase + "/sub".
func resolveConcat(expr *ast.BinaryExpr, table *constTable) (string, error) {
	if expr.Op != token.ADD {
		return "", fmt.Errorf("%s: unsupported operator %q in endpoint value", table.position(expr.Pos()), expr.Op)
	}

	left, err := resolveExpr(expr.X, table)
	if err != nil {
		return "", err
	}

	right, err := resolveExpr(expr.Y, table)
	if err != nil {
		return "", err
	}

	return left + right, nil
}

// typeName names the type instantiated by a composite literal value. Anonymous
// types such as struct{} or map[...] carry no name and yield an empty value.
func typeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", t.X, t.Sel.Name)
	case *ast.ArrayType:
		return "[]" + typeName(t.Elt)
	default:
		return ""
	}
}
