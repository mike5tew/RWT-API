#!/bin/bash
set -euo pipefail

# --- Simplified Buildx Setup ---
BUILDX_BUILDER_NAME="multiarch-builder"
# Ensure buildx is available (simple check)
if ! docker buildx version &>/dev/null; then
    echo "ERROR: docker buildx is not available. Please install/enable it."
    # Add specific instructions if needed, e.g., link to Docker docs
    exit 1
fi
# Create or use the builder instance
if ! docker buildx use "${BUILDX_BUILDER_NAME}" >/dev/null 2>&1; then
    echo "Creating buildx builder instance: ${BUILDX_BUILDER_NAME}..."
    docker buildx create --use --name "${BUILDX_BUILDER_NAME}" || {
        echo "ERROR: Failed to create buildx builder instance."
        exit 1
    }
fi
docker buildx inspect --bootstrap || {
    echo "ERROR: Failed to bootstrap buildx builder."
    exit 1
}
echo "Using buildx builder: $(docker buildx current)"
BUILD_COMMAND="docker buildx build --platform linux/amd64"
# --- End Buildx Setup ---

IMAGE_NAME="m1ke57ew/espthinking-api"
BUILD_DATE=$(date +%Y%m%d-%H%M%S)  # More readable timestamp
AMD64_TAG="${IMAGE_NAME}:amd64"
TIMESTAMP_TAG="${IMAGE_NAME}:${BUILD_DATE}"

echo "===== API IMAGE BUILD PROCESS ====="
echo "Building for: AMD64 (Vultr production)"
echo "Build timestamp: ${BUILD_DATE}"

# Optional cleanup
if [[ "${1:-}" == "--clean" ]]; then
    echo "Performing cleanup..."
    docker system prune -af --volumes
fi

# Verify Docker context
if ! docker context inspect &>/dev/null; then
    echo "ERROR: Docker context not properly configured"
    exit 1
fi

# Check SSH key (more robust check)
SSH_KEY="${HOME}/.ssh/id_rsa"
if [[ ! -f "$SSH_KEY" ]]; then
    echo "WARNING: SSH key not found at ${SSH_KEY}"
    echo "Note: Needed for server deployment but not for Docker builds"
fi

# Build the image with proper architecture
echo "Building API image..."
cd ../RWTAPI || { echo "ERROR: Could not find RWTAPI directory"; exit 1; }

# Use the determined BUILD_COMMAND
${BUILD_COMMAND} \
    --load \
    --no-cache \
    --build-arg GOAMD64=v2 \
    -t "${AMD64_TAG}" \
    -t "${TIMESTAMP_TAG}" \
    . || { echo "ERROR: Build failed"; exit 1; }

# --- Removed incompatible binary test ---
echo "Skipping binary test (incompatible with distroless image)."

# Verify Docker Hub authentication
echo "Checking Docker Hub access..."
if ! docker pull alpine:latest >/dev/null 2>&1; then
    echo "ERROR: Cannot access Docker Hub. Please run 'docker login'"
    exit 1
fi

# Push images with retry logic
echo "Pushing images to Docker Hub..."
for tag in "${AMD64_TAG}" "${TIMESTAMP_TAG}"; do
    for attempt in {1..3}; do
        if docker push "${tag}"; then
            break
        else
            echo "WARNING: Push attempt ${attempt} failed"
            if [[ ${attempt} -lt 3 ]]; then
                sleep $((attempt * 2))
            else
                echo "ERROR: Failed to push ${tag} after 3 attempts"
                exit 1
            fi
        fi
    done
done

echo "===== BUILD SUCCESSFUL ====="
echo "Production image: ${AMD64_TAG}"
echo "Versioned image: ${TIMESTAMP_TAG}"
echo ""
echo "Deployment checklist:"
echo "1. Test the image locally: docker run -it --rm ${AMD64_TAG}"
echo "2. Deploy to server:"
echo "   a. cd .. && ./scp-deploy.sh"
echo "   b. ssh your-server './force-cleanup.sh && ./deploy.sh'"
echo "3. Verify deployment: curl https://your-domain.com/api/health"