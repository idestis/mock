.PHONY: build run test clean docker-build docker-run docker-run-custom \
	ghcr-login ghcr-build ghcr-build-push ghcr-push ghcr-pull ghcr-run \
	helm-install helm-upgrade helm-uninstall helm-template help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build parameters
BINARY_NAME=server
BINARY_PATH=bin/$(BINARY_NAME)
MAIN_PATH=cmd/server/main.go

# Docker parameters
DOCKER_IMAGE=mock-backend-jobs
DOCKER_TAG=latest

# GHCR parameters
GHCR_REGISTRY=ghcr.io
GITHUB_USER?=$(shell git config user.name | tr '[:upper:]' '[:lower:]')
GITHUB_REPO?=mock-backend-jobs
GHCR_IMAGE=$(GHCR_REGISTRY)/$(GITHUB_USER)/$(GITHUB_REPO)
VERSION?=latest

# Helm parameters
HELM_RELEASE=mock-backend-jobs
HELM_CHART=./helm/mock-backend-jobs

all: test build

build:
	$(GOBUILD) -o $(BINARY_PATH) $(MAIN_PATH)

run:
	$(GORUN) $(MAIN_PATH)

test:
	$(GOTEST) -v ./...

clean:
	rm -rf bin/
	rm -f $(BINARY_NAME)

deps:
	$(GOMOD) download
	$(GOMOD) tidy

fmt:
	$(GOCMD) fmt ./...

docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-run-custom:
	docker run -p 8080:8080 \
		-e USER_ANNONYMIZATION_DURATION=10s \
		-e USER_ANNONYMIZATION_SHOULD_FAIL=false \
		-e ZENDESK_IMPORT_DURATION=15s \
		-e ZENDESK_IMPORT_SHOULD_FAIL=false \
		-e BLUESHIFT_EXPORT_DURATION=20s \
		-e BLUESHIFT_EXPORT_SHOULD_FAIL=true \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

# GHCR targets
ghcr-login:
	@echo "Logging into GitHub Container Registry..."
	@echo "Make sure GITHUB_TOKEN is set in your environment"
	@echo $$GITHUB_TOKEN | docker login $(GHCR_REGISTRY) -u $(GITHUB_USER) --password-stdin

ghcr-build:
	docker buildx build --platform linux/amd64,linux/arm64 \
		-t $(GHCR_IMAGE):$(VERSION) \
		-t $(GHCR_IMAGE):latest \
		.

ghcr-build-push:
	docker buildx build --platform linux/amd64,linux/arm64 \
		-t $(GHCR_IMAGE):$(VERSION) \
		-t $(GHCR_IMAGE):latest \
		--push \
		.

ghcr-push:
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(GHCR_IMAGE):$(VERSION)
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(GHCR_IMAGE):latest
	docker push $(GHCR_IMAGE):$(VERSION)
	docker push $(GHCR_IMAGE):latest

ghcr-pull:
	docker pull $(GHCR_IMAGE):$(VERSION)

ghcr-run:
	docker run -p 8080:8080 $(GHCR_IMAGE):$(VERSION)

helm-install:
	helm install $(HELM_RELEASE) $(HELM_CHART)

helm-upgrade:
	helm upgrade $(HELM_RELEASE) $(HELM_CHART)

helm-uninstall:
	helm uninstall $(HELM_RELEASE)

helm-template:
	helm template $(HELM_RELEASE) $(HELM_CHART)

help:
	@echo "Available targets:"
	@echo ""
	@echo "Go commands:"
	@echo "  build              - Build the binary"
	@echo "  run                - Run the application"
	@echo "  test               - Run tests"
	@echo "  clean              - Remove build artifacts"
	@echo "  deps               - Download and tidy dependencies"
	@echo "  fmt                - Format code"
	@echo ""
	@echo "Docker commands:"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run Docker container"
	@echo "  docker-run-custom  - Run Docker with custom env vars"
	@echo ""
	@echo "GHCR commands:"
	@echo "  ghcr-login         - Login to GitHub Container Registry"
	@echo "  ghcr-build         - Build multi-arch image for GHCR"
	@echo "  ghcr-build-push    - Build and push multi-arch image to GHCR"
	@echo "  ghcr-push          - Push image to GHCR"
	@echo "  ghcr-pull          - Pull image from GHCR"
	@echo "  ghcr-run           - Run GHCR image"
	@echo ""
	@echo "Helm commands:"
	@echo "  helm-install       - Install Helm chart"
	@echo "  helm-upgrade       - Upgrade Helm release"
	@echo "  helm-uninstall     - Uninstall Helm release"
	@echo "  helm-template      - Render Helm templates"
	@echo ""
	@echo "Environment variables:"
	@echo "  GITHUB_USER        - GitHub username (default: git config user.name)"
	@echo "  GITHUB_TOKEN       - GitHub personal access token (required for GHCR)"
	@echo "  VERSION            - Image version tag (default: latest)"
