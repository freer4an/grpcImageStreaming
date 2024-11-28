package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/freer4an/image-storage/internal/config"
	"github.com/freer4an/image-storage/internal/services"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.New("configs.yml")

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())

	cc, err := grpc.NewClient(os.Getenv("APP_ADDR"), opts)
	if err != nil {
		panic(err)
	}
	slog.Info("Client listening on " + os.Getenv("APP_ADDR"))
	client := services.NewImageClient(cc, cfg.App.ImageFormats)
	filePaths, err := getFilesFromPath(cfg.Paths.Images)
	if err != nil {
		panic(err)
	}
	err = client.UploadImage(context.Background(), filePaths)
	if err != nil {
		slog.Error("client.UploadImage failed", slog.String("error", err.Error()))
		return
	}
}

func getFilesFromPath(dir string) ([]string, error) {
	var filePaths []string

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filePaths = append(filePaths, filepath.Join(dir, file.Name()))
		}
	}

	return filePaths, nil
}
