# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files from backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code from backend
COPY backend/ .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o marketplace .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies for SQLite
RUN apk add --no-cache libc6-compat

# Copy binary from builder
COPY --from=builder /app/marketplace .

# Create uploads directory
RUN mkdir -p uploads

# Expose port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release
ENV PORT=8080

# Run the application
CMD ["./marketplace"]
