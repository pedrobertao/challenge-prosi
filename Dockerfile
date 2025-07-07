# -------- BUILD STAGE --------
FROM golang:1.23-alpine AS builder

# Install git and CA certificates
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy entire project
COPY . .

# Build the Go binary (statically linked for Alpine)
# Change to cmd directory and build from there
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# -------- RUNTIME STAGE --------
FROM alpine:latest

# Install certs for HTTPS
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy built binary from builder
COPY --from=builder /app/cmd/server .

# Make binary executable
RUN chmod +x ./server

# Set binary as entrypoint
CMD ["./server"]

EXPOSE 8080