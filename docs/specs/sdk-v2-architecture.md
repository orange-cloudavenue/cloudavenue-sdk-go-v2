---
title: CloudAvenue SDK v2 Architecture Specification
status: stable
created: 2026-08-18
---

# CloudAvenue SDK v2 Architecture Specification

## 0. Purpose and Scope

This document specifies the target architecture for a from-scratch rewrite of
`cloudavenue-sdk-go-v2`. It supersedes the implicit design of the current
(pre-rewrite) codebase, which is reflection-heavy, string-keyed at its public
API boundary, and carries significant runtime cost for zero compile-time
safety benefit (see §1 for the specific defects being corrected).

**Primary confirmed consumer:** `terraform-provider-cloudavenue`
(`github.com/orange-cloudavenue/terraform-provider-cloudavenue`). This SDK
rewrite exists specifically to back that provider. This is not a hypothetical
future consumer — it is the reason this rework is happening — and its real,
observed integration patterns (documented in §9) are binding requirements on
this spec, not just informative color.

**Authority relationship:** this SDK specification is authoritative. Where
the current provider's code conflicts with the shape defined here, the
provider adapts to the SDK, not the reverse. The provider's Terraform schema
definitions themselves are the one exception — they remain hand-written via
the existing `superschema.Schema{}` DSL (see §9.3) and are a Terraform-side
concern outside this SDK's scope.

---

## 1. Defects Being Corrected

The current implementation's public API surface routes every call through a
reflection-based, string-keyed command registry
(`cmds.Get(Namespace, Resource, Verb).Run(ctx, client, params)`), with `any`
typed parameters and return values recovered only via a type assertion in a
thin generated wrapper. Validation is performed by marshaling real params to
JSON, unmarshaling into a throwaway struct built at runtime via
`ompluscator/dynamic-struct`, and validating that dynamic struct instead of
the real one. Endpoint definitions embed `resty`-specific types directly in a
public struct. A subclient type-switch lives in shared request-dispatch code.
Mock mode is toggled via a magic sentinel string flipping a global mutable
bool. Full defect catalogue and code-smell inventory is preserved in the
project history; this spec does not re-derive it, it corrects it.

This rewrite eliminates all of the above. The replacement design is specified
below.

---

## 2. Core Abstraction: `Operation[P, R]` / `Execute[P, R]`

Leaf (single-HTTP-call) operations are represented as a generic,
compile-time-typed value:

```go
type Operation[P any, R any] struct {
	// Name is a stable, generator-populated identifier used for
	// observability (logging, tracing, metrics) and error context.
	// Format: "<Group>.<Verb><Resource>", e.g. "EdgeGateway.Get".
	Name string

	// Backend declares which backend this leaf operation targets.
	// This is a required, structured field — not inferred, not left
	// to convention — because backend migration (infrapi replacing
	// vmware verb-by-verb) is an explicit, ongoing project reality.
	Backend BackendTarget

	Endpoint *cav.Endpoint

	// Validate performs pure, structural validation only (see §5).
	// It must never perform I/O.
	Validate func(P) error

	// Transform maps typed params to the wire request body/path/query
	// values consumed by Endpoint.
	Transform func(P) (any, error)

	// Extract maps the raw response back to a typed result.
	Extract func(*cav.Response, P) (R, error)
}

type BackendTarget int

const (
	BackendInfrapi BackendTarget = iota
	BackendVMware
	BackendOSE      // Object Storage Extension (S3), see §8.4
	BackendNetBackup
)

func Execute[P, R any](ctx context.Context, c *Client, op Operation[P, R], params P) (R, error) {
	var zero R
	if op.Validate != nil {
		if err := op.Validate(params); err != nil {
			return zero, fmt.Errorf("%s: validate: %w", op.Name, err)
		}
	}
	body, err := op.Transform(params)
	if err != nil {
		return zero, fmt.Errorf("%s: transform: %w", op.Name, err)
	}
	resp, err := c.do(ctx, op.Endpoint, body)
	if err != nil {
		return zero, fmt.Errorf("%s: %w", op.Name, err)
	}
	result, err := op.Extract(resp, params)
	if err != nil {
		return zero, fmt.Errorf("%s: extract: %w", op.Name, err)
	}
	return result, nil
}
```

### 2.1 Design notes

