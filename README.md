# Provider UpCloud

`provider-upcloud` is a [Crossplane](https://crossplane.io/) provider for
[UpCloud](https://upcloud.com/), built with
[Upjet](https://github.com/crossplane/upjet) and backed by the
[`UpCloudLtd/upcloud`](https://github.com/UpCloudLtd/terraform-provider-upcloud)
Terraform provider. It is maintained jointly by UpCloud and Upbound in the
crossplane-contrib organization.

It exposes XRM-conformant managed resources for every resource of the
Terraform provider: servers, server groups, firewall rules and rulesets,
tags, storage with backups and templates, file storage with shares and ACLs,
networks, routers, network peerings, floating IPs, gateways with VPN
connections and tunnels, load balancers with backends, frontends, rules,
resolvers and certificate bundles, managed databases with users, logical
databases and connection pools, managed object storage with users, policies,
buckets, custom domains and static sites, and the UpCloud Kubernetes Service.
Every resource is available in two flavors: cluster-scoped
(`*.upcloud.crossplane.io`) and namespaced (`*.upcloud.m.crossplane.io`).

## Authentication

The provider authenticates with an UpCloud API token (recommended) or an API
username and password read from the referenced `Secret`. Create a token in the
[UpCloud Hub](https://hub.upcloud.com/) under Account > API tokens and store
it as JSON:

```json
{
  "token": "ucat_..."
}
```

or, for a username and password pair:

```json
{
  "username": "<UpCloud API username>",
  "password": "<UpCloud API password>"
}
```

The same JSON may carry the optional client settings of the upstream
Terraform provider, as numbers or numeric strings. They bound how long a
reconcile can wait on the UpCloud API; the defaults are the upstream ones.

```json
{
  "token": "ucat_...",
  "request_timeout_sec": 30,
  "retry_max": 2,
  "retry_wait_min_sec": 1,
  "retry_wait_max_sec": 30
}
```

## Getting Started

### 1. Install the provider

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-upcloud
spec:
  package: xpkg.crossplane.io/crossplane-contrib/provider-upcloud:v0.2.0
```

### 2. Create a credentials Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: upcloud-creds
  namespace: crossplane-system
type: Opaque
stringData:
  creds: |
    {
      "token": "ucat_..."
    }
```

### 3. Create a ProviderConfig

```yaml
apiVersion: upcloud.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: upcloud-creds
      namespace: crossplane-system
      key: creds
```

Namespaced managed resources reference a `ClusterProviderConfig` (or a
namespaced `ProviderConfig`) in the `upcloud.m.crossplane.io` group instead;
see [`examples/namespaced/providerconfig/`](examples/namespaced/providerconfig/).

### 4. Create a managed resource

Examples for each resource are available under
[`examples/cluster/`](examples/cluster/) and
[`examples/namespaced/`](examples/namespaced/).

## ProviderConfig fields

| Field | Required | Description |
|---|---|---|
| `spec.credentials.source` | Yes | One of `Secret`, `InjectedIdentity`, `Environment`, `Filesystem` |
| `spec.credentials.secretRef` | When source=Secret | Reference to the credentials Secret |
| `spec.reconciliationPolicy` | No | Rate-limiting policy for reconciliation |

## Developing

### Code generation

```console
make generate
```

This runs the Upjet code generator against the pinned `UpCloudLtd/upcloud`
Terraform provider schema and docs, and writes the generated APIs,
controllers, and CRDs for both the cluster-scoped and namespaced variants.

### Run locally against a cluster

```console
make run
```

### Run end-to-end tests

```console
UPTEST_EXAMPLE_LIST="examples/cluster/storage/v1alpha1/storage.yaml" \
UPTEST_CLOUD_CREDENTIALS='{"token": "ucat_..."}' \
make e2e
```

See [AGENTS.md](AGENTS.md) for the repository layout and the full development
workflow.

## Reporting issues

Please open an [issue](https://github.com/crossplane-contrib/provider-upcloud/issues)
for bug reports, feature requests, or questions.
