// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	manageddatabaseconnectionpool "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/database/manageddatabaseconnectionpool"
	manageddatabaselogicaldatabase "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/database/manageddatabaselogicaldatabase"
	manageddatabasemysql "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/database/manageddatabasemysql"
	manageddatabaseopensearch "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/database/manageddatabaseopensearch"
	manageddatabasepostgresql "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/database/manageddatabasepostgresql"
	manageddatabaseuser "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/database/manageddatabaseuser"
	manageddatabasevalkey "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/database/manageddatabasevalkey"
	firewallruleset "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/firewallruleset"
	floatingipaddress "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/floatingipaddress"
	gateway "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/gateway"
	gatewayconnection "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/gatewayconnection"
	gatewayconnectiontunnel "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/gatewayconnectiontunnel"
	loadbalancer "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/loadbalancer"
	loadbalancerbackend "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/loadbalancerbackend"
	loadbalancerbackendtlsconfig "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/loadbalancerbackendtlsconfig"
	loadbalancerdynamicbackendmember "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/loadbalancerdynamicbackendmember"
	loadbalancerdynamiccertificatebundle "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/loadbalancerdynamiccertificatebundle"
	loadbalancerfrontend "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/loadbalancerfrontend"
	loadbalancerfrontendrule "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/loadbalancerfrontendrule"
	loadbalancerfrontendtlsconfig "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/loadbalancerfrontendtlsconfig"
	loadbalancermanualcertificatebundle "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/loadbalancermanualcertificatebundle"
	loadbalancerresolver "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/loadbalancerresolver"
	loadbalancerstaticbackendmember "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/loadbalancerstaticbackendmember"
	network "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/network"
	networkpeering "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/networkpeering"
	router "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/network/router"
	managedobjectstorage "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/objectstorage/managedobjectstorage"
	managedobjectstoragebucket "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/objectstorage/managedobjectstoragebucket"
	managedobjectstoragecustomdomain "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/objectstorage/managedobjectstoragecustomdomain"
	managedobjectstoragepolicy "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/objectstorage/managedobjectstoragepolicy"
	managedobjectstoragestaticsite "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/objectstorage/managedobjectstoragestaticsite"
	managedobjectstorageuser "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/objectstorage/managedobjectstorageuser"
	managedobjectstorageuseraccesskey "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/objectstorage/managedobjectstorageuseraccesskey"
	managedobjectstorageuserpolicy "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/objectstorage/managedobjectstorageuserpolicy"
	providerconfig "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/providerconfig"
	firewallrules "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/server/firewallrules"
	server "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/server/server"
	servergroup "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/server/servergroup"
	serverprivatefirewallruleset "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/server/serverprivatefirewallruleset"
	tag "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/server/tag"
	filestorage "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/storage/filestorage"
	filestorageshare "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/storage/filestorageshare"
	filestorageshareacl "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/storage/filestorageshareacl"
	storage "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/storage/storage"
	storagebackup "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/storage/storagebackup"
	storagetemplate "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/storage/storagetemplate"
	kubernetescluster "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/uks/kubernetescluster"
	kubernetesnodegroup "github.com/crossplane-contrib/provider-upcloud/internal/controller/namespaced/uks/kubernetesnodegroup"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		manageddatabaseconnectionpool.Setup,
		manageddatabaselogicaldatabase.Setup,
		manageddatabasemysql.Setup,
		manageddatabaseopensearch.Setup,
		manageddatabasepostgresql.Setup,
		manageddatabaseuser.Setup,
		manageddatabasevalkey.Setup,
		firewallruleset.Setup,
		floatingipaddress.Setup,
		gateway.Setup,
		gatewayconnection.Setup,
		gatewayconnectiontunnel.Setup,
		loadbalancer.Setup,
		loadbalancerbackend.Setup,
		loadbalancerbackendtlsconfig.Setup,
		loadbalancerdynamicbackendmember.Setup,
		loadbalancerdynamiccertificatebundle.Setup,
		loadbalancerfrontend.Setup,
		loadbalancerfrontendrule.Setup,
		loadbalancerfrontendtlsconfig.Setup,
		loadbalancermanualcertificatebundle.Setup,
		loadbalancerresolver.Setup,
		loadbalancerstaticbackendmember.Setup,
		network.Setup,
		networkpeering.Setup,
		router.Setup,
		managedobjectstorage.Setup,
		managedobjectstoragebucket.Setup,
		managedobjectstoragecustomdomain.Setup,
		managedobjectstoragepolicy.Setup,
		managedobjectstoragestaticsite.Setup,
		managedobjectstorageuser.Setup,
		managedobjectstorageuseraccesskey.Setup,
		managedobjectstorageuserpolicy.Setup,
		providerconfig.Setup,
		firewallrules.Setup,
		server.Setup,
		servergroup.Setup,
		serverprivatefirewallruleset.Setup,
		tag.Setup,
		filestorage.Setup,
		filestorageshare.Setup,
		filestorageshareacl.Setup,
		storage.Setup,
		storagebackup.Setup,
		storagetemplate.Setup,
		kubernetescluster.Setup,
		kubernetesnodegroup.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		manageddatabaseconnectionpool.SetupGated,
		manageddatabaselogicaldatabase.SetupGated,
		manageddatabasemysql.SetupGated,
		manageddatabaseopensearch.SetupGated,
		manageddatabasepostgresql.SetupGated,
		manageddatabaseuser.SetupGated,
		manageddatabasevalkey.SetupGated,
		firewallruleset.SetupGated,
		floatingipaddress.SetupGated,
		gateway.SetupGated,
		gatewayconnection.SetupGated,
		gatewayconnectiontunnel.SetupGated,
		loadbalancer.SetupGated,
		loadbalancerbackend.SetupGated,
		loadbalancerbackendtlsconfig.SetupGated,
		loadbalancerdynamicbackendmember.SetupGated,
		loadbalancerdynamiccertificatebundle.SetupGated,
		loadbalancerfrontend.SetupGated,
		loadbalancerfrontendrule.SetupGated,
		loadbalancerfrontendtlsconfig.SetupGated,
		loadbalancermanualcertificatebundle.SetupGated,
		loadbalancerresolver.SetupGated,
		loadbalancerstaticbackendmember.SetupGated,
		network.SetupGated,
		networkpeering.SetupGated,
		router.SetupGated,
		managedobjectstorage.SetupGated,
		managedobjectstoragebucket.SetupGated,
		managedobjectstoragecustomdomain.SetupGated,
		managedobjectstoragepolicy.SetupGated,
		managedobjectstoragestaticsite.SetupGated,
		managedobjectstorageuser.SetupGated,
		managedobjectstorageuseraccesskey.SetupGated,
		managedobjectstorageuserpolicy.SetupGated,
		providerconfig.SetupGated,
		firewallrules.SetupGated,
		server.SetupGated,
		servergroup.SetupGated,
		serverprivatefirewallruleset.SetupGated,
		tag.SetupGated,
		filestorage.SetupGated,
		filestorageshare.SetupGated,
		filestorageshareacl.SetupGated,
		storage.SetupGated,
		storagebackup.SetupGated,
		storagetemplate.SetupGated,
		kubernetescluster.SetupGated,
		kubernetesnodegroup.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupWebhookWithManager registers conversion webhooks for all resource kinds in the group.
