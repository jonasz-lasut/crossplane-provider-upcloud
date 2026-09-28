// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// GroupKindOverrides overrides the group and kind of the resource if it matches
// any entry in the GroupMap.
func GroupKindOverrides() config.ResourceOption {
	return func(r *config.Resource) {
		if f, ok := GroupMap[r.Name]; ok {
			r.ShortGroup, r.Kind = f(r.Name)
		}
	}
}

// GroupKindCalculator returns the correct group and kind name for given TF
// resource.
type GroupKindCalculator func(resource string) (string, string)

// KnownGroupKind returns a GroupKindCalculator that assigns the given static
// group and kind regardless of the resource name.
func KnownGroupKind(group, kind string) GroupKindCalculator {
	return func(string) (string, string) { return group, kind }
}

// GroupMap assigns UpCloud resources to API groups and kinds. Resources are
// added here together with their entry in the external name tables.
var GroupMap = map[string]GroupKindCalculator{
	// database
	"upcloud_managed_database_logical_database": KnownGroupKind("database", "ManagedDatabaseLogicalDatabase"),
	"upcloud_managed_database_mysql":            KnownGroupKind("database", "ManagedDatabaseMysql"),
	"upcloud_managed_database_opensearch":       KnownGroupKind("database", "ManagedDatabaseOpensearch"),
	"upcloud_managed_database_postgresql":       KnownGroupKind("database", "ManagedDatabasePostgresql"),
	"upcloud_managed_database_user":             KnownGroupKind("database", "ManagedDatabaseUser"),
	"upcloud_managed_database_valkey":           KnownGroupKind("database", "ManagedDatabaseValkey"),
	"upcloud_managed_database_connection_pool":  KnownGroupKind("database", "ManagedDatabaseConnectionPool"),

	// network
	"upcloud_firewall_ruleset":                        KnownGroupKind("network", "FirewallRuleset"),
	"upcloud_floating_ip_address":                     KnownGroupKind("network", "FloatingIPAddress"),
	"upcloud_gateway":                                 KnownGroupKind("network", "Gateway"),
	"upcloud_gateway_connection":                      KnownGroupKind("network", "GatewayConnection"),
	"upcloud_gateway_connection_tunnel":               KnownGroupKind("network", "GatewayConnectionTunnel"),
	"upcloud_loadbalancer":                            KnownGroupKind("network", "LoadBalancer"),
	"upcloud_loadbalancer_backend":                    KnownGroupKind("network", "LoadBalancerBackend"),
	"upcloud_loadbalancer_backend_tls_config":         KnownGroupKind("network", "LoadBalancerBackendTLSConfig"),
	"upcloud_loadbalancer_dynamic_backend_member":     KnownGroupKind("network", "LoadBalancerDynamicBackendMember"),
	"upcloud_loadbalancer_dynamic_certificate_bundle": KnownGroupKind("network", "LoadBalancerDynamicCertificateBundle"),
	"upcloud_loadbalancer_frontend":                   KnownGroupKind("network", "LoadBalancerFrontend"),
	"upcloud_loadbalancer_frontend_rule":              KnownGroupKind("network", "LoadBalancerFrontendRule"),
	"upcloud_loadbalancer_frontend_tls_config":        KnownGroupKind("network", "LoadBalancerFrontendTLSConfig"),
	"upcloud_loadbalancer_manual_certificate_bundle":  KnownGroupKind("network", "LoadBalancerManualCertificateBundle"),
	"upcloud_loadbalancer_resolver":                   KnownGroupKind("network", "LoadBalancerResolver"),
	"upcloud_loadbalancer_static_backend_member":      KnownGroupKind("network", "LoadBalancerStaticBackendMember"),
	"upcloud_network":                                 KnownGroupKind("network", "Network"),
	"upcloud_network_peering":                         KnownGroupKind("network", "NetworkPeering"),
	"upcloud_router":                                  KnownGroupKind("network", "Router"),

	// objectstorage
	"upcloud_managed_object_storage":                 KnownGroupKind("objectstorage", "ManagedObjectStorage"),
	"upcloud_managed_object_storage_policy":          KnownGroupKind("objectstorage", "ManagedObjectStoragePolicy"),
	"upcloud_managed_object_storage_user":            KnownGroupKind("objectstorage", "ManagedObjectStorageUser"),
	"upcloud_managed_object_storage_user_access_key": KnownGroupKind("objectstorage", "ManagedObjectStorageUserAccessKey"),
	"upcloud_managed_object_storage_user_policy":     KnownGroupKind("objectstorage", "ManagedObjectStorageUserPolicy"),
	"upcloud_managed_object_storage_bucket":          KnownGroupKind("objectstorage", "ManagedObjectStorageBucket"),
	"upcloud_managed_object_storage_custom_domain":   KnownGroupKind("objectstorage", "ManagedObjectStorageCustomDomain"),
	"upcloud_managed_object_storage_static_site":     KnownGroupKind("objectstorage", "ManagedObjectStorageStaticSite"),

	// server
	"upcloud_firewall_rules":                  KnownGroupKind("server", "FirewallRules"),
	"upcloud_server":                          KnownGroupKind("server", "Server"),
	"upcloud_server_group":                    KnownGroupKind("server", "ServerGroup"),
	"upcloud_server_private_firewall_ruleset": KnownGroupKind("server", "ServerPrivateFirewallRuleset"),
	"upcloud_tag":                             KnownGroupKind("server", "Tag"),

	// storage
	"upcloud_file_storage":           KnownGroupKind("storage", "FileStorage"),
	"upcloud_file_storage_share":     KnownGroupKind("storage", "FileStorageShare"),
	"upcloud_file_storage_share_acl": KnownGroupKind("storage", "FileStorageShareACL"),
	"upcloud_storage":                KnownGroupKind("storage", "Storage"),
	"upcloud_storage_backup":         KnownGroupKind("storage", "StorageBackup"),
	"upcloud_storage_template":       KnownGroupKind("storage", "StorageTemplate"),

	// uks
	"upcloud_kubernetes_cluster":    KnownGroupKind("uks", "KubernetesCluster"),
	"upcloud_kubernetes_node_group": KnownGroupKind("uks", "KubernetesNodeGroup"),
}
