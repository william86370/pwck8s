package api

import (
	"encoding/json"
	"log"
	"net/http"

	rancher "pwck8s/rancher"
)

// /api/v1/user
func UserHandler(Config GlobalConfig, w http.ResponseWriter, r *http.Request) {

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, UserDN")

	// Check if the request is for CORS preflight
	if r.Method == "OPTIONS" {
		// Just return header with no body, as preflight is just to check the CORS setting of the server
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get the UserDN from the request
	UserDN, err := GetUserDn(r)
	if err != nil {
		http.Error(w, Logboi(r, err.Error()), http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" {
		handleGetUser(Config, w, r, UserDN)
	} else if r.Method == "POST" {
		handlePostUser(Config, w, r, UserDN)
	} else if r.Method == "DELETE" {
		HandleDeleteUser(Config, w, r, UserDN)
	} else {
		http.Error(w, Logboi(r, "Invalid request method"), http.StatusMethodNotAllowed)
	}
}

func handleGetUser(Config GlobalConfig, w http.ResponseWriter, r *http.Request, UserDN string) {

	client := Config.Client
	// Get the user from the UserDN
	user, err := rancher.GetRancherUser(client, UserDN)
	if err != nil {
		http.Error(w, Logboi(r, err.Error()), http.StatusInternalServerError)
		return
	}
	// Log
	Logboi(r, "["+UserDN+"]")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding user: %v", err)
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

func HandleCleanupUser(Config GlobalConfig, UserDN string) error {
	//TODO Implement

	return nil
}

func handlePostUser(Config GlobalConfig, w http.ResponseWriter, r *http.Request, UserDN string) {
	client := Config.Client

	// Check if the user already has a project
	exists, err := rancher.UserExists(client, UserDN)
	if err != nil {
		http.Error(w, Logboi(r, err.Error()), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, Logboi(r, "User already exists"), http.StatusConflict)
		return
	}

	// Generate a new user object
	user := rancher.GenerateUser(UserDN, Config.AuthProvider)

	// Create the user in Rancher
	err = rancher.CreateRancherUser(client, user)
	if err != nil {
		HandleCleanupUser(Config, UserDN)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the GlobalRoleBinding in Rancher
	err = rancher.CreateGlobalRoleBinding(client, user, Config.DefaultGlobalRole)
	if err != nil {
		HandleCleanupUser(Config, UserDN)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the user object
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding user: %v", err)
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}

}

func HandleDeleteUser(Config GlobalConfig, w http.ResponseWriter, r *http.Request, UserDN string) {
	client := Config.Client

	// Check if the user already has a user
	exists, err := rancher.UserExists(client, UserDN)
	if err != nil {
		http.Error(w, Logboi(r, err.Error()), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, Logboi(r, "User does not exist"), http.StatusNotFound)
		return
	}

	// Delete the user in Rancher
	err = rancher.DeleteRancherUser(client, UserDN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the GlobalRoleBinding in Rancher
	err = rancher.DeleteGlobalRoleBinding(client, UserDN)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the user object
	w.WriteHeader(http.StatusOK)
}