- **`Extract` receives `P` (the original params), not just the response.**
  This is a deliberate correction versus the earlier proposal reviewed in
  architectural roundtable feedback. The legacy VCD `/api/query` backend
  requires knowing *which endpoint/URL pattern was hit* to know how to parse
  identifiers out of the response (HREF-regex extraction differs from
  cloudapi's direct URN fields). An `Extract` that only sees the raw response
  cannot express this. Passing `P` alongside the response closes this gap
  without requiring a second abstraction layer.
- **This abstraction covers leaf CRUD only.** Multi-call orchestration
  (fan-out lookups, async job submission + polling + resource-name
  extraction, conditional follow-up calls) is explicitly **not** expressed as
  a single `Operation[P,R]`. It is written as an ordinary Go method on
  `*Client` that calls `Execute` as a subroutine, potentially multiple times,
  potentially against operations targeting different backends. See §4.
- **`Name` is mandatory, not optional-and-added-later.** This closes a gap
  identified during design review: today's registry gives Namespace+Resource+
  Verb for free at one chokepoint for OTel span naming, structured log
  fields, and per-operation-type retry/rate-limit config. Scattering
  `Operation[P,R]` values as package-level vars loses this unless the field
  exists from day one. The generator populates it; it is never hand-typed
  for generated leaves.
- **`Backend` is a structured, generator-readable input, not informal
  tracking.** See §7 for how this drives per-verb migration status and
  generator behavior.

---

## 3. Validation: Structural vs. Business Rules — A Hard Boundary

This is a stated architectural principle, not a convention to be inferred:

> **`Validate()` is pure and structural. Anything requiring I/O is
> orchestration, never validation.**

Concretely:

- **Structural/cross-field validation** (e.g. "ID is required if Name is
  null", "Bandwidth must be between 1 and 10000") belongs in `Validate()`
  and/or declarative `validator`-style struct tags on the `Params*` type
  itself. This is pure, deterministic, requires no network access, and can
  run before any HTTP call is attempted.
- **Lookup-dependent business rules** (e.g. "T0Name must reference an
  existing T0 router visible to this org") are categorically NOT validation.
  They require a network call, and their result can be stale between
  validate-time and execute-time (TOCTOU). These live in orchestration
  methods (§4), expressed as ordinary Go function calls with ordinary Go
  error handling — never inside a `Validate()` method.

This distinction must be enforced by code review discipline (and ideally a
lint rule under ruleguard) — the specific failure mode being prevented is
a future engineer "helpfully" adding an API call inside a `Validate()` method
because it's convenient, silently reintroducing TOCTOU bugs and turning a
supposedly pure function into one with hidden I/O and hidden latency.

### 3.1 Struct tags are a required, load-bearing public contract

`Params*` structs retain declarative `validator`-style struct tags
(`validator:"required,min=1"` etc.) as **static, reflectable metadata**. This
is not optional or YAGNI: `terraform-provider-cloudavenue` is a confirmed,
real consumer, and while it does not currently reflect over SDK types to
build its Terraform schemas (see §9.3 — this is a net-new pattern, not
existing behavior being preserved), the SDK must not foreclose this
possibility by discarding tag-based introspection during the rewrite.
Retrofitting struct tags after the fact — once five API groups' worth of
hand-written `Validate()` methods exist with no tag-based residue — is
expensive. Preserving them now is nearly free.

The SDK exposes a formal, tested introspection helper built on Go 1.26's
`reflect.Type.Fields()` iterator API:

```go
// FieldsOf returns an iterator over P's exported struct fields, exposing
// name, type, and validator/documentation struct tags. This is the
// SDK-blessed mechanism for external reflection over Params*/Model* types
// (e.g. future schema-generation tooling), rather than leaving consumers to
// hand-roll reflection against undocumented struct internals.
func FieldsOf[P any]() iter.Seq[reflect.StructField]
```

This is provided as an available, tested contract now. Per explicit user
direction, no schema-generation tooling is built against it in this rework —
the Terraform provider's schema definitions remain hand-written via
`superschema.Schema{}` (§9.3) for the foreseeable future. Struct-tag-driven
schema generation is recorded here as an explicit **long-term possibility**,
not in scope for this rewrite.

### 3.2 Cross-field `ParamsRules` replacement

The current `ParamsRules` reflection-based cross-field rule evaluator is
deleted outright. Its structural cases (e.g. "id required if name null") are
expressed as validator struct tags (`required_without=Name`) where the tag
vocabulary supports it, or as a hand-written `Validate() error` method on the
`Params*` type for anything a tag can't express. This is a net quality
improvement: real Go code with real stack traces beats a dynamically
evaluated rules engine for cross-field logic, and is easier to unit test in
isolation.

---

## 4. Orchestrated (Composite) Operations

Operations requiring more than one HTTP call — parallel disambiguation
lookups, async job submission, conditional follow-up calls — are written as
ordinary hand-written methods on `*Client` (or the relevant per-group
`Client` wrapper), composing one or more calls to `Execute[P,R]` internally.

```go
const opCreateEdgeGateway = "EdgeGateway.Create"

func (c *Client) CreateEdgeGateway(ctx context.Context, params types.ParamsCreateEdgeGateway) (*types.ModelEdgeGateway, error) {
	// ... fan-out lookups, bandwidth validation, job submission via
	// Execute(ctx, c, createEdgeGatewayOp, ...), AwaitJob, conditional
	// follow-up PATCH ...
}
```

### 4.1 Identity convention for hand-written orchestration

Every orchestration method declares a package-level `const op<Name> = "..."`
identifier, used consistently in every wrapped error and log line inside that
method, mirroring the `Operation[P,R].Name` field's purpose for generated
leaves. This is a required convention, not a nicety — it is how orchestrated
operations remain greppable and traceable in production without the
free registry lookup the old design provided.

### 4.2 Backend targeting is per-call, not per-method

An orchestration method has no single "backend." Its constituent `Execute`
calls each carry their own `Backend` via their respective `Operation[P,R]`
values. `CreateEdgeGateway`'s three-way fan-out (VDC groups / VDCs / T0
routers) and its job-based create call and its conditional bandwidth PATCH
may each target different backends during the active infrapi migration. No
`Primary`/`Fallback` field is added to `Operation[P,R]` to model this at the
orchestration level — that model was evaluated and explicitly rejected
during design review because it cannot express control-flow shape
differences (Cerberus's async job pattern vs. VCD's synchronous REST) as a
mere endpoint swap. Migrating a leaf means rewriting/regenerating that one
leaf's `Operation[P,R]` (including its `Backend` field, `Transform`, and
`Extract`) and updating the orchestration method's call site if the leaf's
signature changed; it does not require a parallel-path abstraction.

### 4.3 Partial-failure contract (required, standardized)

This closes a gap flagged as the single riskiest unresolved item in design
review: today's opaque orchestration functions are at least *contained* —
one function, one place to look. The new design must not regress this by
leaving partial-failure semantics to each orchestration method's individual
discretion.

**Standardized contract:** any orchestration method that can fail after a
side effect has already partially succeeded (e.g. `CreateEdgeGateway`'s
async Job succeeds, creating the resource, but a subsequent bandwidth PATCH
fails) must return a non-nil result **and** a wrapped error satisfying
`ErrPartialSuccess`:

```go
// ErrPartialSuccess wraps an error occurring after a resource was already
// partially created/modified. Result carries whatever partial state is
// known; callers (including the Terraform provider) MUST check for this
// via errors.AsType and persist Result even though err != nil, to avoid
// orphaning a resource that exists server-side but is absent from state.
type ErrPartialSuccess[R any] struct {
	Result R
	Err    error
}

func (e *ErrPartialSuccess[R]) Error() string { return e.Err.Error() }
func (e *ErrPartialSuccess[R]) Unwrap() error { return e.Err }
```

Every orchestration method's doc comment must state explicitly whether it
can produce `ErrPartialSuccess` and under what conditions. This is a required
field in the per-operation acceptance criteria for every hand-written
orchestration method added during the rewrite — not an afterthought to be
added if someone remembers.

---

## 5. Errors

- `APIError` (existing rich diagnostic struct: Operation, StatusCode,
  StatusMessage, Duration, Endpoint, Method) is **retained**, not replaced.
  It gains `Unwrap() error` and participates in sentinel-based matching via
  wrapped sentinel errors (`ErrNotFound`, `ErrJobFailed`, `ErrJobTimeout` —
  see §6).
- All error inspection code — within the SDK and recommended to consumers —
  uses Go 1.26's `errors.AsType[E]` as the standard idiom in place of the
  classic `errors.As(err, &target)` two-step:

  ```go
  if apiErr, ok := errors.AsType[*APIError](err); ok { ... }
  ```

- `IsNotFound()`-style boolean helper methods may remain as thin sugar over
  `errors.Is(err, ErrNotFound)` for ergonomic parity with the provider's
  existing `cerrs.IsNotFound(err)` call sites (§9.2), easing the provider's
  adaptation.

---

## 6. Async Jobs: `AwaitJob[R]`

CloudAvenue's Cerberus/infrapi subsystem is fully async (submit → poll job
status → success/fail), while VCD/cloudapi is synchronous. The current
implementation detects this implicitly via reflection on
`BodyResponseType == cav.Job{}`. The new design makes this explicit: any
`Operation[P,R]` whose `Extract` needs to await a job does so by calling
`AwaitJob` directly inside its `Extract` function or inside the enclosing
orchestration method — there is no implicit type-sniffing.

