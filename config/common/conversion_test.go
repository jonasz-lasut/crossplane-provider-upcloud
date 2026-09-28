// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package common

import (
	"testing"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/google/go-cmp/cmp"
)

func TestEmptyValueDefaults(t *testing.T) {
	type args struct {
		params map[string]any
		mode   config.Mode
	}

	cases := map[string]struct {
		args args
		want map[string]any
	}{
		"AbsentArgumentsGetEmptyValues": {
			args: args{params: map[string]any{"name": "router"}, mode: config.ToTerraform},
			want: map[string]any{"name": "router", "labels": map[string]any{}, "static_route": []any{}},
		},
		"PresentArgumentsAreKept": {
			args: args{
				params: map[string]any{"labels": map[string]any{"env": "dev"}, "static_route": []any{map[string]any{"name": "r"}}},
				mode:   config.ToTerraform,
			},
			want: map[string]any{"labels": map[string]any{"env": "dev"}, "static_route": []any{map[string]any{"name": "r"}}},
		},
		"ExplicitNilIsKept": {
			args: args{params: map[string]any{"labels": nil}, mode: config.ToTerraform},
			want: map[string]any{"labels": nil, "static_route": []any{}},
		},
		"FromTerraformIsUntouched": {
			args: args{params: map[string]any{"name": "router"}, mode: config.FromTerraform},
			want: map[string]any{"name": "router"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := EmptyValueDefaults([]string{"labels"}, []string{"static_route"}).Convert(tc.args.params, nil, tc.args.mode)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Convert() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
