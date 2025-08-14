//go:build ruleguard
// +build ruleguard

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package gorules

import (
	"github.com/quasilyte/go-ruleguard/dsl"
)

// apiTypes checks for common API types in the codebase.
// It ensures that the types used in API endpoints are consistent and follow best practices.
// Ensures that struct types in any `api/` directory follow these naming conventions:
//   - **API Response Types:**
//     Must be named `ApiResponse<Object>` (e.g., `ApiResponseEdgeGateway`).
//   - **API Request Body Types:**
//     Must be named `ApiRequest<Object>` (e.g., `ApiRequestEdgeGateway`).
//   - **User-facing Model Types:**
//     Must be named `Model<Object>` (e.g., `ModelEdgeGateway`).
//   - **User-supplied Parameter Types:**
//     Must be named `Params<Object>` (e.g., `ParamsEdgeGateway`).
//   - **Client Types:**
//     Must be named `Client` (exactly, for the main client struct of an API group, e.g., `type Client struct { ... }`).
//
// Regex101 https://regex101.com/r/9Mv2Ak
// If a struct type in an `api/` directory does not follow one of these conventions, the linter will report an error.
func apiResponseTypes(m dsl.Matcher) {
	m.Match(
		`type $name struct { $*_ }`,
	).
		Where(m.File().PkgPath.Matches(`api/`) && !m["name"].Text.Matches(`^(ApiResponse|ApiRequest|Model|Params)[A-Z][A-Za-z0-9]*$|^Client$`)).
		Report(`Struct type names must start with ApiResponse | ApiRequest | Model | Params (See GUIDELINE.md)`)
}

// apiFuncPrefix enforces naming conventions for exported functions in API packages.
// It ensures that exported function names in any `api/` directory start with one of the following prefixes:
//   - Get
//   - Create
//   - List
//   - Delete
//   - Update
//   - Enable
//   - Disable
//   - Add
//   - Remove
//
// If an exported function does not follow this convention, the linter will report an error.
func apiFuncPrefix(m dsl.Matcher) {
	isExported := func(v dsl.Var) bool {
		return v.Text.Matches(`^\p{Lu}`)
	}

	m.Match(`func ($*method) $name($*params) ($*output) { $*body }`).
		Where(
			m.File().PkgPath.Matches(`api/`) &&
				!m["name"].Text.Matches(`^(Get|Create|List|Delete|Update|Enable|Disable|Add|Remove)[A-Z]{1}[A-Za-z0-9_]*$`) &&
				isExported(m["name"]),
		).
		Report(`Function names must start with Get | Create | List | Delete | Update | Enable | Disable | Add | Remove (See CONTRIBUTING.md)`)
}

// boolFunctionNaming enforces naming conventions for functions that return a boolean value.
// It ensures that any function returning a bool starts with one of the following prefixes:
//   - Is
//   - is
//   - Has
//   - has
//   - match
//   - Match
//
// If a bool-returning function does not follow this convention, the linter will report an error.
func boolFunctionNaming(m dsl.Matcher) {
	m.Match(`func $name($*params) bool { $*body }`).
		Where(!m["name"].Text.Matches(`^(Is|is|Has|has|match|Match).*`)).
		Report("bool function name should start with 'Is' | 'is' | 'Has' | 'has' | 'match' | 'Match'")
}
