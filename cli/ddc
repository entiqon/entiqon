#!/bin/bash
set -e

# Default values
IMAGE_NAME="entiqon-docs"
CONTAINER_NAME="EntiqonDocs"
PORT_MAPPING="3001:8000"
DOCKERFILE_NAME="Dockerfile-documentation"

# Parse command line args for --image, --container, --port
while [[ $# -gt 0 ]]; do
  case $1 in
    --image)
      IMAGE_NAME="$2"
      shift 2
      ;;
    --container)
      CONTAINER_NAME="$2"
      shift 2
      ;;
    --port)
      PORT_MAPPING="$2"
      shift 2
      ;;
    *)
      echo "Unknown option $1"
      exit 1
      ;;
  esac
done

echo "Building Docker image: $IMAGE_NAME:latest"
docker build -f $DOCKERFILE_NAME -t $IMAGE_NAME:latest .

echo "Checking for running container: $CONTAINER_NAME"
if [ "$(docker ps -q -f name=^/${CONTAINER_NAME}$)" ]; then
    echo "Stopping container $CONTAINER_NAME"
    docker stop $CONTAINER_NAME
fi

if [ "$(docker ps -aq -f status=exited -f name=^/${CONTAINER_NAME}$)" ]; then
    echo "Removing container $CONTAINER_NAME"
    docker rm $CONTAINER_NAME
fi

echo "Starting new container $CONTAINER_NAME on port $PORT_MAPPING"
docker run -d -p $PORT_MAPPING --name $CONTAINER_NAME $IMAGE_NAME:latest

echo "Deployment complete. Container $CONTAINER_NAME is running."
