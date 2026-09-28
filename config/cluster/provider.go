// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package cluster

import (
	"github.com/crossplane-contrib/provider-upcloud/config/cluster/database"
	"github.com/crossplane-contrib/provider-upcloud/config/cluster/network"
	"github.com/crossplane-contrib/provider-upcloud/config/cluster/objectstorage"
	"github.com/crossplane-contrib/provider-upcloud/config/cluster/server"
	"github.com/crossplane-contrib/provider-upcloud/config/cluster/storage"
	"github.com/crossplane-contrib/provider-upcloud/config/cluster/uks"
)

func init() {
	ProviderConfiguration.AddConfig(database.Configure)
	ProviderConfiguration.AddConfig(network.Configure)
	ProviderConfiguration.AddConfig(objectstorage.Configure)
	ProviderConfiguration.AddConfig(server.Configure)
	ProviderConfiguration.AddConfig(storage.Configure)
	ProviderConfiguration.AddConfig(uks.Configure)
}
