package api

import "k8s.io/client-go/dynamic"

type GlobalConfig struct {
	Client             dynamic.Interface
	ClusterID          string
	DefaultProjectRole string
	DefaultGlobalRole  string
	Debug              bool
}
