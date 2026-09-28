// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package uks

import (
	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/crossplane-contrib/provider-upcloud/config/common"
)

// Configure configures the uks (UpCloud Kubernetes Service) group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("upcloud_kubernetes_cluster", func(r *config.Resource) {
		r.UseAsync = true

		r.References["network"] = config.Reference{
			TerraformName: "upcloud_network",
		}

		// labels is optional+computed upstream; the framework provider
		// needs it sent explicitly (an empty map) to converge, so it is a
		// required, non-computed parameter.
		if s, ok := r.TerraformResource.Schema["labels"]; ok {
			s.Optional = false
			s.Computed = false
			s.Required = true
		}
		r.TerraformConversions = append(r.TerraformConversions, common.EmptyValueDefaults([]string{"labels"}, nil))
	})

	p.AddResourceConfigurator("upcloud_kubernetes_node_group", func(r *config.Resource) {
		r.References["cluster"] = config.Reference{
			TerraformName: "upcloud_kubernetes_cluster",
		}
	})
}
