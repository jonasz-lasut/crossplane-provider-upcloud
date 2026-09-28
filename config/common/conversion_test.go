// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package common

import (
	"testing"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/google/go-cmp/cmp"
)

func TestEmptyListDefaults(t *testing.T) {
	type args struct {
		params map[string]any
		mode   config.Mode
	}

	cases := map[string]struct {
		args args
		want map[string]any
	}{
		"AbsentArgumentGetsEmptyList": {
			args: args{params: map[string]any{"plan": "2xCPU-4GB"}, mode: config.ToTerraform},
			want: map[string]any{"plan": "2xCPU-4GB", "ssh_keys": []any{}},
		},
		"PresentArgumentIsKept": {
			args: args{params: map[string]any{"ssh_keys": []any{"ssh-ed25519 AAAA"}}, mode: config.ToTerraform},
			want: map[string]any{"ssh_keys": []any{"ssh-ed25519 AAAA"}},
		},
		"ExplicitNilIsKept": {
			args: args{params: map[string]any{"ssh_keys": nil}, mode: config.ToTerraform},
			want: map[string]any{"ssh_keys": nil},
		},
		"FromTerraformIsUntouched": {
			args: args{params: map[string]any{"plan": "2xCPU-4GB"}, mode: config.FromTerraform},
			want: map[string]any{"plan": "2xCPU-4GB"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := EmptyListDefaults("ssh_keys").Convert(tc.args.params, nil, tc.args.mode)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Convert() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
