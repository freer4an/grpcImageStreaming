package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/freer4an/image-storage/internal/connections"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	file, err := os.Open("images/pexels-am83-13407872.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ctx := context.Background()
	awsClient := connections.S3client(ctx)

	_, err = awsClient.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(""),
		Key:    aws.String(""),
		Body:   file,
	})
	if err != nil {
		log.Fatal(err)
	}
}
