# Stage 1: Build the application
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk --no-cache add ca-certificates git

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build a static binary with CGO disabled for better portability
# This creates a standalone executable that doesn't depend on external libraries
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -trimpath \
    -o /receipt-processor

# Stage 2: Create a minimal runtime image
FROM gcr.io/distroless/static-debian12

# Copy the binary from the builder stage
COPY --from=builder /receipt-processor /receipt-processor

# Run as non-root user for security
USER nonroot:nonroot

# Expose the port that the application listens on
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/receipt-processor", "-health-check"]

# Command to run the executable
CMD ["/receipt-processor"]