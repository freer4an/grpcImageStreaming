package services

import (
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"mime/multipart"
)

type S3storage struct {
	client   *s3.Client
	uploader *manager.Uploader
}

func NewS3Storage(client *s3.Client) *S3storage {
	return &S3storage{
		client:   client,
		uploader: manager.NewUploader(client),
	}
}

func (s *S3storage) Upload(file *multipart.FileHeader) (string, error) {
	return "", nil
}
