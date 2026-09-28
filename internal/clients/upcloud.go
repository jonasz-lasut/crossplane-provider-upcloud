// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package clients

import (
	"context"
	"encoding/json"
	"math"
	"strconv"

	"github.com/UpCloudLtd/terraform-provider-upcloud/upcloud"
	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	terraformsdk "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"
	tjresource "github.com/crossplane/upjet/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/terraform"

	clusterv1beta1 "github.com/crossplane-contrib/provider-upcloud/apis/cluster/v1beta1"
	namespacedv1beta1 "github.com/crossplane-contrib/provider-upcloud/apis/namespaced/v1beta1"
	"github.com/crossplane-contrib/provider-upcloud/internal/version"
)

const (
	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal upcloud credentials as JSON"

	keyToken    = "token"
	keyUsername = "username"
	keyPassword = "password"

	// Optional client settings, named after the upstream provider attributes
	// they are passed through to.
	keyRequestTimeoutSec = "request_timeout_sec"
	keyRetryMax          = "retry_max"
	keyRetryWaitMinSec   = "retry_wait_min_sec"
	keyRetryWaitMaxSec   = "retry_wait_max_sec"
)

// clientSettingKeys lists the optional integer settings of the upstream
// provider that the credentials JSON may carry next to the credentials.
var clientSettingKeys = []string{keyRequestTimeoutSec, keyRetryMax, keyRetryWaitMinSec, keyRetryWaitMaxSec}

// TerraformSetupBuilder builds a terraform.SetupFn that configures the
// in-process UpCloud Terraform provider (no-fork, no terraform CLI). The setup
// always carries the plugin-framework provider from the upjet configuration:
// upjet configures a fresh provider server from the returned Configuration on
// every connect. The terraform-plugin-sdk/v2 provider meta is configured only
// for the kinds upstream still implements with the SDK (the gateway connection
// resources), because configuring it validates the credentials against the
// UpCloud API on every call.
func TerraformSetupBuilder(p *ujconfig.Provider) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		pcSpec, err := resolveProviderConfig(ctx, client, mg)
		if err != nil {
			return terraform.Setup{}, errors.Wrap(err, "cannot resolve provider config")
		}

		data, err := resource.CommonCredentialExtractor(ctx, pcSpec.Credentials.Source, client, pcSpec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return terraform.Setup{}, errors.Wrap(err, errExtractCredentials)
		}

		creds := map[string]any{}
		if err := json.Unmarshal(data, &creds); err != nil {
			return terraform.Setup{}, errors.Wrap(err, errUnmarshalCredentials)
		}

		cfg, err := buildConfiguration(creds)
		if err != nil {
			return terraform.Setup{}, errors.Wrap(err, "cannot build provider configuration")
		}

		ps := terraform.Setup{
			Configuration:     cfg,
			FrameworkProvider: p.TerraformPluginFrameworkProvider,
		}
		if !usesSDKClient(p, mg) {
			return ps, nil
		}
		return ps, errors.Wrap(configureProviderMeta(ctx, &ps), "cannot configure the upcloud SDK provider client")
	}
}

// usesSDKClient reports whether the managed resource is reconciled with the
// terraform-plugin-sdk/v2 client, which needs the provider meta in the setup.
func usesSDKClient(p *ujconfig.Provider, mg resource.Managed) bool {
	tr, ok := mg.(tjresource.Terraformed)
	if !ok {
		return false
	}
	r, ok := p.Resources[tr.GetTerraformResourceType()]
	return ok && r.ShouldUseTerraformPluginSDKClient()
}

// sdkProvider returns a fresh terraform-plugin-sdk/v2 flavor of the UpCloud
// provider that reports this provider's User-Agent to the UpCloud API.
func sdkProvider() *schema.Provider {
	p := upcloud.Provider()
	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return upcloud.ProviderConfigure(ctx, d, version.UserAgent())
	}
	return p
}

// configureProviderMeta configures a fresh in-process SDKv2 provider instance
// with the resolved configuration and stores the resulting provider meta (the
// UpCloud API client) in the Setup, as required by upjet's TerraformPluginSDK
// connectors. A fresh instance is used per call because provider meta is
// per-ProviderConfig state.
func configureProviderMeta(ctx context.Context, ps *terraform.Setup) error {
	p := sdkProvider()
	diag := p.Configure(context.WithoutCancel(ctx), &terraformsdk.ResourceConfig{Config: ps.Configuration})
	if diag != nil && diag.HasError() {
		return errors.Errorf("failed to configure the upcloud provider: %v", diag)
	}
	ps.Meta = p.Meta()
	return nil
}

// buildConfiguration maps the credentials secret onto the UpCloud provider
// configuration block. An API token takes precedence; otherwise the
// username and password pair is used. The optional client settings
// (request_timeout_sec, retry_max, retry_wait_min_sec, retry_wait_max_sec)
// are passed through when present, as JSON numbers or numeric strings. Only
// the keys that are set are passed through, so the upstream provider applies
// its own defaults for the rest.
func buildConfiguration(creds map[string]any) (map[string]any, error) {
	cfg := map[string]any{}
	token, err := stringValue(creds, keyToken)
	if err != nil {
		return nil, err
	}
	username, err := stringValue(creds, keyUsername)
	if err != nil {
		return nil, err
	}
	password, err := stringValue(creds, keyPassword)
	if err != nil {
		return nil, err
	}
	switch {
	case token != "":
		cfg[keyToken] = token
	case username != "" && password != "":
		cfg[keyUsername] = username
		cfg[keyPassword] = password
	default:
		return nil, errors.New(`credentials secret needs a "token" key or both "username" and "password" keys`)
	}
	for _, key := range clientSettingKeys {
		v, ok := creds[key]
		if !ok {
			continue
		}
		n, err := integerValue(key, v)
		if err != nil {
			return nil, err
		}
		cfg[key] = n
	}
	return cfg, nil
}

