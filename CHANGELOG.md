## 1.2.0 (Unreleased)
## 1.1.0 (November 20, 2026)
### :rotating_light: **Breaking Changes**

* cav: replace the public polling API with `WaitJob(ctx, client, JobWaitRequest)`, requiring an explicit backend and removing the legacy await helpers. (GH-155)

### :rocket: **New Features**

* cav: add `WithJobTimeout` and `WithJobPollInterval` client options. Consumers can tune job wait budgets once per client; the values apply to every `WaitJob` call that does not set an explicit `JobWaitRequest` value, overriding the package defaults (3-minute timeout, 2-second polling). (GH-155)
* feat(vapp): ComposeVApp now accepts an empty vApp without a sourced item. (GH-161)
* feat(vapp): add VMware vApp runtime network and guest properties surfaces for provider migration. Adds runtime networkConfigSection read/create/add-org/update/remove operations, ProductSection guest-properties read/update APIs, public undeploy parameters, the matching public types, VMware endpoint definitions, generated wrappers, and focused tests. (GH-159)
* feat(vapp): add discard suspended state action. Adds DiscardSuspendedStateVApp and its VMware endpoint definition to the vApp lifecycle operations. (GH-159)
* feat(vapp): add vApp ACL operations and models. Adds GetVAppACL, SetVAppACL and RemoveVAppACL with the matching public models and parameters, subject kind/source and access-level validation, VMware endpoint definitions, internal ACL wire types, and focused tests. (GH-159)
* feat(vapp): deliver the vApp lifecycle and payload behavior verified against the live VCD 39.1 contract. Undeploy requests use the versioned VMware XML content type; DeleteVApp ignores only the explicit benign not-running undeploy case and tolerates the 404 answered by removeAllNetworks for vApps without networks; runtime network writes mutate the raw section fetched from VMware so untouched identity, features, router info and sibling network state are preserved. Guest properties round-trip with full fidelity (Key, Label, Type, Value, UserConfigurable) in stable order, with empty sections accepted and nil sections rejected. vApp job polling is configurable through the new WithJobTimeout and WithJobPollInterval client options (5-minute default timeout). Endpoints send the media types required by the API (query records JSON, instantiateVAppTemplateParams XML, productSections JSON) and instantiate/compose network payloads emit the ovf:Info and Configuration (ParentNetwork, FenceMode) blocks mandated by the VCD schema. (GH-159)
* feat(vdc): add a type-agnostic Org VDC Network lookup (GetVDCNetwork by ID, or by name with a VDC/VDCGroup owner URN or name), fetch Org VDC Network listings across all pages, and expose the legacy network HREF on the model through VAppParentNetworkReference for use as a vApp ParentNetwork reference. (GH-161)

### :bug: **Bug Fixes**

* fix(vapp): align legacy Cloud Director vApp reads and runtime network writes with the wire contract. Versioned vendor media types are used on legacy `/api/vApp` endpoints, XML bodies are marshaled for the vendor `+xml` content type, and job polling sends a versioned Accept header. VApp path identifiers keep the required `vapp-` prefix, runtime network configuration follows the XSD element sequence and carries `GuestVlanAllowed` through read-modify-write cycles, and status/owner decoding accepts both the legacy numeric/object and string wire encodings with normalization to semantic string tokens. (GH-161)
* fix(vapp): align vApp creation, undeploy and delete with the legacy wire contract, verified against a live Cloud Director. The composeVApp and undeploy actions now send their mandatory vendor media types (generic XML yielded 415), undeployVAppParams uses the UndeployPowerAction element instead of the XSD-rejected legacy boolean, a synchronous 201 creation response is no longer polled as a task, the vApp URN is derived from the entity HREF when query records return a null id, DELETE sends a versioned Accept header (406 otherwise), and XML error bodies no longer decode to an empty error message. (GH-161)

### :dependabot: **Dependencies**

* deps: bumps github.com/aws/aws-sdk-go-v2, github.com/aws/aws-sdk-go-v2/service/s3 from 1.47.8 to 1.48.0 (GH-160)
* deps: bumps github.com/aws/aws-sdk-go-v2, github.com/aws/aws-sdk-go-v2/service/s3, golang.org/x/sync from 1.47.4 to 1.47.8 (GH-156)
* deps: bumps resty.dev/v3 from 3.0.0-rc.3 to 3.0.0-rc.4 (GH-148)

## 1.0.0 (November 12, 2026)

# 0.0.1 (June 20, 2025)

* Init the project.
