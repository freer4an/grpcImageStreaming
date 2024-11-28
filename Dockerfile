FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o cmd/server/main.go

FROM alpine:latest

WORKDIR app/
COPY --from=builder /app .
EXPOSE 8081
CMD ["./app"]
