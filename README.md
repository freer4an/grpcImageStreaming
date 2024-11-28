# gRPC Image Streaming Service

This repository contains a **gRPC service** for uploading, managing, and retrieving images with associated metadata. The service is designed to handle image streams, generate thumbnails, store metadata in a PostgreSQL database, and store images locally.

## Features

1. **UploadImage**  
   - Accepts a stream of image chunks from the client.
   - Saves the image locally.
   - Generates a thumbnail of the image.
   - Stores metadata about the image in a PostgreSQL database.
   - Processes metadata concurrently for high performance.

2. **ListImages**  
   - Retrieves all image metadata from the database.

3. **GetImage**  
   - Returns an image (original or thumbnail) along with its metadata.

4. **DeleteImage**  
   - Deletes image metadata from the database.

---

## Tech Stack

- **gRPC**: For client-server communication.
- **PostgreSQL**: For metadata storage.
- **Goose & Embed Migrations**: For database schema management.
- **pgxpool**: For database connection pooling.
- **imaging**: For image manipulation (e.g., resizing).
- **grcpurl**: For testing the gRPC service.

---

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) 1.20+
- PostgreSQL database
- `grpcurl` for testing
- `goose` for database migrations
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Clone the Repository

```bash
git clone https://github.com/freer4an/grpcImageStreaming.git
cd grpcImageStreaming
```

# Testing

1. Start the server
    ```bash
    make startapp
    ```
    
    ```text
    go run cmd/server/main.go
    postgresql://postgres:pass@127.0.0.1:5432/image_storage?sslmode=disable
    2024/11/28 14:48:55 OK   20241127152558_add_image_column.sql (4.63ms)
    2024/11/28 14:48:55 goose: successfully migrated database to version: 20241127152558
    2024/11/28 14:48:55 gRPC server started
    ```
2. Start client to send list of images
    ```bash
    make startclient
    ```
    
    ```text
    go run cmd/client/main.go
    2024/11/28 14:51:43 INFO Client listening on 127.0.0.1:8081
    2024/11/28 14:51:44 INFO Successfully uploaded images "total size"=6006901
    ```
    
3. Get all image metada from db
    ```bash
    make listImages
    ``` 
    
    ```text
    grpcurl -d '{}' -plaintext localhost:8081 image.v1.ImageService/ListImages
    {
      "images": [
        {
          "id": "20a84d64-6ee1-4bc2-80b0-0634574dc18e",
          "format": ".jpg",
          "width": 256,
          "height": 256,
          "originalPath": "image_storage/originals/pexels-am83-13407872.jpg",
          "thumbnailPath": "image_storage/thumbnails/20a84d64-6ee1-4bc2-80b0-0634574dc18e/.jpg"
        },
        {
          "id": "1459b1b1-3b6f-4e49-a7a3-febf2dc4c441",
          "format": ".jpg",
          ...
        },
        ...
      ]
    }
    ```

4. Get image and it's metadata 
    ```bash
    grpcurl -d '{"id": "<uuid-from-metadata>"}' -plaintext localhost:8081 image.v1.ImageService/GetImage
    ```
    ```text
    {
      "imageMetada": {
        "id": "1cb2895b-f7db-430d-8dcf-05a2c4c04f11",
        "width": 256,
        "height": 256,
        "originalPath": "image_storage/originals/pexels-werner-hilversum-793171962-21365263.jpg",
        "thumbnailPath": "image_storage/thumbnails/1cb2895b-f7db-430d-8dcf-05a2c4c04f11/.jpg"
      }
    }
    ```
5. Delete image
    ```bash
    grpcurl -d '{"id": "<uuid-from-metadata>"}' -plaintext localhost:8081 image.v1.ImageService/DeleteImage
    ```
