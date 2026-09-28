// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	"context"
	_ "embed"

	"github.com/UpCloudLtd/terraform-provider-upcloud/upcloud"
	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/crossplane-contrib/provider-upcloud/config/cluster"
	"github.com/crossplane-contrib/provider-upcloud/config/templates"
	"github.com/crossplane-contrib/provider-upcloud/internal/version"
)

const (
	resourcePrefix = "upcloud"
	modulePath     = "github.com/crossplane-contrib/provider-upcloud"
	// versionV1Alpha1 is the API version every generated resource is served
	// at. Promote a resource to v1beta1 once its example is covered by uptest.
	versionV1Alpha1 = "v1alpha1"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider(_ context.Context) (*ujconfig.Provider, error) {
	defaultResourceOptions := []ujconfig.ResourceOption{
		GroupKindOverrides(),
		ExternalNameConfigurations(),
	}

	pc := ujconfig.NewProvider(
		[]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("upcloud.crossplane.io"),
		ujconfig.WithIncludeList([]string{}),
		ujconfig.WithControllerTemplate(templates.ControllerTemplate),
		ujconfig.WithTerraformPluginSDKIncludeList(terraformPluginSDKResourceList()),
		ujconfig.WithTerraformPluginFrameworkIncludeList(terraformPluginFrameworkResourceList()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithTerraformProvider(terraformProvider()),
		ujconfig.WithTerraformPluginFrameworkProvider(terraformPluginFrameworkProvider()),
		ujconfig.WithSchemaTraversers(&ujconfig.SingletonListEmbedder{}),
		ujconfig.WithDefaultResourceOptions(defaultResourceOptions...),
	)

	// add custom config functions
	for _, configure := range cluster.ProviderConfiguration {
		configure(pc)
	}

	pc.ConfigureResources()

	registerTFConversions(pc)

	return pc, nil
}

// terraformProvider returns the in-process terraform-plugin-sdk/v2 flavor of
// the UpCloud provider, configured to report this provider's User-Agent.
func terraformProvider() *schema.Provider {
	p := upcloud.Provider()
	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return upcloud.ProviderConfigure(ctx, d, version.UserAgent())
	}
	return p
}

// terraformPluginFrameworkProvider returns the in-process
// terraform-plugin-framework flavor of the UpCloud provider, configured to
// report this provider's User-Agent.
func terraformPluginFrameworkProvider() fwprovider.Provider {
	return upcloud.NewWithUserAgent(version.UserAgent())
}
