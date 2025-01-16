#!/bin/bash

IMAGE_NAME="m1ke57ew/espthinking-api"
BUILD_DATE=$(date +%s)
TIMESTAMP_TAG="${IMAGE_NAME}:${BUILD_DATE}"
LATEST_TAG="${IMAGE_NAME}:latest"

echo "Starting build process for ${IMAGE_NAME}..."

# More aggressive cleanup
echo "Cleaning up Docker system..."
docker system prune -af --volumes
docker buildx prune -af
docker builder prune -af

# Remove all buildx builders except default
echo "Resetting buildx environment..."
docker buildx ls | grep -v default | awk '{print $1}' | xargs -r docker buildx rm

# Create fresh builder
echo "Creating fresh builder..."
docker buildx create --use --name fresh-builder
docker buildx inspect --bootstrap

# Stop and remove all containers using the image
echo "Stopping and removing existing containers..."
docker ps -a | grep "${IMAGE_NAME}" | awk '{print $1}' | xargs -r docker rm -f

# Remove existing images
echo "Removing existing images..."
docker images | grep "${IMAGE_NAME}" | awk '{print $3}' | xargs -r docker rmi -f

# Build and push with enhanced logging
echo "Building and pushing image..."
docker buildx build \
  --no-cache \
  --pull \
  --platform linux/amd64 \
  --build-arg CACHE_BUST=$(date +%s) \
  --build-arg GO_VERSION=$(date +%s) \
  -t "${TIMESTAMP_TAG}" \
  -t "${LATEST_TAG}" \
  --build-arg REACT_APP_API_URL=/api \
  --push \
  .

if [ $? -eq 0 ]; then
    echo "Build successful! Tags created:"
    echo "  - ${TIMESTAMP_TAG}"
    echo "  - ${LATEST_TAG}"
else
    echo "Build failed!"
    exit 1
fi
