package x509toolkit

import (
	"testing"
)

func TestGenerateCACertificate(t *testing.T) {
	commonName := "Test CA"
	country := "US"
	organization := "Test Org"
	organizationUnit := "Test Org Unit"

	cert, _, err := GenerateCACertificate(commonName, country, organization, organizationUnit)
	if err != nil {
		t.Fatalf("Failed to generate CA certificate: %v", err)
	}

	// Verify that the certificate is valid
	err = cert.CheckSignatureFrom(cert)
	if err != nil {
		t.Fatalf("Failed to verify CA certificate: %v", err)
	}

	// Verify that the certificate is a CA certificate
	if !cert.IsCA {
		t.Fatal("Expected CA certificate, got non-CA certificate")
	}

	// Verify that the certificate has the correct subject details
	if cert.Subject.CommonName != commonName {
		t.Errorf("Expected CommonName to be %s, got %s", commonName, cert.Subject.CommonName)
	}
	if cert.Subject.Country[0] != country {
		t.Errorf("Expected Country to be %s, got %s", country, cert.Subject.Country[0])
	}
	if cert.Subject.Organization[0] != organization {
		t.Errorf("Expected Organization to be %s, got %s", organization, cert.Subject.Organization[0])
	}
	if cert.Subject.OrganizationalUnit[0] != organizationUnit {
		t.Errorf("Expected OrganizationalUnit to be %s, got %s", organizationUnit, cert.Subject.OrganizationalUnit[0])
	}

}
