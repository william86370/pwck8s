package api

import (
	"fmt"
	"net/http"
	x509toolkit "pwck8s/x509"
)

func GetUserDn(r *http.Request) (string, error) {

	//Get the cert from the request
	x509, err := x509toolkit.ParseCertificate(r, "X-Client-Certificate")
	if err != nil {
		// Log the cert error but dont error out
		fmt.Printf("Error parsing certificate: %v\n", err)
	} else {
		fmt.Printf("[X509] %v\n", x509toolkit.ParseDN(x509))
	}

	// Get the UserDN from the request
	UserDN := r.Header.Get("UserDN")
	if UserDN == "" {
		return "", fmt.Errorf("UserDN not found")
	}
	return UserDN, nil
}
