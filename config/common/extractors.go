// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package common

import (
	"fmt"

	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reference"
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/resource"
)

// ExtractObservedExternalName is the reference extractor expression that
// resolves a reference to the target's external name once the target has been
// observed. It is the right extractor for the name-keyed child resources whose
// Terraform ID is <parent>/<name>: their external name is the name itself,
// while their ID carries the parent prefix.
const ExtractObservedExternalName = `github.com/crossplane-contrib/provider-upcloud/config/common.ObservedExternalName()`

// ExtractObservedPath returns the reference extractor expression that resolves
// a reference to the value at the given Terraform attribute path of the
// target's observed state (status.atProvider).
func ExtractObservedPath(path string) string {
	return fmt.Sprintf(`github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath(%q,true)`, path)
}

// ObservedExternalName returns the target's external name once the target has
// been observed, that is once status.atProvider.id is set. Name-keyed
// resources receive their external name from the name initializer on their
// first reconcile, before they exist in UpCloud; resolving through the plain
// external name would let the referencing resource call the API with a name
// that is not there yet.
func ObservedExternalName() reference.ExtractValueFn {
	return func(mg xpresource.Managed) string {
		tr, ok := mg.(resource.Terraformed)
		if !ok || tr.GetID() == "" {
			return ""
		}
		return meta.GetExternalName(mg)
	}
}
