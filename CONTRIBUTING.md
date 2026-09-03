# Contributing to CloudAvenue SDK V2

Thank you for your interest in contributing to this project!

## Prerequisites

- **Go 1.26** or higher.
- Dependencies: `go mod tidy`

## Workflow

1. Fork and clone the repository.
2. Create a feature branch from `main`.
3. Make your changes following the [Coding Guidelines](./GUIDELINES.md).
4. Ensure tests and lint pass:

   ```sh
   go test ./...
   golangci-lint run
   ```

5. Open a PR against `main`.

## Pull Request Requirements

- PR description clearly explains the change and references related issues.
- All CI checks pass (lint, tests, license headers).
- Public symbols have GoDoc comments.
- Error handling uses wrapped context (`%w`).

For naming conventions, lint rules, and detailed coding standards, see [GUIDELINES.md](./GUIDELINES.md).

For architecture context, see [ARCHITECTURE.md](./ARCHITECTURE.md).
