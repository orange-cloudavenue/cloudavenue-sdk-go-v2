# CloudAvenue SDK V2 Guidelines

This document describes the coding standards, lint rules, naming conventions, and best practices for contributing to the CloudAvenue SDK V2 project.

---

## Coding Standards

- **Idiomatic Go**: Follow [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).
- **Documentation**: All exported functions, types, and packages must have clear GoDoc comments.
- **Error Handling**: Use wrapped errors with context. Prefer `fmt.Errorf("...: %w", err)`.
- **Testing**: Write table-driven tests. Use subtests for complex scenarios.
- **Concurrency**: Use sync primitives (e.g., `sync.Mutex`) carefully. Avoid global state.

---

## Lint Rules and Naming Conventions

To maintain code quality and consistency, the following lint rules are enforced in this project.  
These rules are checked automatically by custom linters and ruleguard scripts.

### API Type Naming

- **API Response Types**  
  - Format: `ApiResponse<Object>` (e.g., `ApiResponseEdgeGateway`).
  - Used for types representing API responses from CloudAvenue endpoints.
  - Location: `internal/itypes/`
  - **Note:** These types are for internal use only and must not be exposed to SDK users.

- **API Request Body Types**  
  - Format: `ApiRequest<Object>` (e.g., `ApiRequestEdgeGateway`).
  - Used for types representing request bodies sent to CloudAvenue endpoints.
  - Location: `internal/itypes/`
  - **Note:** These types are for internal use only and must not be exposed to SDK users.

- **User-facing Model Types**  
  - Format: `Model<Object>` (e.g., `ModelEdgeGateway`).
  - Used for types exposed to SDK users, representing domain models.
  - Location: `types/`

- **User-supplied Parameter Types**  
  - Format: `Params<Object>` (e.g., `ParamsEdgeGateway`).
  - Used for types representing parameters supplied by users to SDK methods.
  - Location: `types/`

- **Client Types**  
  - Name: `Client` (exactly, for the main client struct of an API group).
  - Used for the main client struct managing API group operations.
  - Location: `api/<group>/` or `cav/`

> You can visualize and debug the naming convention regex used by the linter here: [https://regex101.com/r/g8Av6t/1](https://regex101.com/r/g8Av6t/1)

### API Function Naming

- Exported functions in any `api/` directory must start with one of the following prefixes:  
  `Get`, `Create`, `List`, `Delete`, `Update`, `Enable`, `Disable`, `Add`, `Remove`
  - Example: `GetEdgeGateway`, `CreateVDC`, `ListVApps`
- Functions returning a boolean value must start with:  
  `Is`, `is`, `Has`, `has`, `Match`, `match`
  - Example: `IsEnabled`, `HasPermission`, `matchURN`

**Details of names of exported functions in `api/` directories :**

- `List` functions should return a list of objects. No error returned if the list is empty.
- `Get` functions should return a single object. Return an error if not found.
- `Create` functions create a new object and return it.
- `Update` functions update the specified object and return it.
- `Delete` functions delete the specified object.
- `Enable` functions enable the specified object. Only errors are returned.
- `Disable` functions disable the specified object. Only errors are returned.
- `Add` functions add a new object. Only errors are returned.
- `Remove` functions remove the specified object. Only errors are returned.

---

## Naming Conventions

- **Packages**: Use short, meaningful, lowercase names (e.g., `client`, `api`, `utils`).
- **Files**: Use snake_case for file names (e.g., `auth.go`, `client_options.go`).
- **Types**: Use PascalCase for exported types (e.g., `CloudAvenueClient`).
- **Functions**: Use camelCase for unexported, PascalCase for exported.
- **Variables**: Use short, descriptive names. Avoid abbreviations unless well-known.
- **Constants**: Use ALL_CAPS or PascalCase for exported constants.

---

## Pull Request Checklist

- [ ] Code is formatted (`gofmt`)
- [ ] Lint checks pass
- [ ] Unit tests added/updated
- [ ] Documentation updated
- [ ] PR description is clear and references related issues

---

For more details, see the [CONTRIBUTING.md](./CONTRIBUTING.md) and [ARCHITECTURE.md](./ARCHITECTURE.md) files.
