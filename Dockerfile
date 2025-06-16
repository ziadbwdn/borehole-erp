# Stage 1: Builder
# Use a Go base image with necessary build tools
FROM golang:1.24.2-alpine AS builder

# Set working directory for the builder stage
WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
# CGO_ENABLED=0 is important for creating a statically linked binary
# -o /usr/local/bin/server specifies the output path for the executable
# ./cmd/server specifies the main package of your application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /usr/local/bin/server ./cmd/server

# Stage 2: Runner
# Use a minimal base image for the final runtime
FROM alpine:latest

# Set working directory for the runner stage
WORKDIR /root/

# Expose the port your application listens on
EXPOSE 8080

# Copy the compiled executable from the builder stage
COPY --from=builder /usr/local/bin/server .

# Copy configuration file(s)
# Assuming your config.go is needed at runtime for ENV vars or similar setup
COPY config/config.go ./config/

# Define environment variables for the application (e.g., database connection)
# These will be overridden by docker-compose.yaml
ENV DB_HOST="localhost"
ENV DB_PORT="3306"
ENV DB_USER="root"
ENV DB_PASSWORD="iwakpeyek23"
# IMPORTANT: Change this to a strong, random key in production!
ENV DB_NAME="borehole_db"
# IMPORTANT: Change this to a strong, random key in production!
# Note: The secret must be enclosed in double quotes, and any literal '$' must be escaped with '$$'.
ENV JWT_SECRET="@iwakpeyeK23_$$eg00p3ceL"

# Command to run the executable
CMD ["./server"]
