// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package common

import (
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	xpfake "github.com/crossplane/crossplane-runtime/v2/pkg/resource/fake"
	"github.com/crossplane/upjet/v2/pkg/resource/fake"
	"github.com/google/go-cmp/cmp"
)

func TestObservedExternalName(t *testing.T) {
	observed := &fake.Terraformed{Observable: fake.Observable{ID: "0a19c0ca-4a03-4cac-8705-be15f72245e7/lb-backend"}}
	meta.SetExternalName(observed, "lb-backend")
	unobserved := &fake.Terraformed{}
	meta.SetExternalName(unobserved, "lb-backend")
	notTerraformed := &xpfake.Managed{}
	meta.SetExternalName(notTerraformed, "lb-backend")

	cases := map[string]struct {
		args xpresource.Managed
		want string
	}{
		"Observed":       {args: observed, want: "lb-backend"},
		"NotYetObserved": {args: unobserved, want: ""},
		"NotTerraformed": {args: notTerraformed, want: ""},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := ObservedExternalName()(tc.args)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("ObservedExternalName(): -want, +got:\n%s", diff)
			}
		})
	}
}
