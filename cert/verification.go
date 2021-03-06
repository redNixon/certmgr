// Package cert contains certificate specifications and
// certificate-specific management.
package cert

import (
	"crypto/x509"
	"net"
	"sort"
)

// CertificateMatchesHostname checks if the Certificates hosts are the same as the given hosts
func CertificateMatchesHostname(hosts []string, cert *x509.Certificate) bool {
	a := make([]string, len(hosts))
	for idx := range hosts {
		// normalize the IPs.
		ip := net.ParseIP(hosts[idx])
		if ip == nil {
			a[idx] = hosts[idx]
		} else {
			a[idx] = ip.String()
		}
	}
	b := make([]string, len(cert.DNSNames), len(cert.DNSNames)+len(cert.IPAddresses))
	copy(b, cert.DNSNames)
	for idx := range cert.IPAddresses {
		b = append(b, cert.IPAddresses[idx].String())
	}

	if len(a) != len(b) {
		return false
	}

	sort.Strings(a)
	sort.Strings(b)
	for idx := range a {
		if a[idx] != b[idx] {
			return false
		}
	}
	return true
}

// CertificateChainVerify validates if a given cert is derived from the given CA
func CertificateChainVerify(ca *x509.Certificate, cert *x509.Certificate) error {
	roots := x509.NewCertPool()
	roots.AddCert(ca)
	_, err := cert.Verify(x509.VerifyOptions{
		Roots: roots,
	})
	return err
}
