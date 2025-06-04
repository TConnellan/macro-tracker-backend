# Start from the official Golang image as a build stage
FROM golang:1.22 AS builder

# Set the working directory
WORKDIR /app

# Download the Go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go app (assuming your main.go is in cmd/api/main.go)
RUN go build -o /api ./cmd/api

# Final runtime image
FROM ubuntu:22.04

# Install certificates and dependencies
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Install migrate CLI
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate.linux-amd64 /usr/local/bin/migrate && chmod +x /usr/local/bin/migrate

# Set up non-root user (optional but recommended)
RUN useradd -m appuser
USER appuser

# Copy the built binary from the builder
COPY --from=builder /api /api

# Expose the port your API listens on
EXPOSE 8080

# Run the Go binary
ENTRYPOINT ["/api"]
