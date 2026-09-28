// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package network

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-upcloud/config/common"
)

// Configure configures the network group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("upcloud_network", func(r *config.Resource) {
		r.AddSingletonListConversion("ip_network", "ipNetwork")
		// Required upstream through a list validator the schema dump drops.
		r.MarkAsRequired("ip_network")

		r.References["router"] = config.Reference{
			TerraformName: "upcloud_router",
		}
	})

	p.AddResourceConfigurator("upcloud_network_peering", func(r *config.Resource) {
		// Required upstream through list validators the schema dump drops.
		r.MarkAsRequired("network", "peer_network")

		r.References["network.uuid"] = config.Reference{
			TerraformName: "upcloud_network",
		}
		r.References["peer_network.uuid"] = config.Reference{
			TerraformName: "upcloud_network",
		}
	})

	p.AddResourceConfigurator("upcloud_floating_ip_address", func(r *config.Resource) {
		r.References["mac_address"] = config.Reference{
			TerraformName: "upcloud_server",
			Extractor:     common.ExtractObservedPath("network_interface[0].mac_address"),
		}
	})

	p.AddResourceConfigurator("upcloud_gateway", func(r *config.Resource) {
		// Required upstream through a list validator the schema dump drops.
		r.MarkAsRequired("router")
		r.UseAsync = true
		r.AddSingletonListConversion("router", "router")
		r.AddSingletonListConversion("address", "address")

		r.References["router.id"] = config.Reference{
			TerraformName: "upcloud_router",
		}
	})

	p.AddResourceConfigurator("upcloud_gateway_connection", func(r *config.Resource) {
		r.References["gateway"] = config.Reference{
			TerraformName: "upcloud_gateway",
		}
	})

	p.AddResourceConfigurator("upcloud_gateway_connection_tunnel", func(r *config.Resource) {
		r.References["connection_id"] = config.Reference{
			TerraformName: "upcloud_gateway_connection",
		}
	})

	p.AddResourceConfigurator("upcloud_loadbalancer", func(r *config.Resource) {
		r.UseAsync = true

		r.References["network"] = config.Reference{
			TerraformName: "upcloud_network",
		}
		r.References["networks.network"] = config.Reference{
			TerraformName: "upcloud_network",
		}
	})

	p.AddResourceConfigurator("upcloud_loadbalancer_backend", func(r *config.Resource) {
		r.AddSingletonListConversion("properties", "properties")

		r.References["loadbalancer"] = config.Reference{
			TerraformName: "upcloud_loadbalancer",
		}
		r.References["resolver_name"] = config.Reference{
			TerraformName: "upcloud_loadbalancer_resolver",
			Extractor:     common.ExtractObservedExternalName,
		}
	})

	p.AddResourceConfigurator("upcloud_loadbalancer_dynamic_backend_member", func(r *config.Resource) {
		r.References["backend"] = config.Reference{
			TerraformName: "upcloud_loadbalancer_backend",
			Extractor:     common.ExtractObservedPath("id"),
		}
	})

	p.AddResourceConfigurator("upcloud_loadbalancer_static_backend_member", func(r *config.Resource) {
		r.References["backend"] = config.Reference{
			TerraformName: "upcloud_loadbalancer_backend",
			Extractor:     common.ExtractObservedPath("id"),
		}
	})

	p.AddResourceConfigurator("upcloud_loadbalancer_backend_tls_config", func(r *config.Resource) {
		r.References["backend"] = config.Reference{
			TerraformName: "upcloud_loadbalancer_backend",
			Extractor:     common.ExtractObservedPath("id"),
		}
		r.References["certificate_bundle"] = config.Reference{
			TerraformName: "upcloud_loadbalancer_manual_certificate_bundle",
		}
	})

	p.AddResourceConfigurator("upcloud_loadbalancer_frontend", func(r *config.Resource) {
		r.AddSingletonListConversion("properties", "properties")

		r.References["loadbalancer"] = config.Reference{
			TerraformName: "upcloud_loadbalancer",
		}
		r.References["default_backend_name"] = config.Reference{
			TerraformName: "upcloud_loadbalancer_backend",
			Extractor:     common.ExtractObservedExternalName,
		}
	})

	p.AddResourceConfigurator("upcloud_loadbalancer_frontend_rule", func(r *config.Resource) {
		r.AddSingletonListConversion("actions", "actions")
		r.AddSingletonListConversion("matchers", "matchers")

		r.References["frontend"] = config.Reference{
			TerraformName: "upcloud_loadbalancer_frontend",
			Extractor:     common.ExtractObservedPath("id"),
		}
		r.References["actions.use_backend.backend_name"] = config.Reference{
			TerraformName: "upcloud_loadbalancer_backend",
			Extractor:     common.ExtractObservedExternalName,
		}
	})

	p.AddResourceConfigurator("upcloud_loadbalancer_frontend_tls_config", func(r *config.Resource) {
		r.References["frontend"] = config.Reference{
			TerraformName: "upcloud_loadbalancer_frontend",
			Extractor:     common.ExtractObservedPath("id"),
		}
		r.References["certificate_bundle"] = config.Reference{
			TerraformName: "upcloud_loadbalancer_manual_certificate_bundle",
		}
	})

	p.AddResourceConfigurator("upcloud_loadbalancer_resolver", func(r *config.Resource) {
		r.References["loadbalancer"] = config.Reference{
			TerraformName: "upcloud_loadbalancer",
		}
	})
}