```go
type JobPollOptions struct {
	Timeout         time.Duration
	PollingInterval time.Duration
	// Jitter adds randomized variance to PollingInterval to avoid
	// thundering-herd polling against Cerberus at scale (multiple API
	// groups, N concurrent operations, fixed-interval polling was
	// identified as a real production risk against Orange's own infra).
	Jitter time.Duration
}

// AwaitJob polls until the job referenced by jobID completes, is canceled,
// or opts.Timeout elapses — whichever comes first.
//
// Timeout vs. ctx deadline precedence: whichever fires first wins, and the
// resulting error is distinguishable via context.Cause:
//   - if opts.Timeout elapses first, ctx is canceled internally with cause
//     ErrJobTimeout (slow backend — alerting-relevant, paging-worthy).
//   - if the caller's own ctx is canceled/deadline-exceeded first, its
//     original cause is preserved and returned as-is (client-side
//     cancellation — a different, non-backend alerting category).
//
// Cancellation semantics: if ctx is canceled (either cause), AwaitJob stops
// polling and returns immediately. It does NOT currently attempt to signal
// Cerberus to cancel the job server-side — the job may continue running
// server-side, orphaned from the caller's perspective. This is a known,
// accepted limitation for this rewrite (Cerberus does not universally
// expose a job-cancel endpoint); it must be documented prominently on
// AwaitJob's doc comment, not silently assumed away.
//
// Observability: each poll attempt emits a debug-level slog record
// (job ID, elapsed, status) via the logger on ctx/Client, so a stuck
// poll is diagnosable without needing to reproduce with custom
// instrumentation.
func AwaitJob[R any](ctx context.Context, c *Client, jobID string, opts JobPollOptions, extract func(Job) (R, error)) (R, error) {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	timer := time.AfterFunc(opts.Timeout, func() { cancel(ErrJobTimeout) })
	defer timer.Stop()

	interval := opts.PollingInterval
	for {
		select {
		case <-ctx.Done():
			var zero R
			return zero, context.Cause(ctx)
		case <-time.After(withJitter(interval, opts.Jitter)):
		}
		job, err := c.getJobStatus(ctx, jobID)
		if err != nil {
			var zero R
			return zero, fmt.Errorf("await job %s: %w", jobID, err)
		}
		c.logger.Debug("job poll", "job_id", jobID, "status", job.Status)
		switch job.Status {
		case JobStatusSuccess:
			return extract(job)
		case JobStatusFailed:
			var zero R
			return zero, fmt.Errorf("job %s: %w", jobID, ErrJobFailed)
		}
		// still running, continue polling
	}
}
```

