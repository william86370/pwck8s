package x509

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
)

// ParseCertificate extracts a client certificate from HTTP headers and returns the parsed certificate
func ParseCertificate(r *http.Request, HttpCertHeader string) (*x509.Certificate, error) {
	// Check if the header is set
	if HttpCertHeader == "" {
		HttpCertHeader = "X-Client-Certificate"
	}
	// Get the certificate from the header
	certHeader := r.Header.Get(HttpCertHeader)
	if certHeader == "" {
		return nil, fmt.Errorf("no client certificate found in header")
	}

	// Decode the base64-encoded certificate
	decodedCert, err := base64.StdEncoding.DecodeString(certHeader)
	if err != nil {
		return nil, fmt.Errorf("error decoding certificate: %v", err)
	}

	// Parse the PEM encoded certificate
	block, _ := pem.Decode(decodedCert)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	return cert, nil
}

// ValidateCACertificate validates a client certificate (*x509.Certificate) against the given CA certificate in PEM format
func ValidateCACertificate(clientCert *x509.Certificate, caCertPEM []byte) error {
	// Decode the CA certificate
	caBlock, _ := pem.Decode(caCertPEM)
	if caBlock == nil {
		return fmt.Errorf("failed to parse CA certificate PEM")
	}
	caCert, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse CA certificate: %v", err)
	}

	// Create a new CertPool and add the CA certificate to it
	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(caCert)

	// Verify the client certificate against the CA certificate
	opts := x509.VerifyOptions{
		Roots: caCertPool,
	}

	if _, err := clientCert.Verify(opts); err != nil {
		return fmt.Errorf("failed to verify client certificate: %v", err)
	}

	return nil
}

// ParseDN returns the Distinguished Name (DN) of the given x509.Certificate
func ParseDN(cert *x509.Certificate) string {
	if cert == nil {
		return ""
	}

	var dnParts []string

	if len(cert.Subject.CommonName) > 0 {
		dnParts = append(dnParts, fmt.Sprintf("CN=%s", cert.Subject.CommonName))
	}
	for _, org := range cert.Subject.Organization {
		dnParts = append(dnParts, fmt.Sprintf("O=%s", org))
	}
	for _, country := range cert.Subject.Country {
		dnParts = append(dnParts, fmt.Sprintf("C=%s", country))
	}
	// Add other components as needed, such as OrganizationalUnit, Locality, etc.

	return fmt.Sprintf("/%s", stringJoin(dnParts, "/"))
}

// stringJoin is a helper function to join strings with a separator
func stringJoin(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for _, part := range parts[1:] {
		result += sep + part
	}
	return result
}






