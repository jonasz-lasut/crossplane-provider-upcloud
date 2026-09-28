// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package objectstorage

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-upcloud/config/common"
)

// Configure configures the objectstorage group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("upcloud_managed_object_storage", func(r *config.Resource) {
		r.UseAsync = true
	})

	p.AddResourceConfigurator("upcloud_managed_object_storage_policy", func(r *config.Resource) {
		r.References["service_uuid"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage",
		}
	})

	p.AddResourceConfigurator("upcloud_managed_object_storage_user", func(r *config.Resource) {
		r.References["service_uuid"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage",
		}
	})

	p.AddResourceConfigurator("upcloud_managed_object_storage_user_access_key", func(r *config.Resource) {
		r.References["service_uuid"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage",
		}
		r.References["username"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage_user",
			Extractor:     common.ExtractObservedExternalName,
		}
	})

	p.AddResourceConfigurator("upcloud_managed_object_storage_bucket", func(r *config.Resource) {
		r.References["service_uuid"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage",
		}
	})

	p.AddResourceConfigurator("upcloud_managed_object_storage_custom_domain", func(r *config.Resource) {
		r.References["service_uuid"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage",
		}
	})

	p.AddResourceConfigurator("upcloud_managed_object_storage_static_site", func(r *config.Resource) {
		r.References["service_uuid"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage",
		}
		r.References["bucket_name"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage_bucket",
			Extractor:     common.ExtractObservedExternalName,
		}
	})

	p.AddResourceConfigurator("upcloud_managed_object_storage_user_policy", func(r *config.Resource) {
		r.References["service_uuid"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage",
		}
		r.References["username"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage_user",
			Extractor:     common.ExtractObservedExternalName,
		}
		r.References["name"] = config.Reference{
			TerraformName: "upcloud_managed_object_storage_policy",
			Extractor:     common.ExtractObservedExternalName,
		}
	})
}
