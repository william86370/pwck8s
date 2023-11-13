# pwck8s - Kubernetes Interface in Go

`pwck8s` is a Go application designed to interface with Kubernetes. It provides API endpoints to handle various Kubernetes-related operations and features a health check endpoint for Kubernetes integration.

## Features

- **Kubernetes Connection**: Connects to Kubernetes either using a local kubeconfig file or in-cluster configuration.
- **API Endpoints**: Provides endpoints for handling projects and users in a Kubernetes cluster.
- **Health Check**: Includes a health check endpoint for Kubernetes liveness and readiness probes.
- **Environment Configuration**: Configurable via environment variables.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go 1.x or higher.
- Access to a Kubernetes cluster.
- Configured kubeconfig file for local development.

### Installing

2. **Run the Application**:
- For local development:
  ```
  ./pwck8s -debug=true
  ```
- For in-cluster execution:
  ```
  ./pwck8s
  ```

### Configuration

Configure the application using the following environment variables:

- `CLUSTER_ID`: ID of the Kubernetes cluster.
- `AUTH_PROVIDER`: Authentication provider for the cluster.
- `DEFAULT_PROJECT_ROLE`: Default role for a project.
- `DEFAULT_GLOBAL_ROLE`: Default global role.
```go
	GlobalConfig := api.GlobalConfig{
		Client:             dynamicClient,
		ClusterID:          "local",
		AuthProvider:       "local",
		DefaultProjectRole: "project-member",
		DefaultGlobalRole:  "user-base",
		Debug:              *debug,
  }
```

## Usage

The application provides several endpoints:

- `/api/v1/project`: Endpoint for project-related operations.
- `/api/v1/user`: Endpoint for user-related operations.
- `/healthcheck`: Health check endpoint for Kubernetes.

## Kubernetes Integration

For Kubernetes deployments, use the `/healthcheck` endpoint in your liveness and readiness probes.