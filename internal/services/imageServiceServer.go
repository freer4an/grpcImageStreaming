package services

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/freer4an/image-storage/internal/models"
	"github.com/freer4an/image-storage/protos/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MAX_IMAGE_SIZE = 2 << 20
	WIDTH          = 256
	HEIGHT         = 256
)

type IMCache interface {
	SaveImage(uid string, image []byte)
	GetImage(uid string) ([]byte, bool)
	DeleteImage(uid string)
	ListImages() map[string][]byte
}

type Storage interface {
	SaveImage(ctx context.Context, image *models.Image) (string, error)
	GetImage(ctx context.Context, id string) (*models.Image, error)
	DeleteImage(ctx context.Context, id string) error
	ListImages(ctx context.Context) ([]models.Image, error)
}

type ImageServer struct {
	gen.UnimplementedImageServiceServer
	oImagesStorage    string
	thumbnailsStorage string
	storage           Storage
}

func NewImageServer(oImagesStorage, thumbnailsStorage string, storage Storage) *ImageServer {
	return &ImageServer{
		oImagesStorage:    oImagesStorage,
		thumbnailsStorage: thumbnailsStorage,
		storage:           storage,
	}
}

// Receive client streaming
func (s *ImageServer) UploadImage(stream gen.ImageService_UploadImageServer) error {
	var totalSize uint32
	workerPool := NewWorkerPool()
	ctx, cancel := context.WithTimeout(stream.Context(), time.Second*20)
	defer cancel()
	workerPool.Start(ctx, s.createThumbnail)

	if err := createDir([]string{s.oImagesStorage, s.thumbnailsStorage}); err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	for {
		slog.Debug("Receiving stream")
		req, err := stream.Recv()
		if err == io.EOF {
			slog.Info("All images received")
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive request: %s", err.Error())
		}

		imageID := uuid.New().String()
		imagePath := fmt.Sprintf("%s/%s", s.oImagesStorage, req.GetName())
		thumbnailPath := filepath.Join(s.thumbnailsStorage, imageID, req.GetFormat())

		err = os.WriteFile(imagePath, req.ImgChunk, 0644)
		if err != nil {
			slog.Error("Can't write image",
				slog.String("error", err.Error()))
			return status.Errorf(codes.Internal, "cannot save image %s: %s", req.Name, err.Error())
		}

		slog.Info("Image saved",
			slog.String("receivedName", req.GetName()),
			slog.String("uuid", imageID))

		totalSize += uint32(len(req.GetImgChunk()))

		workerPool.AddTask(Task{
			ImageID:       imageID,
			ImagePath:     imagePath,
			ImageFormat:   req.GetFormat(),
			ThumbnailPath: thumbnailPath,
			Width:         WIDTH,
			Height:        HEIGHT,
		})

	}
	workerPool.Wait()
	res := &gen.UploadImageResponse{
		Size: totalSize,
	}

	err := stream.SendAndClose(res)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot send response: %v", err)
	}

	slog.Info("Upload completed",
		slog.String("size", fmt.Sprintf("%v", res.Size)))

	return nil
}

func (s *ImageServer) GetImage(ctx context.Context, in *gen.GetImageRequest) (*gen.GetImageResponse, error) {
	img, err := s.storage.GetImage(ctx, in.GetId())
	if err != nil {
		slog.Error("Can't get image", slog.String("error", err.Error()))
		return nil, status.Errorf(codes.NotFound, "404")
	}

	res := &gen.GetImageResponse{
		ImageMetada: imageToGenImage(img),
	}

	return res, nil
}

func (s *ImageServer) DeleteImage(ctx context.Context, req *gen.DeleteImageRequest) (*gen.DeleteImageResponse, error) {
	if err := s.storage.DeleteImage(ctx, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "cannot delete image %s: %s", req.GetId(), err.Error())
	}
	return &gen.DeleteImageResponse{}, nil
}

func (s *ImageServer) ListImages(ctx context.Context, req *gen.ListImagesRequest) (*gen.ListImagesResponse, error) {
	images, err := s.storage.ListImages(ctx)
	if err != nil {
		slog.Error("Can't list images", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, err.Error())
	}

	genImages := make([]*gen.Image, 0, len(images))
	for _, img := range images {
		genImages = append(genImages, imageToGenImage(&img))
	}

	res := gen.ListImagesResponse{Images: genImages}
	return &res, nil
}

func (s *ImageServer) createThumbnail(ctx context.Context, task *Task) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("thumbnail creation canceled")
	default:
	}

	src, err := imaging.Open(task.ImagePath)
	if err != nil {
		return fmt.Errorf("cannot open image: %w", err)
	}

	src = imaging.Resize(src, WIDTH, 0, imaging.Lanczos)
	destPath := fmt.Sprintf("%s/%s", s.thumbnailsStorage, task.ImageID+task.ImageFormat)
	if err = imaging.Save(src, destPath); err != nil {
		return fmt.Errorf("cannot save image: %w", err)
	}
	img := taskToImage(task)
	_, err = s.storage.SaveImage(ctx, img)
	if err != nil {
		return fmt.Errorf("cannot save image with id %s: %w", img.Id, err)
	}
	return nil
}

func (s *ImageServer) loadImage(uid string, basePath string) ([]byte, error) {
	filePath := filepath.Join(basePath, uid+".jpg")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read image: %v", err)
	}
	return data, nil
}

func createDir(paths []string) error {
	for _, path := range paths {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			slog.Error("Can't create dir",
				slog.String("path", path),
				slog.String("error", err.Error()))
			return fmt.Errorf("cannot create dir: %w", err)
		}
	}
	return nil
}

func imageToGenImage(img *models.Image) *gen.Image {
	return &gen.Image{
		Id:            img.Id,
		Width:         int32(img.Width),
		Height:        int32(img.Height),
		Format:        img.Format,
		OriginalPath:  img.OriginalPath,
		ThumbnailPath: img.ThumbnailPath,
	}
}

func taskToImage(gen *Task) *models.Image {
	return &models.Image{
		Id:            gen.ImageID,
		Format:        gen.ImageFormat,
		OriginalPath:  gen.ImagePath,
		ThumbnailPath: gen.ThumbnailPath,
		Width:         gen.Width,
		Height:        gen.Height,
	}
}
