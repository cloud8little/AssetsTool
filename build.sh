#!/bin/bash

# Set the project name
PROJECT_NAME="assetcli"

# Build for AMD64
echo "Building for AMD64..."
GOARCH=amd64 GOOS=linux go build -o ${PROJECT_NAME}_amd64
echo "AMD64 build completed."

# Build for ARM64
echo "Building for ARM64..."
GOARCH=arm64 GOOS=linux go build -o ${PROJECT_NAME}_arm64
echo "ARM64 build completed."

echo "All builds completed."