### 6.1 Design notes

- **`AwaitJob` standardizes only the poll loop** (timeout/interval/jitter/
  cancellation/status-check/logging). Resource-specific extraction — e.g.
  regex-parsing a newly created resource's name out of `Actions[].Details`
  — remains hand-written per call site via the `extract` callback. This is
  correct and intentional, not an oversight: extraction logic is
  fundamentally per-resource and does not generalize.
- Recommended testing approach: wrap `AwaitJob` tests in Go 1.26's GA
  `testing/synctest.Test(t, func(t *testing.T) { ... })` to get a virtualized
  clock, making cancellation-mid-poll and timeout-fires scenarios fast and
  deterministic to test without real sleeps.
- `ErrJobTimeout` is a new sentinel added alongside the existing
  `ErrJobFailed`.

---

## 7. Backend Migration Tracking (`Backend` field, generator input)

Per explicit direction, per-verb backend targeting is a **structured,
named generator input**, not informal/spreadsheet tracking. Every leaf
operation definition (the source the generator consumes to emit
`Operation[P,R]` values) declares its `Backend` explicitly. This makes:

- migration status queryable at build/lint time (e.g., "list all verbs still
  on `BackendVMware`" is a mechanical scan over generated `Operation[P,R]`
  values' `Name`+`Backend` fields, not a manual audit);
- migrating a single verb a scoped, mechanical change: update that leaf's
  source definition's `Backend` field, rewrite its `Transform`/`Extract`/
  `PathTemplate` to match the new backend's request/response shape,
  regenerate. No other leaf or orchestration method is touched unless it
  directly calls the migrated leaf.

This directly supports the reality that migration is happening **verb-by-
verb within a resource**, not resource-by-resource — e.g. `GetEdgeGateway`
might move to infrapi while `PatchEdgeGatewayBandwidth` stays on vmware for
longer, exactly the shape already seen in `CreateEdgeGateway`'s existing
conditional bandwidth-PATCH step.

---

## 8. Authentication and Multi-Backend Credentials

### 8.1 Core principle: infrapi is the single credential entrypoint

**Cerberus/infrapi is authoritative for authentication.** Users configure
credentials once, against infrapi. Every other backend used by the SDK
(VCD/vmware, S3/OSE, and NetBackup where feasible) derives its own access
credentials from the infrapi-authenticated session via a backend-specific
credential-exchange step, rather than requiring the user to separately
configure credentials for each backend added to the stack. This is a
deliberate usability and security-surface decision: adding a new backend to
the SDK must not mean asking users for a new, independent credential to
manage.

Bootstrap sequence:

1. Authenticate against Cerberus/infrapi using user-supplied credentials.
   This is the only credential the user directly provides.
2. Derive/attach a VCD session from the infrapi-authenticated context (VCD
   session is a secondary, derived step — not an independent auth flow).
3. For backends requiring additional protocol-specific credentials (e.g. S3
   access/secret keys, see §8.4), perform a backend-specific credential
   exchange call against infrapi to obtain them, on demand or eagerly at
   client construction, per backend's own needs.

This replaces the current design's implicit assumption that `ClientVmware`
and `ClientCerberus` each independently implement the `auth` interface with
symmetrical, unrelated credential logic. Under this spec, infrapi auth is
first-class and primary; other subclients' auth logic is expressed in terms
of it.

### 8.2 `SubClient` interface

```go
type SubClient interface {
	// ContextData returns backend-specific values (e.g. OrganizationID,
	// SiteID for vmware) needed by request construction. Replaces the
	// current concrete type-switch (`switch sc := sc.(type) { case *vmware: ... }`)
	// living in shared request-dispatch code — a direct polymorphism fix.
	ContextData(ctx context.Context) ContextData
}
```

No `Transport` interface is introduced (evaluated and rejected — see §11).

### 8.3 Mock mode

Mock mode is configured via a constructor functional option, not a global
mutable bool keyed off a magic sentinel string:

```go
func NewClient(organization string, opts ...ClientOption) (*Client, error)

func WithMockMode(enabled bool) ClientOption
```

This directly fixes a documented, acknowledged gotcha in the current
codebase (`isMockClient` package-level global) and, as a direct consequence,
makes it safe to run a mock client and a real client concurrently in the
same process — not possible today.

### 8.4 S3 / Object Storage Extension (OSE)

S3 is a structurally different backend (AWS S3-protocol, not VCD/Cerberus
REST) and is the first backend to exercise the infrapi-derived-credential
pattern end to end:

- A new `ClientOSE` subclient implements `SubClient`.
- Credential exchange: (1) org lookup against infrapi, (2) access/secret key
  fetch against infrapi for that org — both derived from the already-
  authenticated infrapi session, no separate user-supplied S3 credential.
- A separate AWS-S3-protocol client factory is constructed from the
  exchanged credentials, sitting outside the `cav.Client`/`Endpoint`
  abstraction (S3 uses AWS's own SDK and request-signing, not `cav.Endpoint`
  request/response shapes).
**Decided: `aws-sdk-go-v2`.** Confirmed — consistent with this rewrite's
general modernization direction (Go 1.26, generics-first design).

### 8.5 NetBackup — independent auth, confirmed by design

Investigated (v1 SDK's `pkg/clients/netbackup`, `pkg/consoles`, and the
infrapi/Cerberus OpenAPI spec): NetBackup **cannot** be brought under the
infrapi-derived-credential pattern used for VCD/S3. Concretely:

- NetBackup's current auth is a standalone OAuth2 `grant_type=password` flow
  (`POST /auth/token` with `username`/`password` directly against a
  per-console NetBackup endpoint, e.g. `backup1.cloudavenue.orange-business.com`,
  resolved from the console/region table, not from infrapi).
- The infrapi/Cerberus OpenAPI spec (`docs/infrapi/NGP_Api_for_Customer_Cerberus_Cloud_Avenue.yml`)
  exposes no NetBackup credential-exchange surface at all — there is nothing
  analogous to S3's org-lookup → access/secret-key exchange to derive
  NetBackup credentials from an infrapi session.
- NetBackup is a distinct third-party backup product (Veritas NetBackup)
  with its own IdP, not a service infrapi proxies or brokers credentials
  for.

**This is a confirmed, explicit exception to §8.1's infrapi-single-
entrypoint principle**, not an open question: NetBackup requires the user to
separately configure its own username/password (as today), and the SDK's
`ClientNetBackup` subclient's `auth` implementation manages this
independently, deriving nothing from the infrapi-authenticated session. This
should be documented prominently wherever SDK-wide credential configuration
is documented, so it doesn't read as an inconsistency.

---

## 9. Reconciliation with the Real Consumer (`terraform-provider-cloudavenue`)

This section records concrete, observed integration requirements gathered by
reading the actual provider codebase, to ensure this spec is not designed in
a vacuum.

### 9.1 "Fat model" objects with behavior methods → become named operations

The current v1 SDK (still in production use by the provider today) returns
stateful domain objects with behavior methods from `Get` calls, e.g.:

```go
edge, _ := edgegatewayClient.GetEdgeGateway(ctx, idOrName)
edge.NetworkServiceIsEnabled()
edge.EnableNetworkService(ctx)
```

Under this spec, `Execute[P,R]` returns plain `Model*` data structs with no
behavior methods — this is intentional and consistent with the rest of the
design. The equivalent of `EnableNetworkService` becomes its own named leaf
`Operation[P,R]` (e.g. `EdgeGateway.EnableNetworkService`) or, if it requires
orchestration, a hand-written `*Client` method — never a method on the
returned `Model*` value. The provider's adaptation work replaces
`edge.EnableNetworkService(ctx)` call sites with
`sdk.EnableEdgeGatewayNetworkService(ctx, params)`-shaped calls against the
new SDK.

### 9.2 Existing provider-layer conventions this SDK should stay compatible with

- **`cerrs.IsNotFound(err)`** — used in ~14+ files across
  `internal/provider/edgegw/*`. The new SDK's `errors.Is(err, ErrNotFound)`
  (with an optional `IsNotFound()` sugar helper, §5) gives the provider a
  direct, mechanical replacement at each call site.
- **`metrics.New(resourceName, orgName, verb)()` deferred instrumentation** —
  already present on every provider CRUD method today. This means the
  provider already tags every call with resource-type+org+verb at its own
  call sites, independent of the SDK. This reduces (but does not eliminate)
  the urgency of the SDK's own `Operation[P,R].Name` field (§2) — the
  provider does not strictly need it for its own metrics today — but the
  field remains required regardless, since the SDK's own internal logging/
  tracing (independent of what the provider does) still needs it, and other
  future SDK consumers should not be assumed to replicate the provider's
  instrumentation discipline.
- **`cloudavenue.Lock(ctx)`/`Unlock(ctx)`** — a provider-layer mutex around
  mutating calls, presumably serializing concurrent Terraform operations
  against the same edge gateway. This is a provider-layer concurrency
  concern; this SDK does not need to replicate it internally, but should not
  assume the provider will remove it — the SDK's own request/job-polling
  code must remain safe to call concurrently regardless.

### 9.3 Terraform schema generation: explicitly out of scope, long-term only

Confirmed via direct inspection: the provider hand-writes **every** Terraform
schema via Orange's own `superschema.Schema{}` DSL (78 files use it
directly), with **zero** existing reflection over SDK types anywhere in the
provider. The separate `cmd/types-generator` tool goes the opposite
direction (schema-first: generates Go structs *from* an already-hand-written
schema, not the reverse).

