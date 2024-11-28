FROM golang:1.22

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
COPY . ./
RUN go build cmd/server/main.go

EXPOSE 8081
CMD ["./main"]
