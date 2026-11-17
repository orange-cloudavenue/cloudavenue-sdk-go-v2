## 1.1.0 (Unreleased)
### :rotating_light: **Breaking Changes**

* cav: replace the public polling API with `WaitJob(ctx, client, JobWaitRequest)`, requiring an explicit backend and removing the legacy await helpers. (GH-155)

### :rocket: **New Features**

* cav: add `WithJobTimeout` and `WithJobPollInterval` client options. Consumers can tune job wait budgets once per client; the values apply to every `WaitJob` call that does not set an explicit `JobWaitRequest` value, overriding the package defaults (3-minute timeout, 2-second polling). (GH-155)
* feat(vapp): add VMware vApp runtime network and guest properties surfaces for provider migration. Adds runtime networkConfigSection read/create/add-org/update/remove operations, ProductSection guest-properties read/update APIs, public undeploy parameters, the matching public types, VMware endpoint definitions, generated wrappers, and focused tests. (GH-159)
* feat(vapp): add discard suspended state action. Adds DiscardSuspendedStateVApp and its VMware endpoint definition to the vApp lifecycle operations. (GH-159)
* feat(vapp): add vApp ACL operations and models. Adds GetVAppACL, SetVAppACL and RemoveVAppACL with the matching public models and parameters, subject kind/source and access-level validation, VMware endpoint definitions, internal ACL wire types, and focused tests. (GH-159)
* feat(vapp): deliver the vApp lifecycle and payload behavior verified against the live VCD 39.1 contract. Undeploy requests use the versioned VMware XML content type; DeleteVApp ignores only the explicit benign not-running undeploy case and tolerates the 404 answered by removeAllNetworks for vApps without networks; runtime network writes mutate the raw section fetched from VMware so untouched identity, features, router info and sibling network state are preserved. Guest properties round-trip with full fidelity (Key, Label, Type, Value, UserConfigurable) in stable order, with empty sections accepted and nil sections rejected. vApp job polling is configurable through the new WithJobTimeout and WithJobPollInterval client options (5-minute default timeout). Endpoints send the media types required by the API (query records JSON, instantiateVAppTemplateParams XML, productSections JSON) and instantiate/compose network payloads emit the ovf:Info and Configuration (ParentNetwork, FenceMode) blocks mandated by the VCD schema. (GH-159)

### :dependabot: **Dependencies**

* deps: bumps github.com/aws/aws-sdk-go-v2, github.com/aws/aws-sdk-go-v2/service/s3 from 1.47.8 to 1.48.0 (GH-160)
* deps: bumps github.com/aws/aws-sdk-go-v2, github.com/aws/aws-sdk-go-v2/service/s3, golang.org/x/sync from 1.47.4 to 1.47.8 (GH-156)
* deps: bumps resty.dev/v3 from 3.0.0-rc.3 to 3.0.0-rc.4 (GH-148)

## 1.0.0 (November 12, 2026)

# 0.0.1 (June 20, 2025)

* Init the project.
