FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy all source files first
COPY . .

# Download dependencies and create go.sum
RUN go mod tidy && go mod download

# Copy source code
# COPY . .

# Build the application
RUN go build -o api-gateway .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/api-gateway .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./api-gateway"]
