package main

import (
	"context"
	"github.com/freer4an/image-storage/internal/connections"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	migrations "github.com/freer4an/image-storage/goose"
	"github.com/freer4an/image-storage/internal/config"
	"github.com/freer4an/image-storage/internal/repository"
	"github.com/freer4an/image-storage/internal/services"
	"github.com/freer4an/image-storage/protos/gen"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.New("configs.yml")

	ctx := context.Background()
	if err := migrations.MakeMigrations(cfg.GetDbUrl()); err != nil {
		log.Fatalf("migration failed", err)
	}
	pgxPool := connections.Postgres(ctx, cfg.GetDbUrl())
	storage := repository.NewImageRepository(pgxPool)
	imageService := services.NewImageServer(cfg.Paths.OImagesStorage, cfg.Paths.ThumbnailsStorage, storage)
	grpcServer := grpc.NewServer()

	gen.RegisterImageServiceServer(grpcServer, imageService)

	lis, err := net.Listen("tcp", cfg.GetAddr())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	reflection.Register(grpcServer)
	log.Println("gRPC server started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
}
