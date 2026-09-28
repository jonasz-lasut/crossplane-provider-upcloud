// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package common

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// emptyListConversion sends an explicit empty list for the given Terraform
// arguments whenever the managed resource omits them.
//
// Upjet keeps omitempty on every generated field, so an empty set never
// survives the round trip through status.atProvider: after a provider restart
// the prior state is rebuilt with the argument as null while a plugin-framework
// resource plans its empty default, and upjet's replacement filter compares the
// two raw values and refuses every update as a replacement
// (upcloud_kubernetes_node_group ssh_keys). Upjet runs the conversion on the
// configuration and on the rebuilt prior state, so both sides carry the same
// empty list.
//
// Use it only for attributes whose upstream Read tolerates an empty non-null
// value, and never for singleton-list blocks: the singleton conversion runs
// afterwards and would wrap the empty list into a one-element list.
type emptyListConversion struct {
	names []string
}

// EmptyListDefaults returns a TerraformConversion that adds the given
// arguments as empty lists whenever they are absent from the values handed to
// Terraform.
func EmptyListDefaults(names ...string) config.TerraformConversion {
	return emptyListConversion{names: names}
}

func (c emptyListConversion) Convert(params map[string]any, _ *config.Resource, mode config.Mode) (map[string]any, error) {
	if mode != config.ToTerraform {
		return params, nil
	}
	for _, name := range c.names {
		if _, set := params[name]; !set {
			params[name] = []any{}
		}
	}
	return params, nil
}
