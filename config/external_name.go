// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package config

import (
	"context"
	"fmt"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// terraformPluginSDKExternalNameConfigs contains all external name
// configurations for the resources the UpCloud Terraform provider still
// implements with terraform-plugin-sdk/v2 (as of v5.44.1 only
// upcloud_gateway_connection and upcloud_gateway_connection_tunnel, neither
// generated yet). A resource listed here is generated with the in-process SDK
// client and needs the SDK provider meta configured in internal/clients.
var terraformPluginSDKExternalNameConfigs = map[string]config.ExternalName{}

// terraformPluginFrameworkExternalNameConfigs contains all external name
// configurations for the resources the UpCloud Terraform provider implements
// with terraform-plugin-framework. A resource listed here is generated with
// the in-process framework client.
var terraformPluginFrameworkExternalNameConfigs = map[string]config.ExternalName{
	// database
	"upcloud_managed_database_logical_database": config.TemplatedStringAsIdentifier("name", "{{ .parameters.service }}/{{ .external_name }}"),
	"upcloud_managed_database_mysql":            config.IdentifierFromProvider,
	"upcloud_managed_database_opensearch":       config.IdentifierFromProvider,
	"upcloud_managed_database_postgresql":       config.IdentifierFromProvider,
	"upcloud_managed_database_user":             config.TemplatedStringAsIdentifier("username", "{{ .parameters.service }}/{{ .external_name }}"),
	"upcloud_managed_database_valkey":           config.IdentifierFromProvider,

	// network
	"upcloud_network": config.IdentifierFromProvider,
	"upcloud_router":  config.IdentifierFromProvider,

	// objectstorage
	"upcloud_managed_object_storage":                 config.IdentifierFromProvider,
	"upcloud_managed_object_storage_policy":          config.TemplatedStringAsIdentifier("name", "{{ .parameters.service_uuid }}/{{ .external_name }}"),
	"upcloud_managed_object_storage_user":            config.TemplatedStringAsIdentifier("username", "{{ .parameters.service_uuid }}/{{ .external_name }}"),
	"upcloud_managed_object_storage_user_access_key": managedObjectStorageUserAccessKey(),
	"upcloud_managed_object_storage_user_policy":     managedObjectStorageUserPolicy(),

	// server
	"upcloud_firewall_rules": config.IdentifierFromProvider,
	"upcloud_server":         config.IdentifierFromProvider,
	"upcloud_server_group":   config.IdentifierFromProvider,

	// storage
	"upcloud_storage": config.IdentifierFromProvider,

	// uks
	"upcloud_kubernetes_cluster":    config.IdentifierFromProvider,
	"upcloud_kubernetes_node_group": config.TemplatedStringAsIdentifier("name", "{{ .parameters.cluster }}/{{ .external_name }}"),
}

// managedObjectStorageUserAccessKey returns the external name configuration
// of upcloud_managed_object_storage_user_access_key. The Terraform ID is
// <service_uuid>/<username>/<access_key_id>, but access_key_id is assigned by
// the API and must never be seeded from the Kubernetes object name, so the
// identifier argument setter is a no-op and no argument is omitted from the
// generated parameters.
func managedObjectStorageUserAccessKey() config.ExternalName {
	e := config.TemplatedStringAsIdentifier("access_key_id", "{{ .parameters.service_uuid }}/{{ .parameters.username }}/{{ .external_name }}")
	e.SetIdentifierArgumentFn = func(_ map[string]any, _ string) {}
	e.OmittedFields = []string{}
	return e
}

// managedObjectStorageUserPolicy returns the external name configuration of
// upcloud_managed_object_storage_user_policy. Its only naming field (name)
// is at the same time a reference to the policy, so the external name is the
// policy name and the Terraform ID <service_uuid>/<username>/<name> is
// assembled from the parameters.
func managedObjectStorageUserPolicy() config.ExternalName {
	return config.ExternalName{
		SetIdentifierArgumentFn: func(_ map[string]any, _ string) {},
		GetIDFn: func(_ context.Context, _ string, parameters map[string]any, _ map[string]any) (string, error) {
			serviceUUID, ok := parameters["service_uuid"].(string)
			if !ok {
				return "", fmt.Errorf("service_uuid is not a string or is empty")
			}
			username, ok := parameters["username"].(string)
			if !ok {
				return "", fmt.Errorf("username is not a string or is empty")
			}
			name, ok := parameters["name"].(string)
			if !ok {
				return "", fmt.Errorf("name is not a string or is empty")
			}
			return fmt.Sprintf("%s/%s/%s", serviceUUID, username, name), nil
		},
		GetExternalNameFn: func(tfstate map[string]any) (string, error) {
			name, ok := tfstate["name"].(string)
			if !ok {
				return "", fmt.Errorf("name is not a string or is empty")
			}
			return name, nil
		},
	}
}

// ExternalNameConfigs contains all external name configurations for this
// provider, merged from the per-client tables above.
var ExternalNameConfigs = func() map[string]config.ExternalName {
	merged := make(map[string]config.ExternalName, len(terraformPluginSDKExternalNameConfigs)+len(terraformPluginFrameworkExternalNameConfigs))
	for name, e := range terraformPluginSDKExternalNameConfigs {
		merged[name] = e
	}
	for name, e := range terraformPluginFrameworkExternalNameConfigs {
		merged[name] = e
	}
	return merged
}()

// ExternalNameConfigurations applies all external name configs listed in the
// tables terraformPluginSDKExternalNameConfigs and
// terraformPluginFrameworkExternalNameConfigs and pins the API version of
// those resources to v1alpha1.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		e, configured := ExternalNameConfigs[r.Name]
		if !configured {
			return
		}
		r.ExternalName = e
		r.Version = versionV1Alpha1
	}
}

// terraformPluginSDKResourceList returns the list of resources generated with
// the terraform-plugin-sdk/v2 client, anchored for the include-list regex.
func terraformPluginSDKResourceList() []string {
	return resourceList(terraformPluginSDKExternalNameConfigs)
}

// terraformPluginFrameworkResourceList returns the list of resources
// generated with the terraform-plugin-framework client, anchored for the
// include-list regex.
func terraformPluginFrameworkResourceList() []string {
	return resourceList(terraformPluginFrameworkExternalNameConfigs)
}

func resourceList(t map[string]config.ExternalName) []string {
	l := make([]string, 0, len(t))
	for name := range t {
		// $ is added to match the exact string since the format is regex.
		l = append(l, name+"$")
	}
	return l
}
