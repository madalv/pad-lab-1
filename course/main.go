package main

import (
	api "course/api/rpc"
	"course/repository"
	"course/service"
	"course/storage"
	"net"

	"github.com/gookit/slog"
)

func main() {
	// init db
	db := storage.NewDatabase()

	// repositories
	courseRepo := repository.NewCourseRepository(db)
	_ = repository.NewChapterRepository(db)

	// services
	courseSvc := service.NewCourseService(courseRepo)

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