**Per explicit user decision:** this remains unchanged. Terraform schema
definitions stay hand-written. This SDK rewrite does not build, and is not
gated on, any struct-tag-driven schema generator. §3.1's `FieldsOf[P]`
introspection helper exists purely to keep that door open for a **long-term,
not-yet-scheduled** future capability — it is deliberately not consumed by
anything in this rewrite's scope.

### 9.4 `internal/client/client.go` — provider-side migration target

The provider's current central SDK integration point
(`internal/client/client.go`, `CloudAvenue` struct wrapping `*clientca.Client`
[v1] and a deprecated `*govcd.VCDClient`) is what the provider's adaptation
work eventually replaces with a v2 `*Client`. This is provider-side follow-on
work, sequenced after this SDK spec is implemented (§10), not a constraint
on this spec's own shape.

### 9.5 Domain-specific requirements carried forward from prior migration
planning

The following concrete, previously-identified technical requirements are
carried into this spec's domain-level acceptance criteria (not re-litigated
here, recorded for traceability):

- **IPSet/SecurityGroup inconsistency**: today, IPSet bypasses the SDK
  entirely via raw `govcd` calls in the provider, while SecurityGroup uses
  proper SDK types. Resolution: standardize both on the SecurityGroup-style
  higher-level SDK type pattern when implemented in v2 — IPSet should not
  require the provider to bypass the SDK.
