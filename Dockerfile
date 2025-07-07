# -------- BUILD STAGE --------
FROM golang:1.24-alpine AS builder

# Install git and CA certificates
RUN apk add --no-cache git ca-certificates

WORKDIR /app
# Copy entire project
COPY . .

# Download dependencies
RUN go mod download

# Set to directory where your main.go lives
WORKDIR /app/cmd

# Build the Go binary (statically linked for Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# -------- RUNTIME STAGE --------
FROM alpine:latest

# Install certs for HTTPS
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy built binary from builder
COPY --from=builder /app/cmd/server .

# Set binary as entrypoint
CMD ["./server"]

EXPOSE 8080