func SetupWebhookWithManager(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		manageddatabaseconnectionpool.SetupWebhookWithManager,
		manageddatabaselogicaldatabase.SetupWebhookWithManager,
		manageddatabasemysql.SetupWebhookWithManager,
		manageddatabaseopensearch.SetupWebhookWithManager,
		manageddatabasepostgresql.SetupWebhookWithManager,
		manageddatabaseuser.SetupWebhookWithManager,
		manageddatabasevalkey.SetupWebhookWithManager,
		firewallruleset.SetupWebhookWithManager,
		floatingipaddress.SetupWebhookWithManager,
		gateway.SetupWebhookWithManager,
		gatewayconnection.SetupWebhookWithManager,
		gatewayconnectiontunnel.SetupWebhookWithManager,
		loadbalancer.SetupWebhookWithManager,
		loadbalancerbackend.SetupWebhookWithManager,
		loadbalancerbackendtlsconfig.SetupWebhookWithManager,
		loadbalancerdynamicbackendmember.SetupWebhookWithManager,
		loadbalancerdynamiccertificatebundle.SetupWebhookWithManager,
		loadbalancerfrontend.SetupWebhookWithManager,
		loadbalancerfrontendrule.SetupWebhookWithManager,
		loadbalancerfrontendtlsconfig.SetupWebhookWithManager,
		loadbalancermanualcertificatebundle.SetupWebhookWithManager,
		loadbalancerresolver.SetupWebhookWithManager,
		loadbalancerstaticbackendmember.SetupWebhookWithManager,
		network.SetupWebhookWithManager,
		networkpeering.SetupWebhookWithManager,
		router.SetupWebhookWithManager,
		managedobjectstorage.SetupWebhookWithManager,
		managedobjectstoragebucket.SetupWebhookWithManager,
		managedobjectstoragecustomdomain.SetupWebhookWithManager,
		managedobjectstoragepolicy.SetupWebhookWithManager,
		managedobjectstoragestaticsite.SetupWebhookWithManager,
		managedobjectstorageuser.SetupWebhookWithManager,
		managedobjectstorageuseraccesskey.SetupWebhookWithManager,
		managedobjectstorageuserpolicy.SetupWebhookWithManager,
		providerconfig.SetupWebhookWithManager,
		firewallrules.SetupWebhookWithManager,
		server.SetupWebhookWithManager,
		servergroup.SetupWebhookWithManager,
		serverprivatefirewallruleset.SetupWebhookWithManager,
		tag.SetupWebhookWithManager,
		filestorage.SetupWebhookWithManager,
		filestorageshare.SetupWebhookWithManager,
		filestorageshareacl.SetupWebhookWithManager,
		storage.SetupWebhookWithManager,
		storagebackup.SetupWebhookWithManager,
		storagetemplate.SetupWebhookWithManager,
		kubernetescluster.SetupWebhookWithManager,
		kubernetesnodegroup.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