- **`vdc_or_vdc_group` polymorphism**: a prior migration-planning decision
  pushed this polymorphism entirely to the provider layer, reasoning that
  v2's (then reflection/registry-based) architecture had no shared-
  interface/generic-dispatch mechanism to support it in the SDK. That
  reasoning is stale under this spec — `Operation[P,R]`'s generics may
  support a shared abstraction better than the old architecture did. This
  should be **re-examined**, not assumed settled, when the vdc/vdcgroup
  domain is implemented.
- **IAM (Users)**: **confirmed — no infrapi/Cerberus REST endpoint exists
  for user CRUD.** Per §8.1's infrapi-first principle, this domain would
  normally target `BackendInfrapi`, but since no such endpoint exists, IAM
  Users is a **confirmed, explicit exception**: `Backend` is set to
  `BackendVMware` (routing through the legacy XML AdminOrg API, as v1 does
  today) for all IAM Users leaf operations, for as long as no REST surface
  exists. This does not require a bespoke non-generated client path — the
  `Operation[P,R]`/`Execute[P,R]` abstraction and generator convention are
  backend-agnostic already (§2, §7); IAM Users leaves are generated exactly
  like any other `BackendVMware`-targeted leaf, just with XML-API-shaped
  `Transform`/`Extract` functions instead of JSON. If/when infrapi adds user
  CRUD, this is a normal per-verb migration (§7), not a structural change.
