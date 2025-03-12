PROJECT_NAME := bedrock-sidecar

# VERSION is the version we should download and use.
VERSION:=$(shell git rev-parse --short HEAD)
# DOCKER is the docker image repo we need to push to.
DOCKER_REPO:=defangio
DOCKER_USER:=defangio
DOCKER_IMAGE_NAME:=$(DOCKER_REPO)/$(PROJECT_NAME)

DOCKER_IMAGE_ARM64:=$(DOCKER_IMAGE_NAME):arm64-$(VERSION)
DOCKER_IMAGE_AMD64:=$(DOCKER_IMAGE_NAME):amd64-$(VERSION)
BUILD_FLAGS:=-ldflags "-s -w -X main.version=$(VERSION)" -trimpath

.PHONY: dependencies
dependencies:
	go mod tidy

# Build the application
build: dependencies $(shell find . -name "*.go")
	go build -o bin/${PROJECT_NAME}

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run the server
run: build
	go run .

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin/
	go clean

# Run linter (requires golangci-lint)
.PHONY: lint
lint:
	golangci-lint run 

.PHONY: ensure
ensure: ## Run go get -u
	go get -t -u ./...

.PHONY: build-amd64
build-amd64: ensure
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o bin/${PROJECT_NAME}

.PHONY: build-arm64
build-arm64: ensure
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o bin/${PROJECT_NAME}

.PHONY: image-amd64
image-amd64: build-amd64
	docker build --platform linux/amd64 -t ${PROJECT_NAME} -t $(DOCKER_IMAGE_AMD64) --provenance false .

.PHONY: image-arm64
image-arm64: build-arm64
	docker build --platform linux/arm64 -t ${PROJECT_NAME} -t $(DOCKER_IMAGE_ARM64) --provenance false .

.PHONY: image
images: image-amd64 image-arm64 ## Build all docker images and manifest

.PHONY: push-images
push-images: images login ## Push all docker images
	docker push $(DOCKER_IMAGE_AMD64)
	docker push $(DOCKER_IMAGE_ARM64)
	docker manifest create --amend $(DOCKER_IMAGE_NAME):$(VERSION) $(DOCKER_IMAGE_AMD64) $(DOCKER_IMAGE_ARM64)
	docker manifest create --amend $(DOCKER_IMAGE_NAME):latest $(DOCKER_IMAGE_AMD64) $(DOCKER_IMAGE_ARM64)
	docker manifest push --purge $(DOCKER_IMAGE_NAME):$(VERSION)
	docker manifest push --purge $(DOCKER_IMAGE_NAME):latest

.PHONY: login
login: ## Login to docker
	@docker login -u $(DOCKER_USER)