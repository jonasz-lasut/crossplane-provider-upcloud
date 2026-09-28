// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package namespaced

import (
	"github.com/crossplane-contrib/provider-upcloud/config/namespaced/database"
	"github.com/crossplane-contrib/provider-upcloud/config/namespaced/network"
	"github.com/crossplane-contrib/provider-upcloud/config/namespaced/objectstorage"
	"github.com/crossplane-contrib/provider-upcloud/config/namespaced/server"
	"github.com/crossplane-contrib/provider-upcloud/config/namespaced/storage"
	"github.com/crossplane-contrib/provider-upcloud/config/namespaced/uks"
)

func init() {
	ProviderConfiguration.AddConfig(database.Configure)
	ProviderConfiguration.AddConfig(network.Configure)
	ProviderConfiguration.AddConfig(objectstorage.Configure)
	ProviderConfiguration.AddConfig(server.Configure)
	ProviderConfiguration.AddConfig(storage.Configure)
	ProviderConfiguration.AddConfig(uks.Configure)
}
