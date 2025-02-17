/*
 * SPDX-FileCopyrightText: 2019 SAP SE or an SAP affiliate company and Gardener contributors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package core

import (
	"github.com/gardener/controller-manager-library/pkg/resources"
	v1 "k8s.io/api/core/v1"
)

type state struct {
	secrets      ReferencedSecrets
	eabSecrets   ReferencedSecrets
	certificates AssociatedObjects
	quotas       Quotas
	overdueCerts objectNameSet
	revokedCerts objectNameSet
}

func newState() *state {
	return &state{secrets: *NewReferencedSecrets(), eabSecrets: *NewReferencedSecrets(),
		certificates: *NewAssociatedObjects(), quotas: *NewQuotas(),
		overdueCerts: *newObjectNameSet(), revokedCerts: *newObjectNameSet()}
}

func (s *state) RemoveIssuer(name resources.ObjectName) bool {
	s.certificates.RemoveBySource(name)
	s.quotas.RemoveIssuer(name)
	s.eabSecrets.RemoveIssuer(name)
	return s.secrets.RemoveIssuer(name)
}

func (s *state) AddCertAssoc(issuer resources.ObjectName, cert resources.ObjectName) {
	s.certificates.AddAssoc(issuer, cert)
}

func (s *state) RemoveCertAssoc(cert resources.ObjectName) {
	s.certificates.RemoveByDest(cert)
}

func (s *state) CertificateNamesForIssuer(issuer resources.ObjectName) []resources.ObjectName {
	return s.certificates.DestinationsAsArray(issuer)
}

func (s *state) CertificateCountForIssuer(issuer resources.ObjectName) int {
	return s.certificates.DestinationsCount(issuer)
}

func (s *state) KnownIssuers() []resources.ObjectName {
	return s.certificates.Sources()
}

func (s *state) RememberIssuerQuotas(issuer resources.ObjectName, requestsPerDay int) {
	s.quotas.RememberQuotas(issuer, requestsPerDay)
}

// TryAcceptCertificateRequest tries to accept a certificate request according to the quotas.
// Return true if accepted and the requests per days quota value
func (s *state) TryAcceptCertificateRequest(issuer resources.ObjectName) (bool, int) {
	return s.quotas.TryAccept(issuer)
}

func (s *state) IssuerNamesForSecret(secretName resources.ObjectName) resources.ObjectNameSet {
	return s.secrets.IssuerNamesFor(secretName)
}

func (s *state) RememberIssuerSecret(issuer resources.ObjectName, secretRef *v1.SecretReference, hash string) {
	s.secrets.RememberIssuerSecret(issuer, secretRef, hash)
}

func (s *state) GetIssuerSecretHash(issuerName resources.ObjectName) string {
	return s.secrets.GetIssuerSecretHash(issuerName)
}

func (s *state) IssuerNamesForEABSecret(secretName resources.ObjectName) resources.ObjectNameSet {
	return s.eabSecrets.IssuerNamesFor(secretName)
}

func (s *state) RememberIssuerEABSecret(issuer resources.ObjectName, secretRef *v1.SecretReference, hash string) {
	s.eabSecrets.RememberIssuerSecret(issuer, secretRef, hash)
}

func (s *state) AddRenewalOverdue(certName resources.ObjectName) bool {
	return s.overdueCerts.Add(certName)
}

func (s *state) RemoveRenewalOverdue(certName resources.ObjectName) bool {
	return s.overdueCerts.Remove(certName)
}

func (s *state) GetAllRenewalOverdue() []resources.ObjectName {
	return s.overdueCerts.AsArray()
}

func (s *state) GetRenewalOverdueCount() int {
	return s.overdueCerts.Size()
}

func (s *state) AddRevoked(certName resources.ObjectName) bool {
	return s.revokedCerts.Add(certName)
}

func (s *state) RemoveRevoked(certName resources.ObjectName) bool {
	return s.revokedCerts.Remove(certName)
}

func (s *state) GetAllRevoked() []resources.ObjectName {
	return s.revokedCerts.AsArray()
}

func (s *state) GetRevokedCount() int {
	return s.revokedCerts.Size()
}
