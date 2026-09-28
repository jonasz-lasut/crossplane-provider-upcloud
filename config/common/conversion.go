// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

// Package common holds configuration helpers shared by the cluster-scoped
// and namespaced resource configurations.
package common

import (
	"maps"
	"slices"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// emptyValueConversion sends an explicit empty value for the given Terraform
// arguments when the managed resource omits them. Upjet keeps omitempty on
// top-level required parameters, so an empty map or list set in the spec is
// dropped from the Terraform configuration; the UpCloud plugin-framework
// resources treat the resulting null as unknown and never converge. The
// upstream provider used a forked upjet to strip omitempty for these fields;
// this conversion has the same effect without the fork.
type emptyValueConversion struct {
	defaults map[string]func() any
}

// EmptyValueDefaults returns a TerraformConversion that adds the given
// arguments with an empty value of their kind whenever they are absent from
// the parameters handed to Terraform.
func EmptyValueDefaults(emptyMaps []string, emptyLists []string) config.TerraformConversion {
	defaults := make(map[string]func() any, len(emptyMaps)+len(emptyLists))
	for _, name := range emptyMaps {
		defaults[name] = func() any { return map[string]any{} }
	}
	for _, name := range emptyLists {
		defaults[name] = func() any { return []any{} }
	}
	return emptyValueConversion{defaults: defaults}
}

func (c emptyValueConversion) Convert(params map[string]any, _ *config.Resource, mode config.Mode) (map[string]any, error) {
	if mode != config.ToTerraform {
		return params, nil
	}
	for _, name := range slices.Sorted(maps.Keys(c.defaults)) {
		if _, set := params[name]; !set {
			params[name] = c.defaults[name]()
		}
	}
	return params, nil
}
