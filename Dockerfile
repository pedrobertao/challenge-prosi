FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o blog-api app/cmd/main.go

FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/blog-api .

# Copy .env file (optional - for default values)
COPY .env* ./

EXPOSE 3030

CMD ["./blog-api"]