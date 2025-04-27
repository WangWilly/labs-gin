#/bin/bash

# Build docker image
docker build -t labs-gin-app:latest .
if [ $? -ne 0 ]; then
    echo "Failed to build the Docker image. Please check the errors above."
    exit 1
fi

echo "Docker image built successfully. Image name: labs-gin-app:latest"
