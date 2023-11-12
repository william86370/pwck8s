package api

import (
	"encoding/json"
	"log"
	"net/http"

	rancher "pwck8s/rancher"
)

// /api/v1/user
func UserHandler(Config GlobalConfig, w http.ResponseWriter, r *http.Request) {

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
	//Set http code to deleted
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding user: %v", err)
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
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
	user := rancher.GenerateUser(UserDN)

	// Create the user in Rancher
	err = rancher.CreateRancherUser(client, user)
	if err != nil {
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

	// Return the user object
	w.WriteHeader(http.StatusOK)
}
