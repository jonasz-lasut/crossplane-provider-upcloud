#!/usr/bin/env bash
set -aeuo pipefail

# Delete the gateway chain, waiting for each kind, before the Routers and Networks the gateways are attached to are deleted.
${KUBECTL} delete gatewayconnectiontunnel.network.upcloud.m.crossplane.io --all --all-namespaces
${KUBECTL} delete gatewayconnection.network.upcloud.m.crossplane.io --all --all-namespaces
${KUBECTL} delete gateway.network.upcloud.m.crossplane.io --all --all-namespaces
