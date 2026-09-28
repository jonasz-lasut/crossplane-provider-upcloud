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
		r.References["router"] = config.Reference{
			TerraformName: "upcloud_router",
		}
	})

	p.AddResourceConfigurator("upcloud_router", func(r *config.Resource) {
		// labels and static_route are optional+computed upstream; the
		// framework provider needs them sent explicitly (an empty map or
		// list) to converge, so they are required, non-computed parameters.
		for _, field := range []string{"labels", "static_route"} {
			if s, ok := r.TerraformResource.Schema[field]; ok {
				s.Optional = false
				s.Computed = false
				s.Required = true
			}
		}
		r.TerraformConversions = append(r.TerraformConversions, common.EmptyValueDefaults([]string{"labels"}, []string{"static_route"}))
	})
}
