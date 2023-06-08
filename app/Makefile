include default.env

# GLOBAL ENVIRONMENT VARIABLES #################################################
APP_NAME = todo-api

# GENERAL PURPOSE TARGETS ######################################################
.PHONY: help
help: ## Display available commands.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets.

.PHONY: build
build:
	./scripts/build.sh

.PHONY: build-darwin
build-darwin:
	go mod tidy
	GOOS=darwin GOARCH=arm64 go build -a -o bin/todo-api main.go

.PHONY: run
run: ## Run application with default arugments
	go run main.go --config default.env

.PHONY: clean
clean: ## Clean repository
	@test ! -e bin || rm -r bin

# Docker targets.

VERSION = local
REPOSITORY = docker.io/pedropeixoto6/todo-api
.PHONY: build-docker-local
build-docker-local: ## Build local docker image
	docker build \
		--tag $(REPOSITORY):$(VERSION) \
		-f deploy/local/Dockerfile .

.PHONY: build-docker
build-docker: ## Build docker image
	docker build \
		--tag $(REPOSITORY):$(VERSION) \
		-f Dockerfile .

.PHONY: docker-compose-up-local
docker-compose-up-local: ## Run all services locally using docker compose
	docker-compose -f deploy/local/docker-compose.yaml up -d --build

.PHONY: docker-compose-down-local
docker-compose-down-local: ## Run infra services locally using docker compose
	docker-compose -f deploy/local/docker-compose.yaml down