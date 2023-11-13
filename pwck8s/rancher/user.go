package rancher

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// User defines the structure for the user
type User struct {
	// User Object
	UserID         string    `json:"userId"`
	DisplayName    string    `json:"displayName"`
	PrincipalIds   []string  `json:"principalIds"`
	UserDN         string    `json:"userDn"`
	CreationTime   time.Time `json:"creationTime"`
	ExpirationTime time.Time `json:"expirationTime"`
}

func GenerateUserId() string {
	// Generate a new User ID similar to Rancher User ID
	// u-<random 5 char string>
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	const letterBytes = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return "pwck8s-" + string(b)
}

func GenerateUser(UserDN string, AuthProvider string) User {
	// Generate a new user object
	UserID := GenerateUserId()
	user := User{
		UserID:         UserID,
		DisplayName:    UserID,
		PrincipalIds:   []string{AuthProvider + "://" + UserID},
		UserDN:         UserDN,
		CreationTime:   time.Now(),
		ExpirationTime: time.Now().AddDate(0, 0, 1),
	}
	return user
}

func CreateRancherUser(client dynamic.Interface, newUser User) error {
	// Define the User CRD we want to create
	userGVR := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "users",
	}
	user := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "management.cattle.io/v3",
			"kind":       "User",
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
			"principalIds": newUser.PrincipalIds,
			"description":  "Created by pwck8s",
			"username":     newUser.DisplayName,
		},
	}

	// Create the user in Rancher
	_, err := client.Resource(userGVR).Namespace("").Create(context.TODO(), user, v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	fmt.Printf("User created: %s\n", newUser.UserID)
	return nil
}

func DeleteRancherUser(client dynamic.Interface, OwnerDN string) error {

	// Define the User CRD we want to delete
	userGVR := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "users",
	}
	// Get the user in Rancher
	user, err := GetRancherUser(client, OwnerDN)
	if err != nil {
		return err
	}

	// Delete the user in Rancher
	err2 := client.Resource(userGVR).Namespace("").Delete(context.TODO(), user.UserID, v1.DeleteOptions{})
	if err2 != nil {
		return fmt.Errorf("failed to delete user: %v", err2)
	}

	fmt.Printf("User deleted: %s\n", user.UserID)
	return nil
}

func GetRancherUser(client dynamic.Interface, OwnerDN string) (User, error) {
	var tmpuser User
	// Define the User CRD we want to get
	userGVR := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "users",
	}

	labelSelector := labels.Set(map[string]string{"pwck8s/ownerdn": OwnerDN}).AsSelector().String()
	listOptions := v1.ListOptions{LabelSelector: labelSelector}
	// Get the user in Rancher
	userList, err := client.Resource(userGVR).Namespace("").List(context.TODO(), listOptions)
	if err != nil {
		return tmpuser, fmt.Errorf("failed to list users: %v", err)
	}

	// If there are multiple users with the same OwnerDN, return an error
	if len(userList.Items) > 1 {
		return tmpuser, fmt.Errorf("multiple users found")
	}
	// If there are no users with the same OwnerDN, return an error
	if len(userList.Items) == 0 {
		return tmpuser, fmt.Errorf("no users found")
	}
	user := userList.Items[0]

	CreationTime, _ := time.Parse("2006-01-02T15-04-05Z07-00", user.Object["metadata"].(map[string]interface{})["labels"].(map[string]interface{})["pwck8s/creationtime"].(string))
	ExpirationTime, _ := time.Parse("2006-01-02T15-04-05Z07-00", user.Object["metadata"].(map[string]interface{})["labels"].(map[string]interface{})["pwck8s/expirationtime"].(string))
	// Convert the user to our User struct
	userStruct := User{
		UserID:      user.Object["metadata"].(map[string]interface{})["name"].(string),
		DisplayName: user.Object["username"].(string),
		// PrincipalIds:   user.Object["principalIds"].([]interface{}([]string)),
		PrincipalIds:   []string{"local://" + user.Object["metadata"].(map[string]interface{})["name"].(string)},
		UserDN:         user.Object["metadata"].(map[string]interface{})["labels"].(map[string]interface{})["pwck8s/userdn"].(string),
		CreationTime:   CreationTime,
		ExpirationTime: ExpirationTime,
	}

	fmt.Printf("User retrieved: %s\n", userStruct.UserID)
	return userStruct, nil
}

func UserExists(client dynamic.Interface, OwnerDN string) (bool, error) {
	// Define the User CRD we want to get
	userGVR := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "users",
	}

	labelSelector := labels.Set(map[string]string{"pwck8s/ownerdn": OwnerDN}).AsSelector().String()
	listOptions := v1.ListOptions{LabelSelector: labelSelector}
	// Get the user in Rancher
	userList, err := client.Resource(userGVR).Namespace("").List(context.TODO(), listOptions)
	if err != nil {
		return false, fmt.Errorf("failed to list users: %v", err)
	}

	// If there are multiple users with the same OwnerDN, return an error
	if len(userList.Items) > 1 {
		return false, fmt.Errorf("multiple users found")
	}
	// If there are no users with the same OwnerDN, return an error
	if len(userList.Items) == 0 {
		return false, nil
	}

	return true, nil
}
