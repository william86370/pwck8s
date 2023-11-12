package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"net/http"
	"pwck8s/api"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/fatih/color"
)

func printInBox(lines []string) {
	maxLength := 0
	for _, line := range lines {
		if len(line) > maxLength {
			maxLength = len(line)
		}
	}

	topBottom := "+" + strings.Repeat("-", maxLength+2) + "+"
	color.Cyan(topBottom)
	for _, line := range lines {
		spaces := strings.Repeat(" ", maxLength-len(line))
		color.Cyan("| %s%s |", line, spaces)
	}
	color.Cyan(topBottom)
}

func prettyLogBox(title string, content map[string]string) string {
	var maxLength int = len(title)
	for key, value := range content {
		length := len(key + ": " + value)
		if length > maxLength {
			maxLength = length
		}
	}
	topBottomBorder := "+" + strings.Repeat("-", maxLength+2) + "+"
	titleLine := fmt.Sprintf("| %s%s |", title, strings.Repeat(" ", maxLength-len(title)))

	var contentLines string
	for key, value := range content {
		contentLines += fmt.Sprintf("| %s: %s%s |\n", key, value, strings.Repeat(" ", maxLength-len(key+": "+value)))
	}
	return fmt.Sprintf("%s\n%s\n%s%s", topBottomBorder, titleLine, contentLines, topBottomBorder)
}

func main() {
	var kubeconfig *string
	debug := flag.Bool("debug", false, "Set to true to use kubeconfig for local debugging")

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	var config *rest.Config
	var err error

	if *debug {
		// Use kubeconfig from local file
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			color.Red("Error reading kubeconfig: %v", err)
			return
		}
		color.Green("Connected to Kubernetes using local kubeconfig")
	} else {
		// Use in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			color.Red("Error getting in-cluster config: %v", err)
			return
		}
		color.Green("Connected to Kubernetes using in-cluster config")
	}

	// Create a new clientset which includes our CRD schema
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating clientset: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating dynamic client: %v", err)
	}

	// Fetch Kubernetes version info
	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		color.Red("Failed to get Kubernetes server version: %v", err)
		return
	}

	infoLines := []string{
		fmt.Sprintf("Kubernetes Server Version: %s", version.String()),
		fmt.Sprintf("Major: %s, Minor: %s", version.Major, version.Minor),
	}

	printInBox(infoLines)

	// Setup HTTP server and handlers
	http.HandleFunc("/api/v1/project", func(w http.ResponseWriter, r *http.Request) {
		api.ProjectHandler(dynamicClient, w, r)
	})

	http.HandleFunc("/api/v1/user", func(w http.ResponseWriter, r *http.Request) {
		api.UserHandler(dynamicClient, w, r)
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
