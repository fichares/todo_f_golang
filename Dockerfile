FROM golang:1.23.1 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o site main.go
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/awesomeProject .
EXPOSE 8080
CMD ["./awesomeProject"]