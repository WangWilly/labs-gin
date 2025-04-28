#/bin/bash

cd deployments

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "Docker Compose is not installed. Please install Docker Compose and try again."
    exit 1
fi

# Check if the Docker Compose file exists
if [ ! -f docker-compose.yml ]; then
    echo "Docker Compose file (docker-compose.yml) not found. Please ensure it exists in the current directory."
    exit 1
fi

# Start the Docker containers using Docker Compose
docker-compose up -d
if [ $? -ne 0 ]; then
    echo "Failed to start Docker containers. Please check the errors above."
    exit 1
fi

echo "Docker containers started successfully."
