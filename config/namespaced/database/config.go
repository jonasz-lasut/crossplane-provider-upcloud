// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package database

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// serviceRefDescription documents the limitation that the service reference
// fields only resolve PostgreSQL services.
const serviceRefDescription = "The service to which the resource belongs. Please note that reference fields (`serviceRef` and `serviceSelector`) only work for PostgreSQL databases. For other databases you need to leverage compositions and patches to pass the database service ID to the `service` field. See https://docs.crossplane.io/latest/concepts/patch-and-transform/#patching-between-resources for more info."

// Configure configures the database group
func Configure(p *config.Provider) {
	for _, service := range []string{
		"upcloud_managed_database_mysql",
		"upcloud_managed_database_opensearch",
		"upcloud_managed_database_postgresql",
		"upcloud_managed_database_valkey",
	} {
		p.AddResourceConfigurator(service, func(r *config.Resource) {
			r.UseAsync = true

			r.References["network.uuid"] = config.Reference{
				TerraformName: "upcloud_network",
			}

			// node_states is reported by the API only; expose it as an
			// optional, computed attribute so it never becomes a required
			// parameter.
			if s, ok := r.TerraformResource.Schema["node_states"]; ok {
				s.Optional = true
				s.Computed = true
			}
		})
	}

	p.AddResourceConfigurator("upcloud_managed_database_logical_database", func(r *config.Resource) {
		// TODO: use a generic reference type if that is ever implemented in
		// upjet, see https://github.com/crossplane/upjet/issues/95
		r.References["service"] = config.Reference{
			TerraformName: "upcloud_managed_database_postgresql",
		}

		if s, ok := r.TerraformResource.Schema["service"]; ok {
			s.Description = serviceRefDescription
		}
	})

	p.AddResourceConfigurator("upcloud_managed_database_user", func(r *config.Resource) {
		r.References["service"] = config.Reference{
			TerraformName: "upcloud_managed_database_postgresql",
		}

		if s, ok := r.TerraformResource.Schema["service"]; ok {
			s.Description = serviceRefDescription
		}
	})
}
