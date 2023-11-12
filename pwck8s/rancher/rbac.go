package rancher

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func CreateGlobalRoleBinding(client dynamic.Interface, newUser User, globalRoleName string) error {
	// Define the User CRD we want to create
	grbGVR := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "globalrolebindings",
	}

	// Define the GlobalRoleBinding CRD we want to create
	grb := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "management.cattle.io/v3",
			"kind":       "GlobalRoleBinding",
			"metadata": map[string]interface{}{
				"name": newUser.UserID,
				"labels": map[string]string{
					"pwck8s/userid":         newUser.UserID,
					"pwck8s/userdn":         newUser.UserDN,
					"pwck8s/ownerdn":        newUser.UserDN,
					"pwck8s/creationtime":   newUser.CreationTime.Format("2006-01-02T15-04-05Z07-00"),
					"pwck8s/expirationtime": newUser.ExpirationTime.Format("2006-01-02T15-04-05Z07-00"),
				},
			},
			"globalRoleName": globalRoleName,
			"userName":       newUser.UserID,
		},
	}

	// Create the GRB in Rancher
	_, err := client.Resource(grbGVR).Namespace("").Create(context.TODO(), grb, v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create GlobalRoleBinding: %v", err)
	}
	fmt.Printf("GlobalRoleBinding created: %s\n", newUser.UserID)
	return nil
}

func DeleteGlobalRoleBinding(client dynamic.Interface, OwnerDN string) error {

	// Define the GlobalRoleBinding CRD we want to delete
	grbGVR := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "globalrolebindings",
	}
	// Get the GlobalRoleBinding in Rancher
	grb, err := GetGlobalRoleBinding(client, OwnerDN)
	if err != nil {
		return err
	}

	// Delete the GlobalRoleBinding in Rancher
	err2 := client.Resource(grbGVR).Namespace("").Delete(context.TODO(), grb, v1.DeleteOptions{})
	if err2 != nil {
		return fmt.Errorf("failed to delete GlobalRoleBinding: %v", err2)
	}
	fmt.Printf("GlobalRoleBinding deleted: %s\n", grb)
	return nil
}

func GetGlobalRoleBinding(client dynamic.Interface, OwnerDN string) (string, error) {
	// Define the GlobalRoleBinding CRD we want to get
	grbGVR := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "globalrolebindings",
	}

	labelSelector := labels.Set(map[string]string{"pwck8s/ownerdn": OwnerDN}).AsSelector().String()
	listOptions := v1.ListOptions{LabelSelector: labelSelector}
	// Get the user in Rancher
	userList, err := client.Resource(grbGVR).Namespace("").List(context.TODO(), listOptions)
	if err != nil {
		return "", fmt.Errorf("failed to list globalrolebindings: %v", err)
	}

	// If there are multiple users with the same OwnerDN, return an error
	if len(userList.Items) > 1 {
		return "", fmt.Errorf("multiple globalrolebindings found")
	}
	// If there are no users with the same OwnerDN, return an error
	if len(userList.Items) == 0 {
		return "", fmt.Errorf("no globalrolebindings found")
	}
	grb := userList.Items[0]
	return grb.GetName(), nil
}

func CreateProjectRoleBinding(client dynamic.Interface, UserID string, project Project, projectRoleName string) error {
	// Define the User CRD we want to create
	prbGVR := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "projectroletemplatebindings",
	}

	// Define the ProjectRoleBinding CRD we want to create
	prb := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "management.cattle.io/v3",
			"kind":       "ProjectRoleTemplateBinding",
			"metadata": map[string]interface{}{
				"name": "pwck8s-project-owner",
				"labels": map[string]string{
					"pwck8s/userid":         UserID,
					"pwck8s/userdn":         project.OwnerDN,
					"pwck8s/ownerdn":        project.OwnerDN,
					"pwck8s/creationtime":   project.CreationTime.Format("2006-01-02T15-04-05Z07-00"),
					"pwck8s/expirationtime": project.ExpirationTime.Format("2006-01-02T15-04-05Z07-00"),
				},
			},
			"projectName":       project.ClusterID + ":" + project.ProjectID,
			"roleTemplateName":  projectRoleName,
			"userPrincipalName": "local://" + UserID,
			"userName":          UserID,
		},
	}

	// Create the PRB in Rancher
	_, err := client.Resource(prbGVR).Namespace(project.ClusterID).Create(context.TODO(), prb, v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create ProjectRoleBinding: %v", err)
	}
	fmt.Printf("ProjectRoleBinding created: %s\n", UserID)
	return nil
}
