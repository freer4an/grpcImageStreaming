package services

import (
	"context"
	"fmt"
	"github.com/freer4an/image-storage/protos/gen"
	"google.golang.org/grpc"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type ImageServiceClient struct {
	Service      gen.ImageServiceClient
	imageFormats []string
}

func NewImageClient(cc *grpc.ClientConn, imageFormats []string) *ImageServiceClient {
	imageService := gen.NewImageServiceClient(cc)
	return &ImageServiceClient{
		Service:      imageService,
		imageFormats: imageFormats,
	}
}

func (cl *ImageServiceClient) UploadImage(ctx context.Context, filePaths []string) error {
	slog.Debug("Starting stream")
	stream, err := cl.Service.UploadImage(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't start stream: %w", err)
	}

	for _, filePath := range filePaths {
		if err := cl.sendImage(stream, filePath); err != nil {
			return fmt.Errorf("couldn't send file: %w", err)
		}
	}

	slog.Debug("Closing stream")
	res, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("couldn't receive response: %w", err)
	}

	slog.Info("Successfully uploaded images",
		slog.Int64("total size", int64(res.Size)))
	return nil
}

func (cl *ImageServiceClient) sendImage(stream gen.ImageService_UploadImageClient, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("couldn't open file: %w", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			slog.Error("Error closing file: ", slog.String("err", err.Error()))
		}
	}()

	info, err := cl.validateFile(filePath)
	if err != nil {
		return fmt.Errorf("couldn't read file info: %w", err)
	}

	ext := filepath.Ext(filePath)
	buf := make([]byte, MAX_IMAGE_SIZE)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("couldn't read buffer in loop: %w", err)
		}

		req := &gen.UploadImageRequest{
			ImgChunk: buf[:n],
			Name:     info.Name(),
			Format:   ext,
		}
		if err = stream.Send(req); err != nil {
			return fmt.Errorf("couldn't send chunk: %w", err)
		}
	}
	return nil
}

func (cl *ImageServiceClient) validateFile(filePath string) (os.FileInfo, error) {
	ext := filepath.Ext(filePath)
	if ext != ".jpg" && ext != ".png" {
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read stat info: %w", err)
	}
	if info.Size() > 2<<20 {
		return nil, fmt.Errorf("file too large: %d bytes", info.Size())
	}

	return info, nil
}
