
# CloudAvenue SDK V2 Architecture

This document provides a detailed overview of the project architecture, focusing on the main folders (from most to least important), and explains the design and role of Endpoints and Commands in `/api/`.

---

## 1. Main Folders (Ordered by Importance)

### 1.1 `cav/`

**Core SDK logic.**

- Main entry point for the SDK: manages authentication, configuration, and session lifecycle.
- Hosts the main `Client` struct, which exposes sub-clients for each CloudAvenue service (e.g., VMware, Cerberus, Netbackup).
- Handles endpoint discovery, service routing, and credential management.
- Contains sub-client implementations and shared logic for all services.

### 1.2 `api/`

**CloudAvenue API implementations.**

- Organized by functional domain (e.g., `edgegateway/`, `vdc/`, `vdcgroup/`).
- Each subdirectory exposes high-level methods for interacting with a specific API group.
- Defines API types (requests, responses, models, params) and main `Client` for each group.
- Implements endpoint logic and command registration for each resource.

### 1.3 `types/`

**User-facing domain models and parameter types.**

- Contains all exported types exposed to SDK users (e.g., `ModelEdgeGateway`, `ParamsEdgeGateway`).
- Used for input/output in public SDK methods.

### 1.4 `internal/`

**Internal helpers and implementation details.**

- Contains non-exported logic, utilities, and shared code not exposed to SDK users.
- Includes `itypes/` for internal API types (e.g., `ApiResponseEdgeGateway`, `ApiRequestEdgeGateway`).

### 1.5 `cmd/`

**Command-line interface and generators.**

- Hosts CLI tools, code generators, and developer utilities.
- Used for development, testing, and automation.

### 1.6 `pkg/`

**Shared packages.**

- Contains reusable libraries, helpers, and utilities used across the SDK.

### 1.7 `ruleguard/`

**Custom lint rules.**

- Contains ruleguard scripts enforcing naming conventions and best practices.

---

## 2. Endpoint and Commands in `/api/`

### 2.1 Endpoint Design

- Each API group (e.g., `edgegateway`, `vdc`) defines endpoints as Go functions or methods that map to CloudAvenue REST API operations.
- Endpoints are responsible for:
 	- Building and sending HTTP requests.
 	- Parsing and validating responses.
 	- Mapping internal types (`ApiRequest*`, `ApiResponse*`) to user-facing models (`Model*`).
- Endpoints are registered and managed by the main `Client` of each API group.
- Endpoint logic is isolated per resource, making it easy to extend or maintain.

### 2.2 Commands Design

- Commands represent high-level operations exposed to SDK users (e.g., `GetEdgeGateway`, `CreateVDC`).
- Each command is implemented as a public method on the API group's `Client` struct.
- Commands:
 	- Accept user-supplied parameters (`Params*`).
 	- Internally call the appropriate endpoint function.
 	- Return user-facing models (`Model*`) or error.
- Commands are documented and validated by custom linters to ensure naming and usage consistency.

---

For more details, see the [CONTRIBUTING.md](./CONTRIBUTING.md) and [GUIDELINE.md](./GUIDELINE.md) files.
