package main

import (
	"context"
	"github.com/freer4an/image-storage/internal/config"
	"github.com/freer4an/image-storage/internal/db"
	"github.com/freer4an/image-storage/internal/repository"
	"github.com/freer4an/image-storage/internal/services"
	"github.com/freer4an/image-storage/protos/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.New("configs.yml")
	grpcServer := grpc.NewServer()
	ctx := context.Background()
	pgxPool := db.ConnectToPostgres(ctx, cfg.GetDbUrl())
	storage := repository.NewImageRepository(pgxPool)
	imageService := services.NewImageServer(cfg.Paths.OImagesStorage, cfg.Paths.ThumbnailsStorage, storage)
	gen.RegisterImageServiceServer(grpcServer, imageService)
	lis, err := net.Listen("tcp", cfg.App.Addr)
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
