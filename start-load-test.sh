#!/bin/bash

echo "🚀 Starting load tests..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

docker run --rm \
  -e API_URL=http://biletter-app:8081 \
  -v "$PWD":/work -w /work \
  --network biletter-net \
  grafana/k6 run events-load-test.js
