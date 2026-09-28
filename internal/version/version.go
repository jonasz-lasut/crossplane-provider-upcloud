// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

// Package version contains the version of this provider
package version

import (
	"fmt"
	"strings"
)

// Version will be overridden with the current version at build time using the -X linker flag
var Version = "0.0.0"

// UserAgent returns the User-Agent this provider reports to the UpCloud API.
// The crossplane-provider-upcloud prefix predates the move to
// crossplane-contrib and is kept so that UpCloud can keep attributing API
// traffic to this provider.
func UserAgent() string {
	return fmt.Sprintf("crossplane-provider-upcloud/%s", strings.TrimPrefix(Version, "v"))
}
