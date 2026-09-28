// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package objectstorage

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the objectstorage group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("upcloud_managed_object_storage", func(r *config.Resource) {
		r.UseAsync = true
	})

	for _, child := range []string{
		"upcloud_managed_object_storage_policy",
		"upcloud_managed_object_storage_user",
	} {
		p.AddResourceConfigurator(child, func(r *config.Resource) {
			r.References["service_uuid"] = config.Reference{
				TerraformName: "upcloud_managed_object_storage",
			}
		})
	}

	p.AddResourceConfigurator("upcloud_managed_object_storage_user_access_key", func(r *config.Resource) {
		r.References["service_uuid"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage",
		}
		r.References["username"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage_user",
		}
	})

	p.AddResourceConfigurator("upcloud_managed_object_storage_user_policy", func(r *config.Resource) {
		r.References["service_uuid"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage",
		}
		r.References["username"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage_user",
		}
		r.References["name"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage_policy",
		}
	})
}
