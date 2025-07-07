FROM golang:1.24-apline

WORKDIR /app

COPY . .

WORKDIR /app/cmd
RUN go mod download
RUN go build -o server main.go

EXPOSE 8080
CMD ["./server"]