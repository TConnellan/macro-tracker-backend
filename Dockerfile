# Start from the official Golang image as a build stage
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /api ./cmd/api

# Final runtime image
FROM ubuntu:22.04

# Install dependencies in one layer to avoid missing package list errors
RUN apt-get update && \
    apt-get install -y curl ca-certificates tar && \
    rm -rf /var/lib/apt/lists/*

# Install migrate CLI safely
RUN curl -fsSL -o migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz && \
    tar -xvzf migrate.tar.gz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate && \
    rm migrate.tar.gz

# Optional: non-root user
RUN useradd -m appuser
USER appuser

COPY --from=builder /api /api

EXPOSE 8080
ENTRYPOINT ["/api"]
