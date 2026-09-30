// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package uks

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the uks (UpCloud Kubernetes Service) group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("upcloud_kubernetes_cluster", func(r *config.Resource) {
		r.UseAsync = true

		r.References["network"] = config.Reference{
			TerraformName: "upcloud_network",
		}
	})

	p.AddResourceConfigurator("upcloud_kubernetes_node_group", func(r *config.Resource) {
		r.AddSingletonListConversion("custom_plan", "customPlan")
		r.AddSingletonListConversion("gpu_plan", "gpuPlan")
		r.AddSingletonListConversion("cloud_native_plan", "cloudNativePlan")

		r.References["cluster"] = config.Reference{
			TerraformName: "upcloud_kubernetes_cluster",
		}
	})
}
