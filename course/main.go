package main

import (
	api "course/api/rpc"
	"course/service"
	"course/storage"
	"net"

	"github.com/gookit/slog"
)

func main() {
	// init db
	storage.NewDatabase()

	// course service
	courseSvc := service.CourseService{}

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
