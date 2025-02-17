/*
 * SPDX-FileCopyrightText: 2019 SAP SE or an SAP affiliate company and Gardener contributors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package metrics

import (
	"github.com/gardener/controller-manager-library/pkg/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
)

func init() {
	prometheus.MustRegister(ACMEAccountRegistrations)
	prometheus.MustRegister(ACMETotalOrders)
	prometheus.MustRegister(ACMEActiveDNSChallenges)
	prometheus.MustRegister(CertEntries)
	prometheus.MustRegister(OverdueCertificates)
	prometheus.MustRegister(RevokedCertificates)
	prometheus.MustRegister(CertificateSecrets)

	server.RegisterHandler("/metrics", promhttp.Handler())
}

var (
	// ACMEAccountRegistrations is the cert_management_acme_account_registrations gauge.
	ACMEAccountRegistrations = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cert_management_acme_account_registrations",
			Help: "ACME account registrations",
		},
		[]string{"uri", "email"},
	)

	// ACMETotalOrders is the cert_management_acme_orders counter.
	ACMETotalOrders = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cert_management_acme_orders",
			Help: "Number of ACME orders",
		},
		[]string{"issuer", "success", "dns_challenges", "renew"},
	)

	// ACMEActiveDNSChallenges is the cert_management_acme_active_dns_challenges gauge.
	ACMEActiveDNSChallenges = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cert_management_acme_active_dns_challenges",
			Help: "Currently active number of ACME DNS challenges",
		},
		[]string{"issuer"},
	)

	// CertEntries is the cert_management_cert_entries gauge.
	CertEntries = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cert_management_cert_entries",
			Help: "Total number of certificate objects per issuer",
		},
		[]string{"issuertype", "issuer"},
	)

	// OverdueCertificates is the cert_management_overdue_renewal_certificates gauge.
	OverdueCertificates = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cert_management_overdue_renewal_certificates",
			Help: "Number of certificate objects with certificate's renewal overdue",
		},
	)

	// RevokedCertificates is the cert_management_revoked_certificates gauge.
	RevokedCertificates = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cert_management_revoked_certificates",
			Help: "Number of certificate objects with revoked certificate",
		},
	)

	// CertificateSecrets is the cert_management_secrets gauge.
	CertificateSecrets = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cert_management_secrets",
			Help: "Number of certificate secrets per classification",
		},
		[]string{"classification"},
	)
)

// AddACMEAccountRegistration increments the ACMEAccountRegistrations counter.
func AddACMEAccountRegistration(uri, email string) {
	ACMEAccountRegistrations.WithLabelValues(uri, email).Set(1)
}

// AddACMEOrder increments the ACMETotalOrders counter.
func AddACMEOrder(issuer string, success bool, count int, renew bool) {
	if count > 0 {
		ACMETotalOrders.WithLabelValues(issuer, strconv.FormatBool(success), strconv.FormatInt(int64(count), 10), strconv.FormatBool(renew)).Inc()
	}
}

// AddActiveACMEDNSChallenge increments the ACMEActiveDNSChallenges gauge.
func AddActiveACMEDNSChallenge(issuer string) {
	ACMEActiveDNSChallenges.WithLabelValues(issuer).Inc()
}

// RemoveActiveACMEDNSChallenge decrements the ACMEActiveDNSChallenges gauge.
func RemoveActiveACMEDNSChallenge(issuer string) {
	ACMEActiveDNSChallenges.WithLabelValues(issuer).Dec()
}

// ReportCertEntries sets the CertEntries gauge
func ReportCertEntries(issuertype, issuer string, count int) {
	CertEntries.WithLabelValues(issuertype, issuer).Set(float64(count))
}

// DeleteCertEntries deletes a CertEntries gauge entry.
func DeleteCertEntries(issuertype, issuer string) {
	CertEntries.DeleteLabelValues(issuertype, issuer)
}

// ReportOverdueCerts sets the OverdueCertificates gauge
func ReportOverdueCerts(count int) {
	OverdueCertificates.Set(float64(count))
}

// ReportRevokedCerts sets the RevokedCertificates gauge
func ReportRevokedCerts(count int) {
	RevokedCertificates.Set(float64(count))
}

// ReportCertificateSecrets sets the CertificateSecrets gauge
func ReportCertificateSecrets(classification string, count int) {
	CertificateSecrets.WithLabelValues(classification).Set(float64(count))
}
