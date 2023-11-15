package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	rancher "pwck8s/rancher"
)

func Logboi(r *http.Request, s string) string {
	log.Printf("[%s] [%s] %s", r.Method, r.URL.String(), s)
	return fmt.Sprintf("[%s] [%s] %s", r.Method, r.URL.String(), s)
}

// GenerateProject Creates a new rancher project object with default values for UserDN
// Takes UserDN (string)
func GenerateProject(UserDN string, ClusterID string) rancher.Project {
	return rancher.Project{
		ProjectID:      rancher.GenerateProjectId(),
		ClusterID:      ClusterID,
		OwnerDN:        UserDN,
		CreationTime:   time.Now(),
		ExpirationTime: time.Now().Add(time.Hour), // 1 hour from now
		DisplayName:    UserDN,
		Resources:      rancher.DefaultResources(), // Default resources
	}
}

// /api/v1/project
func ProjectHandler(Config GlobalConfig, w http.ResponseWriter, r *http.Request) {

	Logboi(r, fmt.Sprintf("Handling request from [%s]", r.RemoteAddr))

	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, UserDN")

	// Check if the request is for CORS preflight
	if r.Method == "OPTIONS" {
		// Just return header with no body, as preflight is just to check the CORS setting of the server
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == "GET" {
		handleGetProject(Config, w, r)
	} else if r.Method == "POST" {
		handlePostProject(Config, w, r)
	} else if r.Method == "DELETE" {
		HandleDeleteProject(Config, w, r)
	} else {
		http.Error(w, Logboi(r, "Invalid request method"), http.StatusMethodNotAllowed)
	}
}

func HandleDeleteProject(Config GlobalConfig, w http.ResponseWriter, r *http.Request) {
	// Get the UserDN from the request
	UserDN := r.Header.Get("UserDN")
	if UserDN == "" {
		http.Error(w, Logboi(r, "UserDN not found"), http.StatusBadRequest)
		return
	}
	client := Config.Client

	// Get the project from the UserDN
	project, err := rancher.GetProjectByOwner(client, UserDN, "local") //TODO Remove harrdcorded cluster IDs
	if err != nil {
		http.Error(w, Logboi(r, fmt.Sprintf("Error getting project: %v", err)), http.StatusInternalServerError)
		return
	}

	// Delete the project
	err = rancher.DeleteRancherProject(client, project.ProjectID, project.ClusterID)
	if err != nil {
		http.Error(w, Logboi(r, fmt.Sprintf("Error deleting project: %v", err)), http.StatusInternalServerError)
		return
	}
	Logboi(r, fmt.Sprintf("Project Deleted: [%v/%v]", project.ClusterID, project.ProjectID))
	//Set http code to deleted
	w.WriteHeader(http.StatusOK)
}

func handlePostProject(Config GlobalConfig, w http.ResponseWriter, r *http.Request) {

	// Get the UserDN from the request
	UserDN := r.Header.Get("UserDN")
	if UserDN == "" {
		http.Error(w, Logboi(r, "UserDN not found"), http.StatusBadRequest)
		return
	}
	client := Config.Client

	//Get the User object from Rancher using the UserDN
	user, err := rancher.GetRancherUser(client, UserDN)
	if err != nil {
		http.Error(w, Logboi(r, fmt.Sprintf("Error: %v", err)), http.StatusInternalServerError)
		return
	}

	// Check if the user already has a project
	err = rancher.EnsureNoDuplicateProject(client, UserDN, "local") //TODO Remove harrdcorded cluster IDs
	if err != nil {
		http.Error(w, Logboi(r, fmt.Sprintf("Error: %v", err)), http.StatusInternalServerError)
		return
	}

	// Generate a new project object
	project := GenerateProject(UserDN, Config.ClusterID)

	// Create the project in Rancher
	err = rancher.CreateRancherProject(client, project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the ProjectRoleTemplateBinding in Rancher
	err = rancher.CreateProjectRoleBinding(client, user.UserID, project, Config.DefaultProjectRole)
	if err != nil {
		log.Printf("Error creating ProjectRoleBinding: %v", err)
		http.Error(w, "Error creating ProjectRoleBinding", http.StatusInternalServerError)
		return
	}

	// Return the project object
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(project); err != nil {
		log.Printf("[handlePostProject] Error encoding project: %v", err)
		http.Error(w, "[handlePostProject] Error encoding response", http.StatusInternalServerError)
		return
	}

}

// handleGetProjects handles GET requests to /api/v1/project
func handleGetProject(Config GlobalConfig, w http.ResponseWriter, r *http.Request) {

	// Get the UserDN from the request
	UserDN := r.Header.Get("UserDN")
	if UserDN == "" {
		http.Error(w, Logboi(r, "UserDN not found"), http.StatusBadRequest)
		return
	}
	client := Config.Client

	// Get the project from the UserDN
	project, err := rancher.GetProjectByOwner(client, UserDN, "local") //TODO Remove harrdcorded cluster IDs
	if err != nil {
		http.Error(w, Logboi(r, fmt.Sprintf("Error getting project: %v", err)), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}