- **Org Properties**: gap-fill is blocked pending verification of whether
  v1's `infrapicustomerproxy` endpoint and v2's `UpdateOrganization` endpoint
  are literally the same backend operation (determines whether they should
  be merged into one leaf operation or kept separate). **Resolution
  principle (per §8.1, generalized beyond auth to all backend selection):**
  infrapi is authoritative — if Org Properties can be read/written through
  infrapi, `Backend` is `BackendInfrapi`; for any specific property/field
  infrapi does not support, that leaf (or that field's leaf, if split at
  finer granularity) falls back to `BackendVMware`. The infrapicustomerproxy-
  vs-UpdateOrganization identity question determines *whether one leaf or
  two* are needed, not *which backend wins* — the backend-selection rule
  itself is already settled by this principle and requires no further
  decision once the endpoint-identity question is answered during
  implementation.
- **AppPortProfile/SecurityGroup/IPSet shared logic**: **placement decided:**
  a new internal package `internal/inetworkobjects` (sibling to
  `internal/iendpoints`) holds logic shared across these three resources
  (e.g. common validation of port-range/protocol structures, common
  URN-vs-UUID normalization). This follows the existing codebase convention
  of internal, non-exported shared packages scoped by concern rather than
  duplicating logic per-domain package or introducing a new exported `pkg/`
  dependency for what is purely an internal implementation detail. Each of
  the three resources keeps its own `Operation[P,R]` definitions in their
  respective domain files; only the genuinely shared, non-domain-specific
  helper functions live in the shared internal package. This favors Go's
  standard idiom of "a little duplication is better than the wrong
  abstraction" for anything that isn't clearly identical across all three
  resources — only extract what is proven identical, not what merely looks
  similar today.
- **VApp**: requires a new `api/vapp/v1`-equivalent replacing not just a thin
  wrapper but also raw-`govcd` calls the provider currently makes directly.
  Specific behaviors that must be preserved as explicit acceptance criteria:
  delete sequence (`RemoveAllNetworks` → `tryUndeploy`-ignore-if-already-
  undeployed → `Delete`), lease-update-only-on-diff, guest-properties OVF
  `ProductSectionList` comparison logic, `IsVAppOrgNetwork`/`IsVAppNetwork`
  classification (consumed transitively by the provider's
  `internal/provider/common/vm/vm.go`), and status polling adapted to this
  spec's job-based async pattern (§6).
- **S3**: see §8.4.

---

## 10. Rollout Sequencing

Per explicit user decision, this SDK rewrite proceeds as **one coordinated
pass across all API groups** (edgegateway, vdc, vdcgroup, organization, draas,
s3, vapp, iam, netbackup), with both code generators (`cmd/command-generator`,
`cmd/endpoint-generator`) rewritten alongside to emit the new `Operation[P,R]`
shape. There is no per-group pilot and no reference-implementation gate before
fanning out — this overrides the previously-drafted `.opencode/plans/sdk-v2-
migration-master-plan.md`'s domain-by-domain lockstep sequencing (see §10.1).

### 10.1 Status of the prior migration plan document

`.opencode/plans/sdk-v2-migration-master-plan.md` (in the
`terraform-provider-cloudavenue` repo) predates this spec and this
authority decision. Its **sequencing model is superseded**: it was built
around "SDK catches up to v1 domain-by-domain, provider migrates in
lockstep," which assumed the provider's needs should pace and gate the
SDK's rewrite. Under this spec, the SDK is authoritative and is rewritten
in full first; the provider's adaptation is a follow-on phase, not an
interleaved one.

