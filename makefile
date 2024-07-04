# Variables
APP_NAME := noxai
DOCKER_USERNAME := faturfawkes
DOCKER_REPO := $(DOCKER_USERNAME)/$(APP_NAME)
#VERSION := $(shell git describe --tags --always --dirty)
VERSION := latest
DOCKER_IMAGE := $(DOCKER_REPO):$(VERSION)

# Default goal
.PHONY: all
all: build

# Build the Go application
.PHONY: build
build:
	@echo "Building the Go application..."
	GO111MODULE=on go build -o $(APP_NAME) ./...

# Run the application locally
.PHONY: run
run: build
	@echo "Running the application..."
	./$(APP_NAME)

# Clean the build
.PHONY: clean
clean:
	@echo "Cleaning the build..."
	rm -f $(APP_NAME)

# Build the Docker image
.PHONY: docker-build
docker-build:
	@echo "Building the Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Push the Docker image to Docker Hub
.PHONY: docker-push
docker-push: docker-build
	@echo "Pushing the Docker image to Docker Hub..."
	docker push $(DOCKER_IMAGE)

# Clean up Docker images
.PHONY: docker-clean
docker-clean:
	@echo "Cleaning up Docker images..."
	docker rmi $(DOCKER_IMAGE)

# Full deployment process
.PHONY: deploy
deploy: docker-push
	@echo "Deploy complete!"

# Display help
.PHONY: help
help:
	@echo "Makefile targets:"
	@echo "  build        - Build the Go application"
	@echo "  run          - Run the application locally"
	@echo "  clean        - Clean the build"
	@echo "  docker-build - Build the Docker image"
	@echo "  docker-push  - Push the Docker image to Docker Hub"
	@echo "  docker-clean - Clean up Docker images"
	@echo "  deploy       - Full deployment process (build and push Docker image)"
	@echo "  help         - Display this help message"
