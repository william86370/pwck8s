package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	rancher "pwck8s/rancher"
)

func HandelEnvir(Config GlobalConfig, w http.ResponseWriter, r *http.Request) {

	// Run CORS handler for the request first
	HandelCors(w, r)

	// Check if the request is for CORS preflight
	if r.Method == "OPTIONS" {
		// Just return header with no body, as preflight is just to check the CORS setting of the server
		w.WriteHeader(http.StatusOK)
		return
	}

	//Perform User Auth
	UserDN, err := GetUserDn(r)
	if err != nil {
		http.Error(w, Logboi(r, err.Error()), http.StatusUnauthorized)
		return
	}

	if r.Method == "GET" {
		handleGetEnvir(Config, w, r, UserDN)
	} else if r.Method == "POST" {
		handlePostEnvir(Config, w, r, UserDN)
	} else if r.Method == "DELETE" {
		// HandleDeleteEnvir(Config, w, r, UserDN)
	} else {
		http.Error(w, Logboi(r, "Invalid request method"), http.StatusMethodNotAllowed)
	}

}

// handleGetEnvir is going to create both the user object and the project object
func handleGetEnvir(Config GlobalConfig, w http.ResponseWriter, r *http.Request, UserDN string) {

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

func HandleCleanupEnvir(Config GlobalConfig, UserDN string) error {
	//TODO Implement

	return nil
}

func handlePostEnvir(Config GlobalConfig, w http.ResponseWriter, r *http.Request, UserDN string) {
	client := Config.Client

	// Check if the user already has a project
	exists, err := rancher.UserExists(client, UserDN)
	if err != nil {
		http.Error(w, Logboi(r, err.Error()), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, Logboi(r, "User already exists"), http.StatusBadRequest)
		return
	}

	// Check if the user already has a project
	err = rancher.EnsureNoDuplicateProject(client, UserDN, "local") //TODO Remove harrdcorded cluster IDs
	if err != nil {
		http.Error(w, Logboi(r, "Project already exists"), http.StatusBadRequest)
		return
	}

	// Generate a new user object
	user := rancher.GenerateUser(UserDN, Config.AuthProvider)

	// Create the user in Rancher
	err = rancher.CreateRancherUser(client, user)
	if err != nil {
		HandleCleanupEnvir(Config, UserDN)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the GlobalRoleBinding in Rancher
	err = rancher.CreateGlobalRoleBinding(client, user, Config.DefaultGlobalRole)
	if err != nil {
		HandleCleanupEnvir(Config, UserDN)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log
	Logboi(r, fmt.Sprintf("User Created: [%v]", UserDN))
	// Logboi(r, fmt.Sprintf("Project Created: [%v/%v]", project.ClusterID, project.ProjectID))

	// Return the user object
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding user: %v", err)
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}
