# Build stage
FROM eclipse-temurin:17-jdk-alpine AS build

WORKDIR /workspace/app

# Install Maven and required packages
RUN apk add --no-cache wget maven

# Copy pom.xml first for better layer caching
COPY pom.xml .

# Download dependencies (cached layer)
RUN mvn dependency:go-offline -B

# Copy source code
COPY src src

# Build the application
RUN mvn clean package -DskipTests -Dmaven.javadoc.skip=true

# Runtime stage
FROM eclipse-temurin:17-jre-alpine

# Install required packages
RUN apk add --no-cache wget curl tzdata \
    && cp /usr/share/zoneinfo/Asia/Almaty /etc/localtime \
    && echo "Asia/Almaty" > /etc/timezone \
    && apk del tzdata

# Create app user for security
RUN addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

# Set working directory
WORKDIR /app

# Copy jar file from build stage
COPY --from=build --chown=appuser:appuser /workspace/app/target/*.jar app.jar

# Switch to app user
USER appuser

# Expose port
EXPOSE 8081

# Add JVM options for better performance and monitoring
ENV JAVA_OPTS="-XX:+UseContainerSupport \
    -XX:MaxRAMPercentage=75.0 \
    -XX:+UseG1GC \
    -XX:+UnlockExperimentalVMOptions \
    -XX:+UseCGroupMemoryLimitForHeap \
    -Djava.security.egd=file:/dev/./urandom \
    -Dspring.profiles.active=prod"

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=120s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/api/actuator/health || exit 1

# Run the application using standard jar execution
ENTRYPOINT ["sh", "-c", "java $JAVA_OPTS -jar app.jar"]