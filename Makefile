# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOLINT = golangci-lint
DOCKER = docker

# Directories
SRC_DIR = ./src
APP_NAME = myapp
APP_VERSION = 1.0.0
BIN_DIR = ./bin
COVERAGE_DIR = ../coverage

# Docker parameters
DOCKER_IMAGE_NAME = myapp
DOCKER_IMAGE_TAG = $(APP_VERSION)
DOCKERFILE_PATH = ./Dockerfile

.PHONY: all test lint build docker

all: test lint build docker

test:
	@echo "Running tests..."
	@pushd $(SRC_DIR) && \
	$(GOTEST) -race -coverprofile=$(COVERAGE_DIR)/coverage.out

lint:
	@echo "Running linter..."
	$(GOLINT) run --fix $(SRC_DIR)/...

build:
	@echo "Building binary..."
	$(GOBUILD) -o $(BIN_DIR)/$(APP_NAME) $(SRC_DIR)/...

docker:
	@echo "Building Docker image..."
	$(DOCKER) build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) -f $(DOCKERFILE_PATH) .

clean:
	@echo "Cleaning target folder..."
	rm target/main
