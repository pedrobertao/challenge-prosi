# -------- BUILD STAGE --------
FROM golang:1.23-alpine AS builder

# Install dependencies for HTTPS, Git
RUN apk add --no-cache git ca-certificates

# Set root of your project
WORKDIR /build

# Copy everything (go.mod must be at root)
COPY . .

# Download Go modules
RUN go mod download

# Build from full path to main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./app/cmd/main.go

# -------- RUNTIME STAGE --------
FROM alpine:latest

# Install certs for HTTPS
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy binary from build stage
COPY --from=builder /build/server .

# Expose port for Cloud Run
EXPOSE 8080

# Run the binary
CMD ["./server"]
