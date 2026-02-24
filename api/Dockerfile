# Build stage
FROM golang:1.24.2-alpine3.20 AS builder

# Add maintainer info
LABEL maintainer="Parvej <your-email@example.com>"

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

# Run final stage
FROM alpine:3.20

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create a non-root user and logs directory
RUN adduser -D -g '' appuser && \
    mkdir -p /app/logs && \
    chown -R appuser:appuser /app

USER appuser

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose ports
EXPOSE 9090

# Command to run the executable
CMD ["./main", "serve-rest"]