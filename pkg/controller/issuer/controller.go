/*
 * SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package issuer

import (
	"time"

	"github.com/gardener/controller-manager-library/pkg/controllermanager/controller"
	"github.com/gardener/controller-manager-library/pkg/resources"
	"github.com/gardener/controller-manager-library/pkg/resources/apiextensions"

	"github.com/gardener/cert-management/pkg/apis/cert"
	"github.com/gardener/cert-management/pkg/apis/cert/crds"
	api "github.com/gardener/cert-management/pkg/apis/cert/v1alpha1"
	"github.com/gardener/cert-management/pkg/cert/source"
	ctrl "github.com/gardener/cert-management/pkg/controller"
	"github.com/gardener/cert-management/pkg/controller/issuer/core"
)

var certificateGroupKind = resources.NewGroupKind(api.GroupName, api.CertificateKind)
var issuerGroupKind = resources.NewGroupKind(api.GroupName, api.IssuerKind)
var certificateRevocationGroupKind = resources.NewGroupKind(api.GroupName, api.CertificateRevocationKind)

func init() {
	crds.AddToRegistry(apiextensions.DefaultRegistry())

	controller.Configure("issuer").
		DefaultedStringOption(core.OptDefaultIssuer, "default-issuer", "name of default issuer (from default cluster)").
		DefaultedStringOption(core.OptIssuerNamespace, "default", "namespace to lookup issuers on default cluster").
		StringOption(core.OptDefaultIssuerDomainRanges, "domain range restrictions when using default issuer separated by comma").
		StringOption(core.OptDNSNamespace, "namespace for creating challenge DNSEntries (in DNS cluster)").
		StringOption(core.OptDNSClass, "class for creating challenge DNSEntries (in DNS cluster)").
		StringOption(core.OptDNSOwnerID, "ownerId for creating challenge DNSEntries").
		BoolOption(core.OptCascadeDelete, "If true, certificate secrets are deleted if dependent resources (certificate, ingress) are deleted").
		StringOption(source.OptClass, "Identifier used to differentiate responsible controllers for entries").
		DefaultedDurationOption(core.OptRenewalWindow, 30*24*time.Hour, "certificate is renewed if its validity period is shorter").
		DefaultedDurationOption(core.OptRenewalOverdueWindow, 25*24*time.Hour, "certificate is counted as 'renewal overdue' if its validity period is shorter (metrics cert_management_overdue_renewal_certificates)").
		DefaultedStringOption(core.OptPrecheckNameservers, "8.8.8.8:53,8.8.4.4:53",
			"DNS nameservers used for checking DNS propagation. If explicity set empty, it is tried to read them from /etc/resolv.conf").
		DefaultedDurationOption(core.OptPrecheckAdditionalWait, 10*time.Second, "additional wait time after DNS propagation check").
		DefaultedDurationOption(core.OptPropagationTimeout, 120*time.Second, "propagation timeout for DNS challenge").
		DefaultedIntOption(core.OptDefaultRequestsPerDayQuota, 10000,
			"Default value for requestsPerDayQuota if not set explicitly in the issuer spec.").
		FinalizerDomain(cert.GroupName).
		Cluster(ctrl.TargetCluster).
		DefaultWorkerPool(2, 24*time.Hour).
		MainResource(api.GroupName, api.CertificateKind).
		WorkerPool("revocations", 1, 0).
		Watch(api.GroupName, api.CertificateRevocationKind).
		Reconciler(newCompoundReconciler).
		Cluster(ctrl.DefaultCluster).
		WorkerPool("issuers", 1, 0).
		SelectedWatch(selectIssuerNamespaceSelectionFunction, api.GroupName, api.IssuerKind).
		WorkerPool("secrets", 1, 0).
		SelectedWatch(selectIssuerNamespaceSelectionFunction, "core", "Secret").
		Cluster(ctrl.DNSCluster).
		Cluster(ctrl.SourceCluster).
		RequireLease(ctrl.SourceCluster).
		MustRegister(ctrl.ControllerGroupCert)
}

func selectIssuerNamespaceSelectionFunction(c controller.Interface) (string, resources.TweakListOptionsFunc) {
	var options resources.TweakListOptionsFunc
	issuerNamespace, _ := c.GetStringOption(core.OptIssuerNamespace)
	return issuerNamespace, options
}
