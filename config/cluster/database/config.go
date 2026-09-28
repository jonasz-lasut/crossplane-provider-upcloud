// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package database

import (
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/crossplane-contrib/provider-upcloud/config/common"
)

// serviceRefDescription documents the limitation that the service reference
// fields only resolve PostgreSQL services.
const serviceRefDescription = "The service to which the resource belongs. Please note that reference fields (`serviceRef` and `serviceSelector`) only work for PostgreSQL databases. For other databases you need to leverage compositions and patches to pass the database service ID to the `service` field. See https://docs.crossplane.io/latest/concepts/patch-and-transform/#patching-between-resources for more info."

// logicalDatabaseNotFound reports whether a Read failed because the logical
// database is missing. The external name is the database name, so upjet reads
// <service>/<name> before the first create; upstream looks the database up in
// the service's list and reports a miss as an error diagnostic instead of an
// empty state.
func logicalDatabaseNotFound(diags []*tfprotov6.Diagnostic) bool {
	for _, d := range diags {
		if d.Summary == "Unable to read logical databases details" && strings.Contains(d.Detail, "not found") {
			return true
		}
	}
	return false
}

// connectionPoolNotFound reports whether a Read failed because the pool is
// missing. Upstream reads it through the generated V9 client, which hands a
// 404 back as a response with that status rather than an error, so the
// IsNotFoundError branch never fires and the miss ends up as an error
// diagnostic.
func connectionPoolNotFound(diags []*tfprotov6.Diagnostic) bool {
	for _, d := range diags {
		if d.Summary == "Unable to read connection pool" && strings.Contains(d.Detail, "status code 404") {
			return true
		}
	}
	return false
}

