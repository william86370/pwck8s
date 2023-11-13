package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func headersHandler(w http.ResponseWriter, r *http.Request) {

	// Print the incoming request's headers
	log.Println("Handling request for", r.URL.Path)

	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Convert headers to JSON and send response
	json.NewEncoder(w).Encode(r.Header)
}

func main() {
	// Handle all requests with the headersHandler function
	http.HandleFunc("/", headersHandler)

	// Print a log message indicating that the server has started
	log.Println("Server started on port 8080")
	// Listen and serve on port 8080
	http.ListenAndServe(":8080", nil)
}
