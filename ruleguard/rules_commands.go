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

// **Details of names of exported functions in `api/` directories :**

// - `List` functions should return a list of objects. No error returned if the list is empty.
// - `Get` functions should return a single object. Return an error if not found.
// - `Create` functions create a new object and return it.
// - `Update` functions update the specified object and return it.
// - `Delete` functions delete the specified object.
// - `Enable` functions enable the specified object. Only errors are returned.
// - `Disable` functions disable the specified object. Only errors are returned.
// - `Add` functions add a new object. Only errors are returned.
// - `Remove` functions remove the specified object. Only errors are returned.
func commandsModelTypes(m dsl.Matcher) {
	// * Require ModelType
	m.Match(`commands.Command{ $*_, Namespace: $ns, Verb: $verb, $*_ }`).
		Where(
			(m["verb"].Text == `"Create"` ||
				m["verb"].Text == `"Update"` ||
				m["verb"].Text == `"Get"` ||
				m["verb"].Text == `"List"`) &&
				!m["$$"].Text.Matches(`\bModelType\s*:`)).
		Report(`Command with Verb:$verb must include ModelType`)

	// * For other verbs, ModelType must NOT be present
	m.Match(`commands.Command{ $*_, Verb: $verb, $*_ }`).
		Where(
			(m["verb"].Text == `"Delete"` ||
				m["verb"].Text == `"Enable"` ||
				m["verb"].Text == `"Disable"` ||
				m["verb"].Text == `"Add"` ||
				m["verb"].Text == `"Remove"`) &&
				m["$$"].Text.Matches(`\bModelType\s*:`)).
		Report(`Command with Verb:$verb must NOT include ModelType`)
}

// ParamsSpecs should have a Name (only lowercase) and a Description
func commandsParamsSpecs(m dsl.Matcher) {

	// * Name
	// Ensure the ParamsSpec struct contains a Name field.
	m.Match(`commands.ParamsSpec{ $*_ }`).
		Where(!m["$$"].Text.Matches(`\bName\s*:`)).
		Report(`ParamsSpec.Name is required`)
		// Ensure the Name field value matches the required format.
		// The Name must be lowercase, can contain letters a-z, digits, underscores, and optionally {index} or {key} for templating.
	// https://regex101.com/r/YcXFnJ/1
	m.Match(`commands.ParamsSpec{ $*_, Name: $name, $*_ }`).
		Where(!m["name"].Text.Matches(`^"(?:[a-z_][a-z0-9_]*|\{(?:index|key)\})(?:\.(?:[a-z_][a-z0-9_]*|\{(?:index|key)\}))*"$`)).
		At(m["name"]).
		Report(`ParamsSpec.Name ($name) must be lowercase and contain only letters a-z`)

	// * Description
	m.Match(`commands.ParamsSpec{ $*_, Name: $name, $*_ }`).
		Where(!m["$$"].Text.Matches(`(?m)^[^\n/]*\bDescription\s*:`)).
		Report(`ParamsSpec.Description is required`)

		// Ensure the Description field is not empty or only whitespace.
	m.Match(`commands.ParamsSpec{ $*_, Description: $d, $*_ }`).
		Where(m["d"].Text.Matches(`^"\s*"$`)).
		At(m["d"]).
		Report(`ParamsSpec.Description must not be empty`)
}

// commands.Command.Verb
func commandsAllowedVerbs(m dsl.Matcher) {
	// Check if Verb are defined
	m.Match(`commands.Command{ $*_, Namespace: $ns, RunnerFunc: $runner, $*_ }`).
		Where(!m["$$"].Text.Matches(`(?m)^[^\n/]*\bVerb\s*:`)).
		Report(`commands.Command.Verb is required`)

	// Check if Verb is one of the allowed ones
	m.Match(`commands.Command{ $*_, Verb: $verb, RunnerFunc: $runner, $*_ }`).
		Where(!m["verb"].Text.Matches(`^(?:"Create"|"Update"|"Get"|"List"|"Delete"|"Enable"|"Disable"|"Add"|"Remove")$`)).
		At(m["verb"]).
		Report(`commands.Command.Verb ($verb) must be one of "Create", "Update", "Get", "List", "Delete", "Enable", "Disable", "Add", "Remove"`)
}

// commands.Command.ShortDocumentation
func commandsShortDocumentation(m dsl.Matcher) {
	// Check if ShortDocumentation is defined
	m.Match(`commands.Command{ $*_, Namespace: $ns, Verb: $verb, $*_ }`).
		Where(!m["$$"].Text.Matches(`(?m)^[^\n/]*\bShortDocumentation\s*:`)).
		Report(`commands.Command.ShortDocumentation is required`)

	// Check if ShortDocumentation is not empty
	m.Match(`commands.Command{ $*_, ShortDocumentation: $doc, $*_ }`).
		Where(m["doc"].Text.Matches(`^"\s*"$`)).
		At(m["doc"]).
		Report(`commands.Command.ShortDocumentation must not be empty`)
}

// commands.Command.LongDocumentation
func commandsLongDocumentation(m dsl.Matcher) {
	// Check if LongDocumentation is defined
	m.Match(`commands.Command{ $*_, Namespace: $ns, Verb: $verb, $*_ }`).
		Where(!m["$$"].Text.Matches(`(?m)^[^\n/]*\bLongDocumentation\s*:`)).
		Report(`commands.Command.LongDocumentation is required`)

	// Check if LongDocumentation is not empty
	m.Match(`commands.Command{ $*_, LongDocumentation: $doc, $*_ }`).
		Where(m["doc"].Text.Matches(`^"\s*"$`)).
		At(m["doc"]).
		Report(`commands.Command.LongDocumentation must not be empty`)
}

// commands.Command.Namespace
func commandsNamespace(m dsl.Matcher) {
	// Check if Namespace is defined
	m.Match(`commands.Command{ $*_ }`).
		Where(!m["$$"].Text.Matches(`(?m)^[^\n/]*\bNamespace\s*:`)).
		Report(`commands.Command.namespace is required`)

	// Check if Namespace is not empty
	m.Match(`commands.Command{ $*_, Namespace: $ns, $*_ }`).
		Where(m["ns"].Text.Matches(`^"\s*"$`)).
		At(m["ns"]).
		Report(`commands.Command.Namespace must not be empty`)
}

// commands.Command.RunnerFunc
func commandsRunnerFunc(m dsl.Matcher) {
	// Check if RunnerFunc is defined
	m.Match(`commands.Command{ $*_, Namespace: $ns, Verb: $verb, $*_ }`).
		Where(!m["$$"].Text.Matches(`(?m)^[^\n/]*\bRunnerFunc\s*:`)).
		Report(`commands.Command.RunnerFunc is required`)

	// Check if RunnerFunc is not nil
	m.Match(`commands.Command{ $*_, RunnerFunc: $rf, $*_ }`).
		Where(m["rf"].Text == "nil").
		At(m["rf"]).
		Report(`commands.Command.RunnerFunc must not be nil`)
}
