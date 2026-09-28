// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package clients

import (
	"strings"
	"testing"

	xpv2 "github.com/crossplane/crossplane/apis/v2/core/v2"
	"github.com/google/go-cmp/cmp"

	namespacedv1beta1 "github.com/crossplane-contrib/provider-upcloud/apis/namespaced/v1beta1"
)

func TestResolveNamespacedSpec(t *testing.T) {
	type args struct {
		spec      namespacedv1beta1.NamespacedProviderConfigSpec
		namespace string
	}

	cases := map[string]struct {
		args args
		want namespacedv1beta1.ProviderConfigSpec
	}{
		"SecretRefResolvesToReferencerNamespace": {
			args: args{
				namespace: "team-a",
				spec: namespacedv1beta1.NamespacedProviderConfigSpec{
					Credentials: namespacedv1beta1.NamespacedProviderCredentials{
						Source: xpv2.CredentialsSourceSecret,
						SecretRef: &xpv2.LocalSecretKeySelector{
							LocalSecretReference: xpv2.LocalSecretReference{Name: "creds"},
							Key:                  "credentials",
						},
					},
				},
			},
			want: namespacedv1beta1.ProviderConfigSpec{
				Credentials: namespacedv1beta1.ProviderCredentials{
					Source: xpv2.CredentialsSourceSecret,
					CommonCredentialSelectors: xpv2.CommonCredentialSelectors{
						SecretRef: &xpv2.SecretKeySelector{
							SecretReference: xpv2.SecretReference{Name: "creds", Namespace: "team-a"},
							Key:             "credentials",
						},
					},
				},
			},
		},
		"NilSecretRefStaysNil": {
			args: args{
				namespace: "team-a",
				spec: namespacedv1beta1.NamespacedProviderConfigSpec{
					Credentials: namespacedv1beta1.NamespacedProviderCredentials{
						Source: xpv2.CredentialsSourceInjectedIdentity,
					},
				},
			},
			want: namespacedv1beta1.ProviderConfigSpec{
				Credentials: namespacedv1beta1.ProviderCredentials{
					Source: xpv2.CredentialsSourceInjectedIdentity,
				},
			},
		},
		"FsAndEnvPassThrough": {
			args: args{
				namespace: "team-a",
				spec: namespacedv1beta1.NamespacedProviderConfigSpec{
					Credentials: namespacedv1beta1.NamespacedProviderCredentials{
						Source: xpv2.CredentialsSourceFilesystem,
						Fs:     &xpv2.FsSelector{Path: "/creds"},
						Env:    &xpv2.EnvSelector{Name: "UPCLOUD_CREDS"},
					},
				},
			},
			want: namespacedv1beta1.ProviderConfigSpec{
				Credentials: namespacedv1beta1.ProviderCredentials{
					Source: xpv2.CredentialsSourceFilesystem,
					CommonCredentialSelectors: xpv2.CommonCredentialSelectors{
						Fs:  &xpv2.FsSelector{Path: "/creds"},
						Env: &xpv2.EnvSelector{Name: "UPCLOUD_CREDS"},
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := resolveNamespacedSpec(tc.args.spec, tc.args.namespace)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("resolveNamespacedSpec() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestBuildConfiguration(t *testing.T) {
	type args struct {
		creds map[string]any
	}

	cases := map[string]struct {
		args            args
		want            map[string]any
		wantErrContains string
	}{
		"TokenOnly": {
			args: args{creds: map[string]any{"token": "ucat_token"}},
			want: map[string]any{"token": "ucat_token"},
		},
		"TokenWinsOverUsernameAndPassword": {
			args: args{creds: map[string]any{"token": "ucat_token", "username": "user", "password": "pass"}},
			want: map[string]any{"token": "ucat_token"},
		},
		"UsernameAndPassword": {
			args: args{creds: map[string]any{"username": "user", "password": "pass"}},
			want: map[string]any{"username": "user", "password": "pass"},
		},
		"ExtraKeysIgnored": {
			args: args{creds: map[string]any{"token": "ucat_token", "unrelated": "x"}},
			want: map[string]any{"token": "ucat_token"},
		},
		"EmptyTokenFallsBackToUsernameAndPassword": {
			args: args{creds: map[string]any{"token": "", "username": "user", "password": "pass"}},
			want: map[string]any{"username": "user", "password": "pass"},
		},
		"UsernameWithoutPassword": {
			args:            args{creds: map[string]any{"username": "user"}},
			wantErrContains: `needs a "token" key or both "username" and "password" keys`,
		},
		"NoCredentials": {
			args:            args{creds: map[string]any{}},
			wantErrContains: `needs a "token" key or both "username" and "password" keys`,
		},
		"TokenNotAString": {
			args:            args{creds: map[string]any{"token": 42.0}},
			wantErrContains: `credentials key "token" must be a string`,
		},
		"ClientSettingsAsNumbers": {
			args: args{creds: map[string]any{"token": "ucat_token", "request_timeout_sec": 30.0, "retry_max": 2.0, "retry_wait_min_sec": 1.0, "retry_wait_max_sec": 10.0}},
			want: map[string]any{"token": "ucat_token", "request_timeout_sec": int64(30), "retry_max": int64(2), "retry_wait_min_sec": int64(1), "retry_wait_max_sec": int64(10)},
		},
		"ClientSettingsAsNumericStrings": {
			args: args{creds: map[string]any{"username": "user", "password": "pass", "request_timeout_sec": "30", "retry_max": "2"}},
			want: map[string]any{"username": "user", "password": "pass", "request_timeout_sec": int64(30), "retry_max": int64(2)},
		},
		"ClientSettingNotAnInteger": {
			args:            args{creds: map[string]any{"token": "ucat_token", "request_timeout_sec": 2.5}},
			wantErrContains: `credentials key "request_timeout_sec" must be an integer`,
		},
		"ClientSettingNotNumeric": {
			args:            args{creds: map[string]any{"token": "ucat_token", "retry_max": "many"}},
			wantErrContains: `credentials key "retry_max" must be an integer`,
		},
		"ClientSettingWrongType": {
			args:            args{creds: map[string]any{"token": "ucat_token", "retry_max": true}},
			wantErrContains: `credentials key "retry_max" must be an integer, got bool`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := buildConfiguration(tc.args.creds)
			if tc.wantErrContains != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErrContains)
				}
				if !strings.Contains(err.Error(), tc.wantErrContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tc.wantErrContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("buildConfiguration() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
