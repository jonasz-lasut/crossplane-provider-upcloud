#!/usr/bin/env bash
set -aeuo pipefail

# Delete the gateway chain, waiting for each kind, before the Routers and Networks the gateways are attached to are deleted.
${KUBECTL} delete gatewayconnectiontunnel.network.upcloud.crossplane.io --all
${KUBECTL} delete gatewayconnection.network.upcloud.crossplane.io --all
${KUBECTL} delete gateway.network.upcloud.crossplane.io --all
