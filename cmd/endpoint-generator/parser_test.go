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
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

const cerberusDocURL = "https://swagger.cloudavenue.orange-business.com/?urls.primaryName=[NGP+Cerberus]+Cloud+Avenue+API"

// writeFixture creates dir with the given file contents and returns the full path of file.
func writeFixture(t *testing.T, dir, file, content string) string {
	t.Helper()

	path := filepath.Join(dir, file)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write fixture %s: %v", path, err)
	}

	return path
}

// newFixtureTable indexes a two file package where the scanned file declares a
// constant and a sibling file declares the shared ones.
func newFixtureTable(t *testing.T) *constTable {
	t.Helper()

	dir := t.TempDir()

	target := writeFixture(t, dir, "target.go", `package iendpoints

func init() {
	const descLocal = "local description"
	const docURLShared = "from scanned file"
}
`)

	writeFixture(t, dir, "shared.go", `package iendpoints

type BackendTarget int

const (
	BackendInfrapi BackendTarget = iota + 1
	BackendVMware
)

const sharedURL = "https://example.invalid/shared/"

const (
	blockURL   = "https://example.invalid/block/"
	notAString = 42
	aliasesURL = sharedURL
)

const docURLShared = "from sibling file"
`)

	table, err := newConstTable(token.NewFileSet(), target)
	if err != nil {
		t.Fatalf("newConstTable: %v", err)
	}

	return table
}

func Test_constTable_lookup(t *testing.T) {
	table := newFixtureTable(t)

	tests := []struct {
		name  string
		input string
		want  string
		found bool
	}{
		{name: "package const block", input: "sharedURL", want: "https://example.invalid/shared/", found: true},
		{name: "const block", input: "blockURL", want: "https://example.invalid/block/", found: true},
		{name: "file scoped const", input: "descLocal", want: "local description", found: true},
		{name: "scanned file wins over sibling", input: "docURLShared", want: "from scanned file", found: true},
		{name: "non string const resolves to its name", input: "notAString", want: "notAString", found: true},
		{name: "typed backend const", input: "BackendInfrapi", want: "BackendInfrapi", found: true},
		{name: "iota continuation const", input: "BackendVMware", want: "BackendVMware", found: true},
		{name: "const alias resolves to its name", input: "aliasesURL", want: "aliasesURL", found: true},
		{name: "undeclared identifier", input: "missingConst", want: "", found: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, found := table.lookup(tc.input)
			if got != tc.want || found != tc.found {
				t.Errorf("lookup(%q) = (%q, %t), want (%q, %t)", tc.input, got, found, tc.want, tc.found)
			}
		})
	}
}

func Test_findValue(t *testing.T) {
	table := newFixtureTable(t)

	tests := []struct {
		name    string
		source  string
		want    string
		wantErr bool
	}{
		{
			name:   "const identifier",
			source: `DocumentationURL: sharedURL`,
			want:   "https://example.invalid/shared/",
		},
		{
			name:   "string literal",
			source: `DocumentationURL: "https://example.invalid/literal/"`,
			want:   "https://example.invalid/literal/",
		},
		{
			name:   "raw string literal",
			source: "DocumentationURL: `https://example.invalid/raw/`",
			want:   "https://example.invalid/raw/",
		},
		{
			name:   "selector expression",
			source: `Backend: cav.BackendInfrapi`,
			want:   "cav.BackendInfrapi",
		},
		{
			name:   "composite literal",
			source: `ResponseType: itypes.APIResponseT0{}`,
			want:   "itypes.APIResponseT0",
		},
		{
			name:   "empty anonymous struct literal",
			source: `ResponseType: struct{}{}`,
			want:   "",
		},
		{
			name:   "slice literal",
			source: `PathParams: []cav.PathParam{}`,
			want:   "[]cav.PathParam",
		},
		{
			name:   "predeclared nil",
			source: `PathParams: nil`,
			want:   "",
		},
		{
			name:   "concatenated constants",
			source: `PathTemplate: sharedURL + "/tier-0-vrfs"`,
			want:   "https://example.invalid/shared//tier-0-vrfs",
		},
		{
			name:    "unresolvable identifier",
			source:  `DocumentationURL: notDeclaredAnywhere`,
			wantErr: true,
		},
		{
			name:    "unsupported expression",
			source:  `DocumentationURL: fmt.Sprintf("%s", "url")`,
			wantErr: true,
		},
		{
			name:    "non additive binary expression",
			source:  `DocumentationURL: "a" - "b"`,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := parser.ParseExpr("cav.Endpoint{" + tc.source + "}")
			if err != nil {
				t.Fatalf("parse %q: %v", tc.source, err)
			}

			compositeLit, ok := expr.(*ast.CompositeLit)
			if !ok {
				t.Fatalf("parse %q: not a composite literal", tc.source)
			}

			kv, ok := compositeLit.Elts[0].(*ast.KeyValueExpr)
			if !ok {
				t.Fatalf("parse %q: not a key value expression", tc.source)
			}

			got, err := findValue(kv, table)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("findValue(%q) = %q, want error", tc.source, got)
				}

				return
			}

			if err != nil {
				t.Fatalf("findValue(%q): %v", tc.source, err)
			}

			if got != tc.want {
				t.Errorf("findValue(%q) = %q, want %q", tc.source, got, tc.want)
			}
		})
	}
}

func Test_findValue_nilKeyValue(t *testing.T) {
	got, err := findValue(nil, newFixtureTable(t))
	if err != nil {
		t.Fatalf("findValue(nil): %v", err)
	}

	if got != "" {
		t.Errorf("findValue(nil) = %q, want empty", got)
	}
}

// Test_scanEndpoints_resolvesSharedConstants guards the real endpoint files: the
// generator runs per file, so the documentation URL constants declared in
// helpers.go must still reach the generated comments.
func Test_scanEndpoints_resolvesSharedConstants(t *testing.T) {
	t.Parallel()

	fileAST, endpts, err := scanEndpoints(filepath.Join("..", "..", "internal", "iendpoints", "t0.go"), false)
	if err != nil {
		t.Fatalf("scanEndpoints: %v", err)
	}

	if fileAST.Name.Name != "iendpoints" {
		t.Fatalf("scanned package = %q, want iendpoints", fileAST.Name.Name)
	}

	byName := make(map[string]endpoint, len(endpts))
	for _, endpt := range endpts {
		byName[endpt.Name] = endpt
	}

	listT0, ok := byName["ListT0"]
	if !ok {
		t.Fatalf("ListT0 not extracted, got %v", byName)
	}

	if listT0.DocumentationURL != cerberusDocURL {
		t.Errorf("ListT0 DocumentationURL = %q, want %q", listT0.DocumentationURL, cerberusDocURL)
	}

	if listT0.Description != "List T0" {
		t.Errorf("ListT0 Description = %q, want %q", listT0.Description, "List T0")
	}
}