Its **domain-specific technical findings remain valid and are carried
forward** into §9.5 above and should inform per-domain acceptance criteria
during implementation. Its file-count-based domain priority ordering
(vdcg=33 provider files → edgegw=28 → elb=24 → s3=16 → vapp=4 → iam=3 →
org=5, netbackup included) remains useful for sequencing the **provider's**
adaptation phase after this SDK spec is implemented, even though it no
longer governs the SDK rewrite itself.

That document should be updated (status + a note pointing to this spec) once
this spec is reviewed/accepted, to avoid two live documents disagreeing
about sequencing.

### 10.2 Generator rewrite

Both `cmd/command-generator` and `cmd/endpoint-generator` are rewritten (not
incrementally patched) to emit `Operation[P,R]` values, including the
`Name` and `Backend` fields, from their source definitions.

### 10.3 Go version

`go.mod`'s `go` directive targets language version `1.26` (not merely an
installed 1.26 toolchain) — `errors.AsType`, `reflect.Type.Fields()`, and
`testing/synctest` GA all require the bumped language version to be
available. Note: `terraform-provider-cloudavenue`'s own `go.mod` currently
targets `go 1.25.8`, not 1.26 — this discrepancy should be resolved (likely
by bumping the provider alongside its SDK v2 adaptation work) but does not
block this SDK spec, since the SDK and provider are independent Go modules.

### 10.4 NetBackup

NetBackup is included in the coordinated pass. It requires a new subclient,
a new OAuth2 password-grant auth flow independent of the infrapi-first
principle in §8.1 (confirmed by §8.5), and a new job-status model. It is an
optional per-console service and no other domain depends on it.

---

## 11. Explicitly Rejected Alternatives

Recorded for traceability, so these are not re-litigated without new
information:

- **General-purpose `Transport` interface over resty.** Rejected as
  premature abstraction (YAGNI) — resty already abstracts `net/http`; the
  actual, narrow problem was resty types leaking into the public `Endpoint`
  struct. Fixed by translating to cav-native `Request`/`Response` types at
  one adapter boundary (`cav.Response` referenced in §2's `Extract`
  signature). No second `Transport` implementation is named or planned; add
  the interface later only if a concrete second implementation or a proven
  fake-based testing need materializes.
- **`Primary`/`Fallback` field on `Operation[P,R]` as a general migration
  mechanism.** Rejected — cannot express control-flow-shape differences
  between backends (async job vs. sync REST, different identifier schemes).
  Valid only for simple leaf CRUD with no orchestration; superseded by the
  per-leaf `Backend` field (§7) plus rewriting orchestration method
  internals when migrating composite operations (§4.2).
- **Dynamic-struct/JSON-roundtrip validation (`ompluscator/dynamic-struct`).**
  Deleted outright. Real params structs are validated directly.
- **Global mutable mock-mode sentinel.** Deleted outright in favor of
  constructor-option DI (§8.3).
- **Reflection-based `commands/` registry with string-keyed dispatch.**
  Deleted outright in favor of `Operation[P,R]`/`Execute[P,R]` compile-time
  dispatch for leaves and ordinary Go methods for orchestration.

---

## 12. Open Questions

All five questions previously open have been resolved:

1. **`aws-sdk-go` v1 vs v2 for the OSE/S3 subclient** (§8.4) — **decided:
   `aws-sdk-go-v2`.**
2. **NetBackup auth** (§8.5) — **decided: genuinely independent.** Confirmed
   no infrapi credential-exchange surface exists for NetBackup; it keeps its
   own OAuth2 password-grant flow as a documented, explicit exception to
   §8.1.
3. **IAM user CRUD REST availability** (§9.5) — **confirmed: not available.**
   No infrapi/Cerberus REST endpoint exists for user CRUD. IAM Users targets
   `BackendVMware` (legacy XML AdminOrg API) as a normal `Operation[P,R]`
   backend assignment; no bespoke non-generated path is needed.
4. **Org Properties backend identity** (`infrapicustomerproxy` vs.
   `UpdateOrganization`) (§9.5) — the endpoint-identity question itself
   still requires confirmation during implementation (determines one leaf
   vs. two), but the **backend-selection principle is resolved**: infrapi
   first, `BackendVMware` fallback per field/leaf where infrapi lacks
   support, per §8.1 generalized.
5. **AppPortProfile/SecurityGroup/IPSet shared-logic placement** (§9.5) —
   **decided:** a new internal package `internal/inetworkobjects` holds only
   genuinely shared, non-domain-specific helpers; domain-specific
   `Operation[P,R]` definitions stay in their own per-resource files.

No open questions remain blocking the start of implementation. Two items
still require a confirmation step *during* implementation, already scoped
above: the Org Properties endpoint-identity check (item 4) and the general
practice of re-verifying infrapi coverage per field before assigning
`BackendVMware` fallback.
