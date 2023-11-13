# Set the default target to 'all'
.PHONY: all
all: build-prod

# Build all Docker image's for production
.PHONY: build-prod
build-prod:
	cd ./x509-proxy && make build-prod
	cd ./pwck8s && make build-prod

# Clean up
.PHONY: clean
clean:
	cd ./x509-proxy && make clean
	cd ./pwck8s && make clean