//go:build ruleguard
// +build ruleguard

package gorules

import (
	"github.com/quasilyte/go-ruleguard/dsl"
)

// func apiResponseToModel(m dsl.Matcher) {
// 	isExported := func(v dsl.Var) bool {
// 		return v.Text.Matches(`^\p{Lu}`)
// 	}

// 	m.Match(
// 		`func ($object) $funcName($*params) $*outputname { $*_ }`,
// 	).
// 		Where(
// 			m.File().PkgPath.Matches(`api/`) &&
// 				m["object"].Text.Matches(`^apiResponse[A-Z][A-Za-z0-9]*$`) &&
// 				m["funcName"].Text.Matches(`^ToModel$`) &&
// 				isExported(m["funcName"]),
// 		).
// 		Report(`Disallow exported ToModel functions in API response types. Use a private function instead.`).
// 		Suggest("func (apiResponse$name) toModel($*params) $outputname { $*_ }")
// }

// apiTypes checks for common API types in the codebase.
// It ensures that the types used in API endpoints are consistent and follow best practices.
// Ensures that struct types in any `api/` directory follow these naming conventions:
//   - **API Response Types:**
//     Must be named `apiResponse<Object>` (e.g., `apiResponseEdgeGateway`).
//   - **API Request Body Types:**
//     Must be named `apiRequest<Object>` (e.g., `apiRequestEdgeGateway`).
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
		Where(m.File().PkgPath.Matches(`api/`) && !m["name"].Text.Matches(`^(apiResponse|apiRequest|Model|Params)[A-Z][A-Za-z0-9]*$|^Client$`)).
		Report(`Struct type names must start with apiResponse | apiRequest | Model | Params (See CONTRIBUTING.md)`)
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
//
// If an exported function does not follow this convention, the linter will report an error.
func apiFuncPrefix(m dsl.Matcher) {
	isExported := func(v dsl.Var) bool {
		return v.Text.Matches(`^\p{Lu}`)
	}

	m.Match(`func ($*method) $name($*params) ($*output) { $*body }`).
		Where(
			m.File().PkgPath.Matches(`api/`) &&
				!m["name"].Text.Matches(`^(Get|Create|List|Delete|Update|Enable|Disable)[A-Z]{1}[A-Za-z0-9_]*$`) &&
				isExported(m["name"]),
		).
		Report(`Function names must start with Get | Create | List | Delete | Update | Enable | Disable (See CONTRIBUTING.md)`)
}