// stringValue returns the string under key, "" when the key is absent, and
// an error when the value is not a string.
func stringValue(creds map[string]any, key string) (string, error) {
	v, ok := creds[key]
	if !ok || v == nil {
		return "", nil
	}
	s, ok := v.(string)
	if !ok {
		return "", errors.Errorf("credentials key %q must be a string, got %T", key, v)
	}
	return s, nil
}

// integerValue converts a JSON number or a numeric string into the int64 the
// upstream provider attributes expect.
func integerValue(key string, v any) (int64, error) {
	switch v := v.(type) {
	case float64:
		if v != math.Trunc(v) {
			return 0, errors.Errorf("credentials key %q must be an integer, got %v", key, v)
		}
		return int64(v), nil
	case string:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, errors.Wrapf(err, "credentials key %q must be an integer", key)
		}
		return n, nil
	default:
		return 0, errors.Errorf("credentials key %q must be an integer, got %T", key, v)
	}
}

func toSharedPCSpec(pc *clusterv1beta1.ProviderConfig) (*namespacedv1beta1.ProviderConfigSpec, error) {
	if pc == nil {
		return nil, nil
	}
	data, err := json.Marshal(pc.Spec)
	if err != nil {
		return nil, err
	}

	var mSpec namespacedv1beta1.ProviderConfigSpec
	err = json.Unmarshal(data, &mSpec)
	return &mSpec, err
}

func resolveProviderConfig(ctx context.Context, crClient client.Client, mg resource.Managed) (*namespacedv1beta1.ProviderConfigSpec, error) {
	switch managed := mg.(type) {
	case resource.LegacyManaged:
		return resolveLegacy(ctx, crClient, managed)
	case resource.ModernManaged:
		return resolveModern(ctx, crClient, managed)
	default:
		return nil, errors.New("resource is not a managed resource")
	}
}

func resolveLegacy(ctx context.Context, client client.Client, mg resource.LegacyManaged) (*namespacedv1beta1.ProviderConfigSpec, error) {
	configRef := mg.GetProviderConfigReference()
	if configRef == nil {
		return nil, errors.New(errNoProviderConfig)
	}
	pc := &clusterv1beta1.ProviderConfig{}
	if err := client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetProviderConfig)
	}

	t := resource.NewLegacyProviderConfigUsageTracker(client, &clusterv1beta1.ProviderConfigUsage{})
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackUsage)
	}

	return toSharedPCSpec(pc)
}

// resolveNamespacedSpec converts a namespaced ProviderConfig spec into the
// internal, fully-resolved ProviderConfigSpec. The namespaced spec omits the
// secret namespace, so it is resolved to the namespace of the referencing
// managed resource here.
func resolveNamespacedSpec(spec namespacedv1beta1.NamespacedProviderConfigSpec, namespace string) namespacedv1beta1.ProviderConfigSpec {
	resolved := namespacedv1beta1.ProviderConfigSpec{
		ReconciliationPolicy: spec.ReconciliationPolicy,
		Credentials: namespacedv1beta1.ProviderCredentials{
			Source: spec.Credentials.Source,
			CommonCredentialSelectors: xpv2.CommonCredentialSelectors{
				Fs:  spec.Credentials.Fs,
				Env: spec.Credentials.Env,
			},
		},
	}
	if spec.Credentials.SecretRef != nil {
		resolved.Credentials.SecretRef = spec.Credentials.SecretRef.ToSecretKeySelector(namespace)
	}
	return resolved
}

func resolveModern(ctx context.Context, crClient client.Client, mg resource.ModernManaged) (*namespacedv1beta1.ProviderConfigSpec, error) {
	configRef := mg.GetProviderConfigReference()
	if configRef == nil {
		return nil, errors.New(errNoProviderConfig)
	}

	pcRuntimeObj, err := crClient.Scheme().New(namespacedv1beta1.SchemeGroupVersion.WithKind(configRef.Kind))
	if err != nil {
		return nil, errors.Wrap(err, "unknown GVK for ProviderConfig")
	}
	pcObj, ok := pcRuntimeObj.(client.Object)
	if !ok {
		// This indicates a programming error, types are not properly generated
		return nil, errors.New("runtime object is not a client.Object")
	}

	// Namespace will be ignored if the PC is a cluster-scoped type
	if err := crClient.Get(ctx, types.NamespacedName{Name: configRef.Name, Namespace: mg.GetNamespace()}, pcObj); err != nil {
		return nil, errors.Wrap(err, errGetProviderConfig)
	}

	var pcSpec namespacedv1beta1.ProviderConfigSpec
	pcu := &namespacedv1beta1.ProviderConfigUsage{}
	switch pc := pcObj.(type) {
	case *namespacedv1beta1.ProviderConfig:
		pcSpec = resolveNamespacedSpec(pc.Spec, mg.GetNamespace())
	case *namespacedv1beta1.ClusterProviderConfig:
		pcSpec = pc.Spec
	default:
		return nil, errors.New("unknown provider config type")
	}
	t := resource.NewProviderConfigUsageTracker(crClient, pcu)
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackUsage)
	}
	return &pcSpec, nil
}
