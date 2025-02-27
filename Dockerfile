# Build stage
FROM golang:1.23 AS builder
WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN go build -o api ./cmd/api

# Final stage
FROM alpine:3.19
WORKDIR /app

# Add certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/api .

# Expose port and run
EXPOSE 8080
CMD ["./api"]