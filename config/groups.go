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

	// network
	"upcloud_network": KnownGroupKind("network", "Network"),
	"upcloud_router":  KnownGroupKind("network", "Router"),

	// objectstorage
	"upcloud_managed_object_storage":                 KnownGroupKind("objectstorage", "ManagedObjectStorage"),
	"upcloud_managed_object_storage_policy":          KnownGroupKind("objectstorage", "ManagedObjectStoragePolicy"),
	"upcloud_managed_object_storage_user":            KnownGroupKind("objectstorage", "ManagedObjectStorageUser"),
	"upcloud_managed_object_storage_user_access_key": KnownGroupKind("objectstorage", "ManagedObjectStorageUserAccessKey"),
	"upcloud_managed_object_storage_user_policy":     KnownGroupKind("objectstorage", "ManagedObjectStorageUserPolicy"),

	// server
	"upcloud_firewall_rules": KnownGroupKind("server", "FirewallRules"),
	"upcloud_server":         KnownGroupKind("server", "Server"),
	"upcloud_server_group":   KnownGroupKind("server", "ServerGroup"),

	// storage
	"upcloud_storage": KnownGroupKind("storage", "Storage"),

	// uks
	"upcloud_kubernetes_cluster":    KnownGroupKind("uks", "KubernetesCluster"),
	"upcloud_kubernetes_node_group": KnownGroupKind("uks", "KubernetesNodeGroup"),
}
