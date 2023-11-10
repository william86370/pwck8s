package rancher

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// // Defines the structure for Rancher authentication
// type RancherAuth struct {
// 	URL           string `json:"url"`           //"https://your-rancher-server/v3/"
// 	ApiKey        string `json:"apiKey"`        //"your-rancher-api-key"
// 	DemoClusterID string `json:"demoClusterID"` //"your-rancher-demo-cluster-id"
// }

// DefaultResources creates a Resources struct with default Kubernetes-aligned values
func DefaultResources() Resources {
	return Resources{
		Pods:                   "15",
		Services:               "50",
		ReplicationControllers: "50",
		Secrets:                "50",
		ConfigMaps:             "50",
		PersistentVolumeClaims: "50",
		ServicesNodePorts:      "0",
		ServicesLoadBalancers:  "0",
		RequestsStorage:        "10Gi",
		LimitsCPU:              "2",
		LimitsMemory:           "4Gi",
	}
}
func MapProjects(unstructuredProjects []unstructured.Unstructured) ([]Project, error) {
	var projects []Project
	for _, unstructuredProject := range unstructuredProjects {
		project, err := MapToProject(unstructuredProject)
		if err != nil {
			return nil, err // Or handle the error based on your use case
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func MapToProject(tmpProject unstructured.Unstructured) (Project, error) {
	var project Project

	// Extract data using unstructured getters
	projectID, found, err := unstructured.NestedString(tmpProject.Object, "metadata", "name")
	if err != nil || !found {
		return project, fmt.Errorf("ProjectID not found or error in reading: %v", err)
	}
	project.ProjectID = projectID

	// Extract data using unstructured getters
	ClusterID, found, err := unstructured.NestedString(tmpProject.Object, "spec", "clusterName")
	if err != nil || !found {
		return project, fmt.Errorf("ClusterID not found or error in reading: %v", err)
	}
	project.ClusterID = ClusterID

	// Extract data using unstructured getters
	DisplayName, found, err := unstructured.NestedString(tmpProject.Object, "spec", "displayName")
	if err != nil || !found {
		return project, fmt.Errorf("displayname not found or error in reading: %v", err)
	}

	project.DisplayName = DisplayName

	// Extract data using unstructured getters
	OwnerDN, found, err := unstructured.NestedString(tmpProject.Object, "metadata", "labels", "pwck8s/ownerdn")
	if err != nil || !found {
		return project, fmt.Errorf("ownerdn not found or error in reading: %v", err)
	}
	project.OwnerDN = OwnerDN

	// Note: Handle timestamps and other complex types appropriately
	CreationTime, found, err := unstructured.NestedString(tmpProject.Object, "metadata", "labels", "pwck8s/creationtime")
	if err != nil || !found {
		return project, fmt.Errorf("creationtime not found or error in reading: %v", err)
	}
	project.CreationTime, err = time.Parse("2006-01-02T15-04-05Z07-00", CreationTime) //TODO please fix "2006-01-02T15-04-05Z07-00" to a constant
	if err != nil {
		return project, fmt.Errorf("error parsing creationtime: %v", err)
	}

	// Note: Handle timestamps and other complex types appropriately
	ExpirationTime, found, err := unstructured.NestedString(tmpProject.Object, "metadata", "labels", "pwck8s/expirationtime")
	if err != nil || !found {
		return project, fmt.Errorf("expirationtime not found or error in reading: %v", err)
	}
	project.ExpirationTime, err = time.Parse("2006-01-02T15-04-05Z07-00", ExpirationTime)
	if err != nil {
		return project, fmt.Errorf("error parsing expirationtime: %v", err)
	}

	// Extract Resources using unstructured getters from spec.resourceQuota.limit
	Pods, found, err := unstructured.NestedString(tmpProject.Object, "spec", "resourceQuota", "limit", "pods")
	if err != nil || !found {
		return project, fmt.Errorf("pods not found or error in reading: %v", err)
	}
	project.Resources.Pods = Pods

	Services, found, err := unstructured.NestedString(tmpProject.Object, "spec", "resourceQuota", "limit", "services")
	if err != nil || !found {
		return project, fmt.Errorf("services not found or error in reading: %v", err)
	}
	project.Resources.Services = Services

	ReplicationControllers, found, err := unstructured.NestedString(tmpProject.Object, "spec", "resourceQuota", "limit", "replicationControllers")
	if err != nil || !found {
		return project, fmt.Errorf("replicationControllers not found or error in reading: %v", err)
	}
	project.Resources.ReplicationControllers = ReplicationControllers

	Secrets, found, err := unstructured.NestedString(tmpProject.Object, "spec", "resourceQuota", "limit", "secrets")
	if err != nil || !found {
		return project, fmt.Errorf("secrets not found or error in reading: %v", err)
	}
	project.Resources.Secrets = Secrets

	ConfigMaps, found, err := unstructured.NestedString(tmpProject.Object, "spec", "resourceQuota", "limit", "configMaps")
	if err != nil || !found {
		return project, fmt.Errorf("configMaps not found or error in reading: %v", err)
	}
	project.Resources.ConfigMaps = ConfigMaps

	PersistentVolumeClaims, found, err := unstructured.NestedString(tmpProject.Object, "spec", "resourceQuota", "limit", "persistentVolumeClaims")
	if err != nil || !found {
		return project, fmt.Errorf("persistentVolumeClaims not found or error in reading: %v", err)
	}
	project.Resources.PersistentVolumeClaims = PersistentVolumeClaims

	ServicesNodePorts, found, err := unstructured.NestedString(tmpProject.Object, "spec", "resourceQuota", "limit", "servicesNodePorts")
	if err != nil || !found {
		return project, fmt.Errorf("servicesNodePorts not found or error in reading: %v", err)
	}
	project.Resources.ServicesNodePorts = ServicesNodePorts

	ServicesLoadBalancers, found, err := unstructured.NestedString(tmpProject.Object, "spec", "resourceQuota", "limit", "servicesLoadBalancers")
	if err != nil || !found {
		return project, fmt.Errorf("servicesLoadBalancers not found or error in reading: %v", err)
	}
	project.Resources.ServicesLoadBalancers = ServicesLoadBalancers

	RequestsStorage, found, err := unstructured.NestedString(tmpProject.Object, "spec", "resourceQuota", "limit", "requestsStorage")
	if err != nil || !found {
		return project, fmt.Errorf("requestsStorage not found or error in reading: %v", err)
	}
	project.Resources.RequestsStorage = RequestsStorage

	LimitsCPU, found, err := unstructured.NestedString(tmpProject.Object, "spec", "resourceQuota", "limit", "limitsCpu")
	if err != nil || !found {
		return project, fmt.Errorf("limitsCpu not found or error in reading: %v", err)
	}
	project.Resources.LimitsCPU = LimitsCPU

	LimitsMemory, found, err := unstructured.NestedString(tmpProject.Object, "spec", "resourceQuota", "limit", "limitsMemory")
	if err != nil || !found {
		return project, fmt.Errorf("limitsMemory not found or error in reading: %v", err)
	}
	project.Resources.LimitsMemory = LimitsMemory

	return project, nil
}

type Resources struct {
	// ProjectSize Object
	Pods                   string `json:"pods,omitempty"`
	Services               string `json:"services,omitempty"`
	ReplicationControllers string `json:"replicationControllers,omitempty"`
	Secrets                string `json:"secrets,omitempty"`
	ConfigMaps             string `json:"configMaps,omitempty"`
	PersistentVolumeClaims string `json:"persistentVolumeClaims,omitempty"`
	ServicesNodePorts      string `json:"servicesNodePorts,omitempty"`
	ServicesLoadBalancers  string `json:"servicesLoadBalancers,omitempty"`
	RequestsStorage        string `json:"requestsStorage,omitempty"`
	LimitsCPU              string `json:"limitsCpu,omitempty"`
	LimitsMemory           string `json:"limitsMemory,omitempty"`
}

// Project defines the response structure for the project endpoint
type Project struct {
	// Project Object
	ProjectID      string    `json:"projectId"`
	ClusterID      string    `json:"clusterId"`
	DisplayName    string    `json:"displayName"`
	Resources      Resources `json:"resources"`
	CreationTime   time.Time `json:"creationTime"`
	ExpirationTime time.Time `json:"expirationTime"`
	OwnerDN        string    `json:"ownerDn"`
}

func GenerateProjectId() string {
	// Generate a new project ID similar to Rancher project ID
	// p-<random 5 char string>
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return "pwck8s-" + string(b)
}

func EnsureNoDuplicateProject(client dynamic.Interface, OwnerDN string, ClusterID string) error {
	// Function to check if a project with the same OwnerDN already exists
	// If it exists, return an error
	projects, err := GetProjectsByOwner(client, OwnerDN, ClusterID)
	if err != nil {
		return fmt.Errorf("failed to get projects: %v", err)
	}

	if len(projects) > 0 {
		return fmt.Errorf("project already exists")
	}
	return nil
}

func DeleteRancherProject(client dynamic.Interface, ProjectID string, ClusterID string) error {
	projectGVR := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "projects",
	}
	// Delete the project
	err := client.Resource(projectGVR).Namespace(ClusterID).Delete(context.TODO(), ProjectID, v1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete project: %v", err)
	}
	return nil
}

func CreateRancherProject(client dynamic.Interface, newProject Project) error {
	projectGVR := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "projects",
	}

	// Define the ProjectSpec according to Rancher's API specification
	projectSpec := map[string]interface{}{
		"displayName": newProject.DisplayName,
		// Description is a human-readable description of the project.
		"description": "PWCK8S Project",
		// ClusterID is the ID of the cluster where the project will be created.
		"clusterName": newProject.ClusterID,
		// ResourceQuota is a specification for the total amount of quota for standard resources that will be shared by all namespaces in the project.
		"resourceQuota": map[string]interface{}{
			"limit": map[string]interface{}{
				"pods":                   newProject.Resources.Pods,
				"services":               newProject.Resources.Services,
				"replicationControllers": newProject.Resources.ReplicationControllers,
				"secrets":                newProject.Resources.Secrets,
				"configMaps":             newProject.Resources.ConfigMaps,
				"persistentVolumeClaims": newProject.Resources.PersistentVolumeClaims,
				"servicesNodePorts":      newProject.Resources.ServicesNodePorts,
				"servicesLoadBalancers":  newProject.Resources.ServicesLoadBalancers,
				"requestsStorage":        newProject.Resources.RequestsStorage,
				"limitsCpu":              newProject.Resources.LimitsCPU,
				"limitsMemory":           newProject.Resources.LimitsMemory,
			},
		},
		"namespaceDefaultResourceQuota": map[string]interface{}{
			"limit": map[string]interface{}{
				"pods":                   newProject.Resources.Pods,
				"services":               newProject.Resources.Services,
				"replicationControllers": newProject.Resources.ReplicationControllers,
				"secrets":                newProject.Resources.Secrets,
				"configMaps":             newProject.Resources.ConfigMaps,
				"persistentVolumeClaims": newProject.Resources.PersistentVolumeClaims,
				"servicesNodePorts":      newProject.Resources.ServicesNodePorts,
				"servicesLoadBalancers":  newProject.Resources.ServicesLoadBalancers,
				"requestsStorage":        newProject.Resources.RequestsStorage,
				"limitsCpu":              newProject.Resources.LimitsCPU,
				"limitsMemory":           newProject.Resources.LimitsMemory,
			},
		},
	}

	project := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "management.cattle.io/v3",
			"kind":       "Project",
			"metadata": map[string]interface{}{
				"name": newProject.ProjectID,
				"labels": map[string]string{
					"pwck8s/ownerdn":        newProject.OwnerDN, // Add the label with the user's DN
					"pwck8s/displayname":    newProject.DisplayName,
					"pwck8s/projectid":      newProject.ProjectID,
					"pwck8s/clusterid":      newProject.ClusterID,
					"pwck8s/creationtime":   newProject.CreationTime.Format("2006-01-02T15-04-05Z07-00"),
					"pwck8s/expirationtime": newProject.ExpirationTime.Format("2006-01-02T15-04-05Z07-00"),
				},
			},
			"spec": projectSpec,
		},
	}

	// Check for the dryrun flag //TODO - move this to the API
	dryrun := os.Getenv("DRYRUN")
	if dryrun == "true" {
		fmt.Printf("[CreateRancherProject] Project Object:\n%v\n", project)
		return nil
	}

	// Create the project in Rancher
	_, err := client.Resource(projectGVR).Namespace(newProject.ClusterID).Create(context.TODO(), project, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("[CreateRancherProject] failed to create project: %v", err)
	}

	fmt.Printf("[CreateRancherProject] Project created: %s\n", newProject.ProjectID)
	return nil
}

