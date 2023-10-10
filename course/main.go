package main

import (
	api "course/api/rpc"
	"course/repository"
	"course/service"
	"course/storage"
	"net"
	"course/pb/pb_rec"

	"github.com/gookit/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// init db
	db := storage.NewDatabase()

	// TODO get address out of config
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Fatalf("Failed to connect to Rec Service: %s", err)
	}
	defer conn.Close()

	// gRPC client for the Recommendation Service
	client := pb_rec.NewRecServiceClient(conn)

	// repositories
	courseRepo := repository.NewCourseRepository(db)
	recRepo := repository.NewRecommendationRepository(client)
	_ = repository.NewChapterRepository(db)

	// services
	courseSvc := service.NewCourseService(courseRepo, recRepo)

	// TODO get address out of config
	listener, err := net.Listen("tcp", "localhost:50052")
	if err != nil {
		slog.Fatal(err)
	}

	srv, err := api.NewGRPCServer(courseSvc)
	if err != nil {
		slog.Fatalf("Could not init gRPC Server: %s", err)
	}

	if err := srv.Serve(listener); err != nil {
		slog.Fatalf("Could not serve RPC: %s", err)
	}
}
