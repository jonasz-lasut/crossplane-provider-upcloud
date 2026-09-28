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

	p.AddResourceConfigurator("upcloud_firewall_rules", func(r *config.Resource) {
		r.References["server_id"] = config.Reference{
			TerraformName: "upcloud_server",
		}
	})
}
