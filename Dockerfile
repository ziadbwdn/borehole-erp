# Stage 1: Builder
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server ./cmd/server

# Stage 2: Runner
FROM alpine:latest
RUN apk update && apk add --no-cache busybox-extras
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080

# CORRECTED: A command that keeps the container alive for debugging
# It will sleep forever, allowing us to exec into it.
CMD ["sleep", "infinity"]