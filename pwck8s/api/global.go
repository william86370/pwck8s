package api

import (
	"net/http"

	"k8s.io/client-go/dynamic"
)

type GlobalConfig struct {
	Client             dynamic.Interface
	ClusterID          string
	AuthProvider       string
	DefaultProjectRole string
	DefaultGlobalRole  string
	Debug              bool
}

// HandelCors sets the CORS headers for the response
func HandelCors(w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, UserDN")

}