func GetProjectsByOwner(client dynamic.Interface, OwnerDN string, ClusterID string) ([]unstructured.Unstructured, error) {
	projectGVR := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "projects",
	}

	labelSelector := labels.Set(map[string]string{"pwck8s/ownerdn": OwnerDN}).AsSelector().String()
	listOptions := v1.ListOptions{LabelSelector: labelSelector}
	projectList, err := client.Resource(projectGVR).Namespace(ClusterID).List(context.TODO(), listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %v", err)
	}

	return projectList.Items, nil
}

func GetProjectByOwner(client dynamic.Interface, OwnerDN string, ClusterID string) (Project, error) {
	var project Project
	projectGVR := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "projects",
	}

	labelSelector := labels.Set(map[string]string{"pwck8s/ownerdn": OwnerDN}).AsSelector().String()
	listOptions := v1.ListOptions{LabelSelector: labelSelector}
	projectList, err := client.Resource(projectGVR).Namespace(ClusterID).List(context.TODO(), listOptions)
	if err != nil {
		return project, fmt.Errorf("failed to list project: %v", err)
	}
	// If there are multiple projects with the same OwnerDN, return an error
	if len(projectList.Items) > 1 {
		return project, fmt.Errorf("multiple projects found")
	}
	// If there are no projects with the same OwnerDN, return an error
	if len(projectList.Items) == 0 {
		return project, fmt.Errorf("no projects found")
	}

	// Extract the project from the list
	project, err = MapToProject(projectList.Items[0])
	if err != nil {
		return project, fmt.Errorf("failed to map project: %v", err)
	}

	return project, nil
}
