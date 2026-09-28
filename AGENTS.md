# provider-upcloud - Agent Guide

This is a **Crossplane v2 Upjet provider** for [UpCloud](https://upcloud.com/). It wraps the [`UpCloudLtd/upcloud`](https://github.com/UpCloudLtd/terraform-provider-upcloud) Terraform provider (terraform-plugin-framework, imported in-process via `upcloud.NewWithUserAgent()`) and generates cluster-scoped and namespaced Crossplane managed resources from a declarative config, with no hand-written controllers and no Terraform CLI.

> **Adding a resource**: Use the `add-upjet-resource` skill.

---

## Architecture

### Dual-scope generation

Every resource is generated **twice** from the same `config/`:

| Scope | Root group | ProviderConfig kind | API version example |
|---|---|---|---|
| Cluster | `upcloud.crossplane.io` | `ProviderConfig` | `server.upcloud.crossplane.io/v1alpha1` |
| Namespaced | `upcloud.m.crossplane.io` | `ClusterProviderConfig` | `server.upcloud.m.crossplane.io/v1alpha1` |

`config/provider.go` (`GetProvider()`) builds the cluster configuration and `config/provider_namespaced.go` (`GetProviderNamespaced()`) the namespaced one. The generator in `cmd/generator/main.go` runs `pipeline.Run(pc, pns, rootDir)`.

CRD group formula: `<shortGroup>.<rootGroup>`, e.g. short group `server` + root group `upcloud.crossplane.io` gives `server.upcloud.crossplane.io`.

### Terraform client split

The upstream provider serves almost everything through terraform-plugin-framework; only `upcloud_gateway_connection` and `upcloud_gateway_connection_tunnel` are still terraform-plugin-sdk/v2 resources. `config/external_name.go` keeps one external-name table per client (`terraformPluginFrameworkExternalNameConfigs`, `terraformPluginSDKExternalNameConfigs`); the table a resource lands in decides which in-process client reconciles it. `internal/clients/upcloud.go` only wires the framework provider today; add the SDK provider meta there before generating an SDK resource.

### Codegen pipeline

`make generate` runs these steps in order:
1. Downloads Terraform 1.5.x via the `build/` submodule tools.
2. Initialises a minimal `.work/terraform/main.tf.json` and runs `terraform providers schema -json` to produce **`config/schema.json`** (do not edit).
3. Clones the Terraform provider repo (sparse, tag pinned) and scrapes its docs into **`config/provider-metadata.yaml`** (do not edit).
4. Runs `go generate ./apis/...` which invokes:
   - `cmd/generator/main.go`: rewrites `apis/cluster/`, `apis/namespaced/`, `internal/controller/cluster/`, `internal/controller/namespaced/`, `examples-generated/`, `config/generated.lst`.
   - `controller-gen`: `package/crds/`.
   - `angryjet`: managed-resource methodsets.
   - upjet `resolver`: rewrites cross-group reference resolvers to go through `internal/apis` (avoids import cycles between API groups).

---

## Repository structure

```
config/
  provider.go               <- resourcePrefix, modulePath, WithRootGroup, client wiring (cluster)   edit
  provider_namespaced.go    <- same for the namespaced scope                                        edit
  external_name.go          <- per-client external-name tables (also gate generation)               edit
  groups.go                 <- GroupMap: Terraform resource -> API group + Kind                      edit
  schema.json               <- Terraform provider schema                                            generated
  provider-metadata.yaml    <- scraped Terraform provider docs                                      generated
  generated.lst             <- list of generated resource names                                     generated
  templates/                <- controller template with the reconciliation policy wiring            edit rarely
  cluster/<group>/config.go <- per-group AddResourceConfigurator (cluster scope)                     edit
  namespaced/<group>/config.go <- same, namespaced scope (must match cluster)                        edit
  cluster/provider.go       <- registers group Configure functions (cluster)                        edit
  namespaced/provider.go    <- registers group Configure functions (namespaced)                     edit

apis/
  cluster/                  <- generated API types (zz_* files) + hand-written v1beta1 ProviderConfig
  namespaced/               <- generated API types (zz_* files) + hand-written v1beta1 ProviderConfig
  generate.go               <- //go:generate directives that drive make generate                    edit rarely

internal/
  apis/                     <- runtime scheme used by the transformed reference resolvers           edit rarely
  clients/                  <- Terraform setup fn; provider credential extraction                    edit
  controller/cluster/       <- generated controllers + hand-written providerconfig controller
  controller/namespaced/    <- generated controllers + hand-written providerconfig controller
  features/                 <- feature flags                                                        edit rarely
  version/                  <- build version and the User-Agent reported to the UpCloud API         edit rarely

cmd/
  generator/main.go         <- runs pipeline.Run(GetProvider(), GetProviderNamespaced())            edit rarely
  provider/main.go          <- provider binary entry point                                         edit rarely

package/
  crossplane.yaml           <- provider package metadata                                            edit
  crds/                     <- generated CRD manifests                                              generated

examples/
  cluster/<group>/          <- curated E2E examples (cluster scope)                                  edit
  namespaced/<group>/       <- curated E2E examples (namespaced scope)                               edit
  cluster/providerconfig/   <- ProviderConfig + Secret template                                      edit
  namespaced/providerconfig/<- ProviderConfig + ClusterProviderConfig + Secret template              edit
  install.yaml              <- provider install manifest                                            edit

examples-generated/         <- raw generated examples (starting point only)                          generated

cluster/test/setup.sh       <- E2E cluster setup: creates Secret + ProviderConfigs                   edit

Makefile                    <- provider knobs at the top (see below)                                edit
build/                      <- crossplane/build submodule, never edit
```

Never edit generated files or the submodule directly.

---

## Makefile knobs

The provider-specific variables at the top of `Makefile`:

```makefile
PROVIDER_NAME                  # upcloud (PROJECT_NAME is derived as provider-$(PROVIDER_NAME))
TERRAFORM_PROVIDER_SOURCE      # UpCloudLtd/upcloud
TERRAFORM_PROVIDER_REPO        # https://github.com/UpCloudLtd/terraform-provider-upcloud (for pulling docs)
TERRAFORM_PROVIDER_VERSION     # pinned UpCloud Terraform provider version
TERRAFORM_DOCS_PATH            # docs/resources (docs dir inside the Terraform provider repo)
```

`PROJECT_REPO` is derived as `github.com/crossplane-contrib/$(PROJECT_NAME)`, matching the Go module path `github.com/crossplane-contrib/provider-upcloud`.

The Go module pins `github.com/UpCloudLtd/terraform-provider-upcloud` by commit pseudo-version, not by tag: the upstream module path has no `/v5` suffix, so its `v5.x` tags are not valid Go module versions. When bumping `TERRAFORM_PROVIDER_VERSION`, run `go get github.com/UpCloudLtd/terraform-provider-upcloud@<commit of the tag>`.

---

## Key commands

```bash
# After any fresh clone or worktree: run first
git submodule update --init --recursive

# Full codegen (schema pull + doc scrape + code generation)
make generate

# Run only the Go generator (schema + metadata already present)
go run cmd/generator/main.go "$PWD"

# Verify compilation, lint, unit tests and the example manifests
go build ./...
make lint
make test
make example-lint

# Run provider out-of-cluster (needs a kubeconfig)
make run

# E2E test (builds provider, spins up a KinD cluster, runs uptest/chainsaw)
UPTEST_EXAMPLE_LIST="examples/cluster/<group>/v1alpha1/<kind>.yaml" \
UPTEST_CLOUD_CREDENTIALS='{"token": "<UpCloud API token>"}' \
make e2e

# After any E2E run: ALWAYS clean up in this order
kubectl delete managed --all --all-namespaces
kind delete cluster --name local-dev
```

---

## Config key concepts

### The external-name tables gate generation

A resource absent from both tables in `config/external_name.go` is **never generated**, even if it is in `GroupMap`. Adding a resource to the table of the client that implements it upstream is always step one; every generated resource is served at `v1alpha1` until its example is covered by uptest.

### Group and Kind naming

`config/groups.go` `GroupMap` assigns every generated resource an explicit short group and Kind with `KnownGroupKind`. Add the new resource there together with its external-name entry.

### References

Cross-resource references (`r.References["field"] = config.Reference{...}`) go in both `config/cluster/<group>/config.go` AND `config/namespaced/<group>/config.go`. These files must stay in sync.

Pick the extractor by what the Terraform argument takes and how the target gets its external name:
- a UUID of an `IdentifierFromProvider` parent: no `Extractor` (the default external name is set only once the parent has been created);
- the bare name of a `TemplatedStringAsIdentifier` child (for example `default_backend_name`, `share_name`): `common.ExtractObservedExternalName`, which waits for the child to be observed. The plain external name is stamped by the name initializer before the child exists in UpCloud and produces 404s on create;
- a composite id such as `<lb uuid>/<backend name>` (for example a backend member's `backend`): `common.ExtractObservedPath("id")`;
- any other observed attribute: `common.ExtractObservedPath("<tf path>")`.

### Singleton lists

Terraform blocks limited to one item (for example `template` on a server, `ip_network` on a network, or the managed database `properties` block and everything nested in it) are embedded as objects in the CRD through explicit `r.AddSingletonListConversion("<tf_path>", "<crdPath>")` calls in the group configurators. The plugin-framework schema dump does not carry the `MaxItems=1` constraint (upstream uses `SizeAtMost(1)` validators), so `SingletonListEmbedder` cannot detect them: when adding a resource, look up its single-item blocks in the upstream source and list them. Write them as objects, not single-item lists, in examples.

### Empty collection defaults

A plugin-framework set attribute that defaults to an empty set and carries `RequiresReplace` (the node group's `ssh_keys`) becomes un-updatable after a provider restart or the uptest import step: the rebuilt prior state has it as null (omitempty drops the empty set from `status.atProvider`), the plan carries the empty default, and upjet's replacement filter compares the raw values. `r.TerraformConversions = append(r.TerraformConversions, common.EmptyListDefaults("<tf_name>"))` sends the empty list on both sides. Only for attributes whose upstream Read tolerates an empty non-null value; never for singleton-list blocks (Storage `import` is such a block and upstream's Read treats a non-null set as a configured import).

### Sensitive fields

Any Terraform attribute marked `Sensitive: true` becomes `<field>SecretRef` in the CRD (references a `v1/Secret`). The plain field name is absent from `forProvider`. Check the CRD when a field seems missing.

---

## E2E setup

`cluster/test/setup.sh` runs during `make e2e`. It:
1. Creates a `provider-secret` Secret in `crossplane-system` from `$UPTEST_CLOUD_CREDENTIALS` (JSON with a `token` key holding an UpCloud API token, or `username` and `password` keys).
2. Applies a `ProviderConfig` (cluster scope) and a `ClusterProviderConfig` (namespaced scope).

Examples pin their zone (`fi-hel2`) and object storage region (`europe-1`, primary zone `fi-hel2`) directly; no uptest datasource is needed.
Every example that creates a `Network` owns a distinct `/24` within its scope (`10.<n>.0.0/24` for network, server and load balancer examples, `172.16.<n>.0/24` for UKS, gateway and file storage examples) so a batch can run concurrently; pick an unused block for a new example and keep the IPs inside it (`ipAddress`, backend member `ip`, ACL `target`) in the same block.
When a teardown has to be ordered (the Gateway examples: the `Gateway` chain must be gone before its `Router` and `Network` are deleted), the parent docs carry `uptest.upbound.io/pre-delete-hook: ../testhooks/<script>.sh`. uptest runs the script right before that doc's own non-blocking delete, so the script deletes the dependent kinds with `${KUBECTL} delete <kind> --all`, which waits for each to be gone (every resource of a batch is being torn down at that point anyway). Scripts live in `examples/<scope>/<group>/testhooks/`, one copy per scope (group suffix and namespace differ), executable, with a shebang.

---

## Crossplane v2 notes

- Both scopes are `SafeStart`-capable (`package/crossplane.yaml`).
- Cluster-scoped ProviderConfig kind for namespaced resources is **`ClusterProviderConfig`** (not `ProviderConfig`).
- Managed resources in the namespaced scope reference it with `providerConfigRef.kind: ClusterProviderConfig`.
- Namespace for namespaced managed resources and Secrets: **`crossplane-system`** in examples (or whatever namespace the user creates them in).
