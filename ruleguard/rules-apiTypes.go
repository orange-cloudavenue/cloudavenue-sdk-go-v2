//go:build ruleguard
// +build ruleguard

package gorules

import (
	"github.com/quasilyte/go-ruleguard/dsl"
)

func apiResponseTypes(m dsl.Matcher) {
	m.Match(`type $name struct { $*_ }`).
		Where(m.File().PkgPath.Matches(`api/`) && !m["name"].Text.Matches(`^(apiResponse|apiRequest|Model|Params)[A-Z][A-Za-z0-9]*$|^Client$`)).
		Report(`Nom de la structure non conforme : $name`)
}

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
// If a struct type in an `api/` directory does not follow one of these conventions, the linter will report an error.

// // Regex101 https://regex101.com/r/9Mv2Ak
// var apiTypesRe = regexp.MustCompile(`(^(apiResponse|apiRequest|Model|Params)[A-Z][A-Za-z0-9]*$|^Client$)`)

// func regexMatcher(ctx *dsl.VarFilterContext) bool {
// 	return apiTypesRe.MatchString(ctx.GetInterface(ctx.GetType("Name")).String())
// }
