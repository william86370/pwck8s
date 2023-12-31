# Set the default target to 'all'
.PHONY: all
all: build-prod

# Set the Go binary name
BINARY_NAME=x509-proxy

# Docker settings
DOCKER_IMAGE_NAME=x509-proxy
DOCKER_TAG=latest


.PHONY: build-certs
build-certs:
	@echo "Building the certificates..."
	@cd ../certificate-toolkit && make all

# Build the application for development
.PHONY: build-dev
build-dev:
	@echo "Building the application for development..."
	go build -o ${BINARY_NAME}

# Run the application in development
.PHONY: run-dev
run-dev: build-dev build-certs
	@echo "Running the application..."
	@PORT=8443 \
	TLS_CERT=../certificate-toolkit/server/server.crt \
	TLS_KEY=../certificate-toolkit/server/server.key \
	CA_CERT=../certificate-toolkit/ca/ca.crt \
	PROXY_URL=localhost:8080 \
	DEBUG=true \
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
	cd ../certificate-toolkit && make clean