package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorReset  = "\033[0m"
)

type GlobalConfig struct {
	// The port to listen on
	Port int `json:"port"`
	// The path to the TLS certificate
	TLSCert string `json:"tls_cert"`
	// The path to the TLS key
	TLSKey string `json:"tls_key"`
	// CA certificate pool to verify client certificates
	CACert *x509.CertPool `json:"ca_cert"`
	// Service to proxy to
	ProxyURL string `json:"proxy_url"`
	// Debug mode
	Debug bool `json:"debug"`
}

// LoadCACertPool loads a CA certificate from a given file and returns an x509.CertPool.
func LoadCACertPool(caCertPath string) (*x509.CertPool, error) {
	// Read the CA certificate file
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
		return nil, err
	}

	// Create a new CertPool
	caCertPool := x509.NewCertPool()

	// Append the CA certificate to the pool
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("Failed to append CA certificate")
		return nil, err
	}

	return caCertPool, nil
}

func GetConfigFromEnvProduction() GlobalConfig {
	// Get Config From ENV but fail if not found
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("PORT must be set")
	}

	tlsCert := os.Getenv("TLS_CERT")
	if tlsCert == "" {
		log.Fatal("TLS_CERT must be set")
	}

	tlsKey := os.Getenv("TLS_KEY")
	if tlsKey == "" {
		log.Fatal("TLS_KEY must be set")
	}

	caCert := os.Getenv("CA_CERT")
	if caCert == "" {
		log.Fatal("CA_CERT must be set")
	}

	proxyURL := os.Getenv("PROXY_URL")
	if proxyURL == "" {
		log.Fatal("PROXY_URL must be set")
	}

	debug := os.Getenv("DEBUG")
	if debug == "" {
		debug = "false"
	}
	debugMode, err := strconv.ParseBool(debug)
	if err != nil {
		return GlobalConfig{}
	}

	caCertPool, err := LoadCACertPool(caCert)
	if err != nil {
		log.Fatal("Error loading CA certificate pool:", err)
	}

	return GlobalConfig{
		Port:     port,
		TLSCert:  tlsCert,
		TLSKey:   tlsKey,
		CACert:   caCertPool,
		ProxyURL: proxyURL,
		Debug:    debugMode,
	}
}

func HandleConfig() (GlobalConfig, error) {
	//Check if Debug mode is enabled
	debug := os.Getenv("DEBUG")
	if debug == "" {
		debug = "false"
	}
	debugMode, err := strconv.ParseBool(debug)
	if err != nil {
		return GlobalConfig{}, err
	}

	// If debug mode is enabled, use the development config
	if debugMode {
		// Print debug info
		log.Println("Running in debug mode")
		return GetConfigFromEnvProduction(), nil
	}

	// Otherwise, use the production config
	log.Println("Running in production mode")
	return GetConfigFromEnvProduction(), nil

}

type HttpHeaderMap struct {
	CN string `json:"cn"`
	DN string `json:"dn"`
}

func GetHttpHeaderMapFromEnv() HttpHeaderMap {

	cn := os.Getenv("HTTP_HEADER_CN")
	if cn == "" {
		cn = "X-Client-Cn"
	}

	dn := os.Getenv("HTTP_HEADER_DN")
	if dn == "" {
		cn = "X-Client-Dn"
	}

	return HttpHeaderMap{
		CN: cn,
		DN: dn,
	}
}

func HandleProxy(w http.ResponseWriter, r *http.Request, config GlobalConfig, proxy *httputil.ReverseProxy, httpHeaderMap HttpHeaderMap) {

	// Print the incoming request's headers
	// Log the details of the request
	log.Printf("%sRemote Address: %s%s - %sProtocol: %s%s - %sMethod: %s%s - %sURL: %s%s",
		colorCyan, r.RemoteAddr, colorReset,
		colorGreen, r.Proto, colorReset,
		colorYellow, r.Method, colorReset,
		colorPurple, r.URL, colorReset)

	// Check if the request has TLS and a client certificate
	if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {

		cert := r.TLS.PeerCertificates[0]
		commonName := cert.Subject.CommonName

		// Construct the DN from the certificate DN = CN,OU's,O
		dn := cert.Subject.CommonName
		for _, ou := range cert.Subject.OrganizationalUnit {
			dn = dn + "," + ou
		}
		for _, o := range cert.Subject.Organization {
			dn = dn + "," + o
		}

		// Set the headers
		r.Header.Set(httpHeaderMap.CN, commonName)
		r.Header.Set(httpHeaderMap.DN, string(dn))
	}
	// Forward the request to the proxy
	proxy.ServeHTTP(w, r)
}

func main() {

	// Get the config
	config, err := HandleConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Get the HttpHeaderMap
	httpHeaderMap := GetHttpHeaderMapFromEnv()

	// Create the proxy as before
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// req.URL.Scheme = "http"
			req.URL.Host = config.ProxyURL
		},
	}

	// Update the server to use HandleProxy
	server := &http.Server{
		Addr: ":" + strconv.Itoa(config.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			HandleProxy(w, r, config, proxy, httpHeaderMap)
		}),
		TLSConfig: &tls.Config{
			ClientCAs:  config.CACert,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	// Start the server
	log.Printf("%sStarting server on port %d%s\n", colorGreen, config.Port, colorReset)
	log.Fatal(fmt.Sprintf("%sServer stopped with error: %s%s", colorRed, server.ListenAndServeTLS(config.TLSCert, config.TLSKey), colorReset))

}