// Configure configures the database group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("upcloud_managed_database_mysql", func(r *config.Resource) {
		r.UseAsync = true
		// Single-item blocks (SizeAtMost(1) upstream, not in the schema
		// dump): the properties block and its object-typed properties.
		for _, paths := range [][2]string{
			{"properties", "properties"},
			{"properties[*].migration", "properties[*].migration"},
			{"properties[*].mysql_incremental_backup", "properties[*].mysqlIncrementalBackup"},
		} {
			r.AddSingletonListConversion(paths[0], paths[1])
		}

		r.References["network.uuid"] = config.Reference{
			TerraformName: "upcloud_network",
		}

		// An omitted properties block is read back with empty lists
		// (ip_filter) and nothing else, which upjet's late-init turns into
		// an empty properties object; sending that back is rejected by the
		// API ("Minimum 1 properties allowed").
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"properties"},
		}
	})

	p.AddResourceConfigurator("upcloud_managed_database_opensearch", func(r *config.Resource) {
		r.UseAsync = true
		// Single-item blocks (SizeAtMost(1) upstream, not in the schema
		// dump): the properties block and its object-typed properties.
		for _, paths := range [][2]string{
			{"properties", "properties"},
			{"properties[*].auth_failure_listeners", "properties[*].authFailureListeners"},
			{"properties[*].auth_failure_listeners[*].internal_authentication_backend_limiting", "properties[*].authFailureListeners[*].internalAuthenticationBackendLimiting"},
			{"properties[*].cluster_remote_store", "properties[*].clusterRemoteStore"},
			{"properties[*].cluster_search_request_slowlog", "properties[*].clusterSearchRequestSlowlog"},
			{"properties[*].cluster_search_request_slowlog[*].threshold", "properties[*].clusterSearchRequestSlowlog[*].threshold"},
			{"properties[*].disk_watermarks", "properties[*].diskWatermarks"},
			{"properties[*].index_rollup", "properties[*].indexRollup"},
			{"properties[*].index_template", "properties[*].indexTemplate"},
			{"properties[*].jwt", "properties[*].jwt"},
			{"properties[*].openid", "properties[*].openid"},
			{"properties[*].opensearch_dashboards", "properties[*].opensearchDashboards"},
			{"properties[*].remote_store", "properties[*].remoteStore"},
			{"properties[*].saml", "properties[*].saml"},
			{"properties[*].search_backpressure", "properties[*].searchBackpressure"},
			{"properties[*].search_backpressure[*].node_duress", "properties[*].searchBackpressure[*].nodeDuress"},
			{"properties[*].search_backpressure[*].search_shard_task", "properties[*].searchBackpressure[*].searchShardTask"},
			{"properties[*].search_backpressure[*].search_task", "properties[*].searchBackpressure[*].searchTask"},
			{"properties[*].search_insights_top_queries", "properties[*].searchInsightsTopQueries"},
			{"properties[*].search_insights_top_queries[*].cpu", "properties[*].searchInsightsTopQueries[*].cpu"},
			{"properties[*].search_insights_top_queries[*].latency", "properties[*].searchInsightsTopQueries[*].latency"},
			{"properties[*].search_insights_top_queries[*].memory", "properties[*].searchInsightsTopQueries[*].memory"},
			{"properties[*].segrep", "properties[*].segrep"},
			{"properties[*].shard_indexing_pressure", "properties[*].shardIndexingPressure"},
			{"properties[*].shard_indexing_pressure[*].operating_factor", "properties[*].shardIndexingPressure[*].operatingFactor"},
			{"properties[*].shard_indexing_pressure[*].primary_parameter", "properties[*].shardIndexingPressure[*].primaryParameter"},
			{"properties[*].shard_indexing_pressure[*].primary_parameter[*].node", "properties[*].shardIndexingPressure[*].primaryParameter[*].node"},
			{"properties[*].shard_indexing_pressure[*].primary_parameter[*].shard", "properties[*].shardIndexingPressure[*].primaryParameter[*].shard"},
		} {
			r.AddSingletonListConversion(paths[0], paths[1])
		}

		r.References["network.uuid"] = config.Reference{
			TerraformName: "upcloud_network",
		}

		// An omitted properties block is read back with empty lists
		// (ip_filter) and nothing else, which upjet's late-init turns into
		// an empty properties object; sending that back is rejected by the
		// API ("Minimum 1 properties allowed").
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"properties"},
		}
	})

	p.AddResourceConfigurator("upcloud_managed_database_postgresql", func(r *config.Resource) {
		r.UseAsync = true
		// Single-item blocks (SizeAtMost(1) upstream, not in the schema
		// dump): the properties block and its object-typed properties.
		for _, paths := range [][2]string{
			{"properties", "properties"},
			{"properties[*].migration", "properties[*].migration"},
			{"properties[*].pgaudit", "properties[*].pgaudit"},
			{"properties[*].pgbouncer", "properties[*].pgbouncer"},
			{"properties[*].pglookout", "properties[*].pglookout"},
			{"properties[*].timescaledb", "properties[*].timescaledb"},
		} {
			r.AddSingletonListConversion(paths[0], paths[1])
		}

		r.References["network.uuid"] = config.Reference{
			TerraformName: "upcloud_network",
		}

		// An omitted properties block is read back with empty lists
		// (ip_filter) and nothing else, which upjet's late-init turns into
		// an empty properties object; sending that back is rejected by the
		// API ("Minimum 1 properties allowed").
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"properties"},
		}
	})

	p.AddResourceConfigurator("upcloud_managed_database_valkey", func(r *config.Resource) {
		r.UseAsync = true
		// Single-item blocks (SizeAtMost(1) upstream, not in the schema
		// dump): the properties block and its object-typed properties.
		for _, paths := range [][2]string{
			{"properties", "properties"},
			{"properties[*].migration", "properties[*].migration"},
		} {
			r.AddSingletonListConversion(paths[0], paths[1])
		}

		r.References["network.uuid"] = config.Reference{
			TerraformName: "upcloud_network",
		}

		// An omitted properties block is read back with empty lists
		// (ip_filter) and nothing else, which upjet's late-init turns into
		// an empty properties object; sending that back is rejected by the
		// API ("Minimum 1 properties allowed").
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"properties"},
		}
	})

	p.AddResourceConfigurator("upcloud_managed_database_logical_database", func(r *config.Resource) {
		r.ExternalName.IsNotFoundDiagnosticFn = logicalDatabaseNotFound

		// TODO: use a generic reference type if that is ever implemented in
		// upjet, see https://github.com/crossplane/upjet/issues/95
		r.References["service"] = config.Reference{
			TerraformName: "upcloud_managed_database_postgresql",
		}

		if s, ok := r.TerraformResource.Schema["service"]; ok {
			s.Description = serviceRefDescription
		}
	})

	p.AddResourceConfigurator("upcloud_managed_database_connection_pool", func(r *config.Resource) {
		r.ExternalName.IsNotFoundDiagnosticFn = connectionPoolNotFound

		r.References["service"] = config.Reference{
			TerraformName: "upcloud_managed_database_postgresql",
		}
		r.References["database"] = config.Reference{
			TerraformName: "upcloud_managed_database_logical_database",
			Extractor:     common.ExtractObservedExternalName,
		}
		r.References["username"] = config.Reference{
			TerraformName: "upcloud_managed_database_user",
			Extractor:     common.ExtractObservedExternalName,
		}

		if s, ok := r.TerraformResource.Schema["service"]; ok {
			s.Description = serviceRefDescription
		}
	})

	p.AddResourceConfigurator("upcloud_managed_database_user", func(r *config.Resource) {
		r.AddSingletonListConversion("pg_access_control", "pgAccessControl")
		r.AddSingletonListConversion("valkey_access_control", "valkeyAccessControl")
		r.AddSingletonListConversion("opensearch_access_control", "opensearchAccessControl")

		r.References["service"] = config.Reference{
			TerraformName: "upcloud_managed_database_postgresql",
		}

		if s, ok := r.TerraformResource.Schema["service"]; ok {
			s.Description = serviceRefDescription
		}
	})
}
