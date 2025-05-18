# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o mailrelay-api .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from build stage
COPY --from=builder /app/mailrelay-api .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./mailrelay-api"]
