// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package server

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the server group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("upcloud_server", func(r *config.Resource) {
		r.UseAsync = true
		r.AddSingletonListConversion("template", "template")
		r.AddSingletonListConversion("login", "login")
		r.AddSingletonListConversion("simple_backup", "simpleBackup")
		// Required upstream through a list validator the schema dump drops.
		r.MarkAsRequired("network_interface")

		// The API reports cpu and mem for a plan-based server and plan for
		// a custom-sized one; late-initializing them into the spec makes
		// the two settings conflict on the next apply.
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"cpu", "mem", "plan"},
		}

		r.References["network_interface.network"] = config.Reference{
			TerraformName: "upcloud_network",
		}
		r.References["storage_devices.storage"] = config.Reference{
			TerraformName: "upcloud_storage",
		}
	})

	p.AddResourceConfigurator("upcloud_server_group", func(r *config.Resource) {
		r.References["members"] = config.Reference{
			TerraformName: "upcloud_server",
		}
	})

	p.AddResourceConfigurator("upcloud_tag", func(r *config.Resource) {
		r.References["servers"] = config.Reference{
			TerraformName: "upcloud_server",
		}
	})

	p.AddResourceConfigurator("upcloud_firewall_rules", func(r *config.Resource) {
		r.References["server_id"] = config.Reference{
			TerraformName: "upcloud_server",
		}
	})

	p.AddResourceConfigurator("upcloud_server_private_firewall_ruleset", func(r *config.Resource) {
		r.References["server_id"] = config.Reference{
			TerraformName: "upcloud_server",
		}
		r.References["ruleset_id"] = config.Reference{
			TerraformName: "upcloud_firewall_ruleset",
		}
	})
}
