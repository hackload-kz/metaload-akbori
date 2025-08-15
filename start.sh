#!/bin/bash

echo "🚀 Starting Biletter Service..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Check if Maven is installed
if ! command -v mvn &> /dev/null; then
    echo "❌ Maven is not installed. Please install Maven first."
    exit 1
fi

# Check if Java 17+ is installed
JAVA_VERSION=$(java -version 2>&1 | head -n 1 | cut -d'"' -f2 | cut -d'.' -f1)
if [ "$JAVA_VERSION" -lt 17 ]; then
    echo "❌ Java 17+ is required. Current version: $JAVA_VERSION"
    exit 1
fi

echo "✅ Prerequisites check passed"

# Start infrastructure services
echo "🐘 Starting PostgreSQL and Kafka..."
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 30

# Check if services are running
if ! docker-compose ps | grep -q "Up"; then
    echo "❌ Failed to start infrastructure services"
    docker-compose logs
    exit 1
fi

echo "✅ Infrastructure services started successfully"

# Build the project
echo "🔨 Building project..."
mvn clean compile

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Build completed"

# Start the application
echo "🎯 Starting Spring Boot application..."
mvn spring-boot:run

echo "🎉 Biletter Service is running at http://localhost:8081"
echo "📊 Kafka UI is available at http://localhost:8080"
echo "🛑 Press Ctrl+C to stop the application"
