# Contributing to CloudAvenue SDK V2

Thank you for your interest in contributing to this project! Please read the following guidelines carefully before submitting your contribution.

---

## Prerequisites

- **Go Version**: This project requires **Go 1.24** or higher.  
  Check your version with:

  ```sh
  go version
  ```

- Install dependencies with:

  ```sh
  go mod tidy
  ```

---

## Best Practices

- **Commit Convention**:  
  Use the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format for your commit messages (e.g., `feat: add multi-tenant support`, `fix: authentication bug fix`).

- **Unit Tests**:  
  Every new feature or bug fix must be covered by unit tests.  
  Run tests with:

  ```sh
  go test ./...
  ```

- **Linting**:  
  The project uses [golangci-lint](https://golangci-lint.run/).  
  Check linting before submitting a PR:

  ```sh
  golangci-lint run
  ```

- **Respect GitHub Workflows**:  
  PRs are automatically validated (lint, tests, license, etc).  
  Make sure all checks pass before requesting a review.

- **Documentation**:  
  Document your public functions and add usage examples if necessary.

---

Thank you for following these guidelines to help ensure the quality and maintainability of the CloudAvenue SDK. We appreciate your contributions!
  
This diagram reflects the actual dependencies and composition in the [`cav`](cav) package:  

- The main `client` holds a `resty.Client`, a `consoles.Console`, and a map of initialized sub-clients.
- Each sub-client embeds a `subclient` struct, which itself holds a `resty.Client` and an `auth` credential.
- The `auth` interface is implemented by `cloudavenueCredential`.

## 2.2 The `api` Directory

The `api/` directory contains all the APIs consumed by the SDK, organized by major object groups and by version.

- **Object Groups:**  
  Each subdirectory under `api/` represents a major CloudAvenue object or resource type, such as `edgegateway`, `vdc`, `vapp`, etc.  
  This structure helps keep the codebase modular and easy to navigate.

- **Versioning:**  
  Inside each object group, APIs are further organized by version (e.g., `v1`, `v2`, etc).  
  **Important:** These versions do **not** directly reflect the upstream API versions.  
  Instead, they are used internally to allow the SDK to implement new versions when there are significant changes in the API, making it easier to manage breaking changes and maintain backward compatibility.

### Example Structure

```

api/
  edgegateway/
    v1/
      edgegateway.go
    v2/
      edgegateway.go
  vdc/
    v1/
      vdc.go
  vapp/
    v1/
      vapp.go

```

This approach allows the SDK to evolve and support multiple API versions for each object group, ensuring stability for users even as the underlying APIs change.

## 3. Contribution Best Practices

- **Commit Convention**:  
  Use the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format for your commit messages (e.g., `feat: add multi-tenant support`, `fix: authentication bug fix`).

- **Unit Tests**:  
  Every new feature or bug fix must be covered by unit tests.  
  Run tests with:

  ```sh
  go test ./...
  ```

- **Linting**:  
  The project uses [golangci-lint](https://golangci-lint.run/).  
  Check linting before submitting a PR:

  ```sh
  golangci-lint run
  ```

- **Respect GitHub Workflows**:  
  PRs are automatically validated (lint, tests, license, etc).  
  Make sure all checks pass before requesting a review.

- **Documentation**:  
  Document your public functions and add usage examples if necessary.

---

Thank you for following these guidelines to help ensure the quality and maintainability of the CloudAvenue SDK. We appreciate your contributions!
