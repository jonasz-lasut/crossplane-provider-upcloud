// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package storage

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-upcloud/config/common"
)

// Configure configures the storage group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("upcloud_storage", func(r *config.Resource) {
		r.UseAsync = false
		r.AddSingletonListConversion("clone", "clone")
		r.AddSingletonListConversion("import", "import")
	})

	p.AddResourceConfigurator("upcloud_storage_backup", func(r *config.Resource) {
		r.UseAsync = true

		r.References["source_storage"] = config.Reference{
			TerraformName: "upcloud_storage",
		}
	})

	p.AddResourceConfigurator("upcloud_storage_template", func(r *config.Resource) {
		r.UseAsync = true

		r.References["source_storage"] = config.Reference{
			TerraformName: "upcloud_storage",
		}
	})

	p.AddResourceConfigurator("upcloud_file_storage", func(r *config.Resource) {
		r.UseAsync = true
		r.AddSingletonListConversion("network", "network")

		r.References["network.uuid"] = config.Reference{
			TerraformName: "upcloud_network",
		}
	})

	p.AddResourceConfigurator("upcloud_file_storage_share", func(r *config.Resource) {
		r.References["file_storage"] = config.Reference{
			TerraformName: "upcloud_file_storage",
		}
	})

	p.AddResourceConfigurator("upcloud_file_storage_share_acl", func(r *config.Resource) {
		r.References["file_storage"] = config.Reference{
			TerraformName: "upcloud_file_storage",
		}
		r.References["share_name"] = config.Reference{
			TerraformName: "upcloud_file_storage_share",
			Extractor:     common.ExtractObservedExternalName,
		}
	})
}
