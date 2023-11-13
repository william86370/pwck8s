package main

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"testing"
	"x509-proxy/x509toolkit"
)

func TestHandleProxy(t *testing.T) {
	// Create a test ReverseProxy
	proxy := &httputil.ReverseProxy{}

	// Create a test HttpHeaderMap
	httpHeaderMap := HttpHeaderMap{
		CN: "CN",
		DN: "DN",
	}
	// Create a test GlobalConfig
	config := GlobalConfig{}

	// Create a test http request
	req := httptest.NewRequest("GET", "http://example.com", nil)

	// Create a test certificate
	cert := &x509.Certificate{
		Subject: pkix.Name{
			CommonName:         "test",
			OrganizationalUnit: []string{"testOU"},
			Organization:       []string{"testO"},
		},
	}

	// Add the certificate to the request
	req.TLS = &tls.ConnectionState{
		PeerCertificates: []*x509.Certificate{cert},
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call HandleProxy
	HandleProxy(rr, req, config, proxy, httpHeaderMap)

	// Check the headers
	if req.Header.Get(httpHeaderMap.CN) != cert.Subject.CommonName {
		t.Errorf("Expected CN header to be %s, got %s", cert.Subject.CommonName, req.Header.Get(httpHeaderMap.CN))
	}

	expectedDN := "test,testOU,testO"
	if req.Header.Get(httpHeaderMap.DN) != expectedDN {
		t.Errorf("Expected DN header to be %s, got %s", expectedDN, req.Header.Get(httpHeaderMap.DN))
	}
}
func TestGetConfigFromEnvProduction(t *testing.T) {
	// Set up test environment variables
	os.Setenv("PORT", "8080")
	os.Setenv("TLS_CERT", "test-cert")
	os.Setenv("TLS_KEY", "test-key")
	os.Setenv("CA_CERT", "test-ca-cert")
	os.Setenv("PROXY_URL", "http://example.com")
	os.Setenv("DEBUG", "true")

	// Write a test CA certificate to disk for the test environment variables to use as a CA certificate file path
	// generate a test CA certificate
	caCert, _, err := x509toolkit.GenerateCACertificate("test", "US", "test", "testOU")
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCert.Raw,
	})

	err = os.WriteFile("test-ca-cert", certPEM, 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Call the function
	config := GetConfigFromEnvProduction()

	// Check the returned values
	if config.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", config.Port)
	}
	if config.TLSCert != "test-cert" {
		t.Errorf("Expected TLSCert to be 'test-cert', got '%s'", config.TLSCert)
	}
	if config.TLSKey != "test-key" {
		t.Errorf("Expected TLSKey to be 'test-key', got '%s'", config.TLSKey)
	}
	if config.ProxyURL != "http://example.com" {
		t.Errorf("Expected ProxyURL to be 'http://example.com', got '%s'", config.ProxyURL)
	}
	if config.Debug != true {
		t.Errorf("Expected Debug to be true, got %t", config.Debug)
	}

	// Clean up test environment variables
	os.Unsetenv("PORT")
	os.Unsetenv("TLS_CERT")
	os.Unsetenv("TLS_KEY")
	os.Unsetenv("CA_CERT")
	os.Unsetenv("PROXY_URL")
	os.Unsetenv("DEBUG")
	// Clean up test CA certificate
	os.Remove("test-ca-cert")
}
