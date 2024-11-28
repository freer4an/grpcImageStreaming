db_driver=postgres
db_string="postgresql://postgres:pass@localhost:5432/image_storage?sslmode=disable"


proto:
	protoc -I protos/buf/image \
		--go_out=./protos/gen --go_opt=paths=source_relative \
		--go-grpc_out=./protos/gen --go-grpc_opt=paths=source_relative \
		protos/buf/image/*.proto

startapp:
	go run cmd/server/main.go

startclient:
	go run cmd/client/main.go

clearstorage:
	rm -f image_storage/originals/* && rm -f image_storage/thumbnails/*

migup:
	goose -dir migrations ${db_driver} ${db_string} up

migdown:
	goose -dir migrations ${db_driver} ${db_string} down

testListImage:
	grpcurl -d '{}' -plaintext localhost:8081 image.v1.ImageService/ListImages
