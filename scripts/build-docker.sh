#!/usr/bin/env bash
set -e

# Usage:
# ./build-docker.sh prod    -> build production image
# ./build-docker.sh dev     -> build development image
# ./build-docker.sh all     -> build both images

MODE=${1:-prod}  # default to prod if no argument

DOCKER_USERNAME="shahadathhs"
PACKAGE_NAME="knowledge-capsule-api"
PACKAGE_VERSION="latest"
DOCKERFILE_PROD="Dockerfile"
DOCKERFILE_DEV="Dockerfile.dev"

build_prod() {
    APP_IMAGE="$DOCKER_USERNAME/$PACKAGE_NAME:$PACKAGE_VERSION"
    echo "🐳 Building production Docker image: $APP_IMAGE"
    docker build -f "$DOCKERFILE_PROD" -t "$APP_IMAGE" .
    echo "✅ Production image built: $APP_IMAGE"
}

build_dev() {
    APP_IMAGE_DEV="$DOCKER_USERNAME/$PACKAGE_NAME-dev:$PACKAGE_VERSION"
    echo "🐳 Building development Docker image: $APP_IMAGE_DEV"
    docker build -f "$DOCKERFILE_DEV" -t "$APP_IMAGE_DEV" .
    echo "✅ Development image built: $APP_IMAGE_DEV"
}

case "$MODE" in
    prod)
        build_prod
        ;;
    dev)
        build_dev
        ;;
    all)
        build_prod
        build_dev
        ;;
    *)
        echo "❌ Unknown mode: $MODE"
        echo "Usage: $0 [prod|dev|all]"
        exit 1
        ;;
esac
