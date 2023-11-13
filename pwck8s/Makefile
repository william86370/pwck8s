# Set the default target to 'all'
.PHONY: all
all: build-prod

# Set the Go binary name
BINARY_NAME=pwck8s-backend

# Docker settings
DOCKER_IMAGE_NAME=pwck8s-backend
DOCKER_TAG=latest

# Build the application for development
.PHONY: build-dev
build-dev:
	@echo "Building the application for development..."
	go build -o ${BINARY_NAME}

# Run the application in development
.PHONY: run-dev
run-dev: build-dev
	@echo "Running the application..."
	@PORT=8080 \
	./${BINARY_NAME}

# Build the Docker image for production
.PHONY: build-prod
build-prod:
	@echo "Building the Docker image..."
	docker build -t ${DOCKER_IMAGE_NAME}:${DOCKER_TAG} .

# Clean up
.PHONY: clean
clean:
	@echo "Cleaning up..."
	go clean || true
	rm -f ${BINARY_NAME} || true
	docker rmi ${DOCKER_IMAGE_NAME}:${DOCKER_TAG} || true