# Docker settings
PACKAGE_NAME := knowledge-capsule-api
DOCKER_USERNAME := shahadathhs
PACKAGE_VERSION := latest
DOCKERFILE := Dockerfile
DOCKERFILE_DEV := Dockerfile.dev
APP_IMAGE := $(DOCKER_USERNAME)/$(PACKAGE_NAME):$(PACKAGE_VERSION)
APP_IMAGE_DEV := $(DOCKER_USERNAME)/$(PACKAGE_NAME)-dev:$(PACKAGE_VERSION)
COMPOSE_FILE := compose.yaml

# Go / build
BINARY_NAME := server
BUILD_DIR := tmp
BUILD_OUT := $(BUILD_DIR)/$(BINARY_NAME)
GO := go

# Local tools (install into ./.bin by default)
GOBIN ?= $(CURDIR)/.bin
BIN_DIR := $(GOBIN)

# Convenience
.PHONY: all help install hooks run build-local build build-dev push push-dev \
	clean fmt vet tidy up down up-dev down-dev restart restart-dev logs logs-dev \
	containers volumes networks images g-jwt sync-env swagger

all: build-local

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  make all / build-local      Build local binary (default)"
	@echo "  make build                  Build production Docker image ($(APP_IMAGE))"
	@echo "  make build-dev              Build development Docker image ($(APP_IMAGE_DEV))"
	@echo "  make push                   Push production image ($(APP_IMAGE))"
	@echo "  make push-dev               Push development image ($(APP_IMAGE_DEV))"
	@echo "  make up                     Start production compose profile"
	@echo "  make down                   Stop production compose profile"
	@echo "  make up-dev                 Start development compose profile"
	@echo "  make down-dev               Stop development compose profile"
	@echo "  make run                    Run local dev server with air (uses local .bin)"
	@echo "  make hooks                  Install git hooks (lefthook)"
	@echo "  make install                Install dev tools into $(GOBIN)"
	@echo "  make fmt                    Run go fmt ./..."
	@echo "  make vet                    Run go vet ./..."
	@echo "  make tidy                   Run go mod tidy"
	@echo "  make clean                  Remove build artifacts and our Docker images"
	@echo "  make containers             docker compose ps"
	@echo "  make volumes                docker volume ls"
	@echo "  make networks               docker network ls"
	@echo "  make images                 docker compose images"
	@echo "  make swagger                Generate Swagger documentation"

# -------------------------
# Generating JWT secret
# -------------------------
g-jwt:
	@echo "🔑 Generating JWT secret..."
	@./scripts/generate-jwt-secret.sh

# -------------------------
# Environment sync
# -------------------------
SYNC_ENV_SCRIPT := ./scripts/sync-env.sh

sync-env:
	@echo "🌱 Syncing .env variables to environment..."
	@echo "source $(SYNC_ENV_SCRIPT)"

# -------------------------
# Dev tools
# -------------------------
install:
	@echo "📦 Creating bin dir: $(GOBIN)"
	@mkdir -p "$(GOBIN)"
	@echo "⬇️  Installing dev tools into $(GOBIN)..."
	@GOBIN="$(GOBIN)" $(GO) install github.com/air-verse/air@latest
	@GOBIN="$(GOBIN)" $(GO) install github.com/evilmartians/lefthook@latest
	@GOBIN="$(GOBIN)" $(GO) install github.com/swaggo/swag/cmd/swag@latest
	@echo "✅ Installed (air, lefthook, swag) to $(GOBIN). Add $(GOBIN) to PATH to run them globally."

hooks: install
	@echo "🔧 Installing git hooks..."
	@$(GOBIN)/lefthook install || lefthook install

run: install sync-env
	@echo "🚀 Starting API (with live reload locally)..."
	@$(GOBIN)/air || air

# -------------------------
# Build & test
# -------------------------
build-local: sync-env
	@echo "🔨 Building local binary -> $(BUILD_OUT)"
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(GOFLAGS) -o $(BUILD_OUT) $(LDFLAGS) $(BUILD_FLAGS) ./main.go
	@echo "✅ Built: $(BUILD_OUT)"

fmt:
	@echo "🧹 Formatting code..."
	@$(GO) fmt ./...

vet:
	@echo "🔍 Running go vet..."
	@$(GO) vet ./...

tidy:
	@echo "🧹 Running go mod tidy..."
	@$(GO) mod tidy

# -------------------------
# Swagger
# -------------------------
swagger:
	@echo "📝 Generating Swagger docs..."
	@$(GOBIN)/swag init -g main.go --output docs
	@echo "✅ Swagger docs generated in docs/"

# -------------------------
# Docker images (prod/dev)
# -------------------------
build:
	@echo "🐳 Building production Docker image: $(APP_IMAGE)"
	@docker build -f $(DOCKERFILE) -t $(APP_IMAGE) .

build-dev:
	@echo "🐳 Building development Docker image: $(APP_IMAGE_DEV)"
	@docker build -f $(DOCKERFILE_DEV) -t $(APP_IMAGE_DEV) .

push: build
	@echo "📤 Pushing Docker image: $(APP_IMAGE)"
	@docker push $(APP_IMAGE)

push-dev: build-dev
	@echo "📤 Pushing Docker image: $(APP_IMAGE_DEV)"
	@docker push $(APP_IMAGE_DEV)

# -------------------------
# Docker compose (profiles)
# -------------------------
up:
	@echo "🐳 Starting Docker Compose For Production..."
	@docker compose -f $(COMPOSE_FILE) --profile prod up -d --build

down:
	@echo "🛑 Stopping Docker Compose For Production..."
	@docker compose -f $(COMPOSE_FILE) --profile prod down

up-dev:
	@echo "🐳 Starting Docker Compose For Development..."
	@docker compose -f $(COMPOSE_FILE) --profile dev up --build

down-dev:
	@echo "🛑 Stopping Docker Compose For Development..."
	@docker compose -f $(COMPOSE_FILE) --profile dev down

restart: down up

restart-dev: down-dev up-dev

logs:
	@echo "📜 Following compose logs..."
	@docker compose -f $(COMPOSE_FILE) --profile prod logs -f

logs-dev:
	@echo "📜 Following compose logs..."
	@docker compose -f $(COMPOSE_FILE) --profile dev logs -f

containers:
	@echo "📦 Listing Docker containers (compose)..."
	@docker compose ps -a

volumes:
	@echo "📦 Listing Docker volumes..."
	@docker volume ls

networks:
	@echo "🌐 Listing Docker networks..."
	@docker network ls

images:
	@docker compose -f $(COMPOSE_FILE) images

# -------------------------
# Cleanup
# -------------------------
clean: down down-dev
	@echo "🧼 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR) $(BIN_DIR)
	@echo "🧽 Removing docker images (if exist): $(APP_IMAGE) $(APP_IMAGE_DEV)"
	-@docker rm $(shell docker ps -a -q) || true
	-@docker rmi $(APP_IMAGE) || true
	-@docker rmi $(APP_IMAGE_DEV) || true
	@echo "✅ Clean complete"
