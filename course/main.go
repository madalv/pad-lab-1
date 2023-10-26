package main

import (
	api "course/api/rpc"
	"course/pb/pb_rec"
	"course/repository"
	"course/service"
	"course/storage"
	"net"
	"os"

	"github.com/gookit/slog"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	appMode := os.Getenv("APP_MODE")

	if appMode == "prod" {
		godotenv.Load(".prod.env")
	} else if appMode == "dev" {
		godotenv.Load(".local.env")
	} else {
		slog.Fatalf("Invalid mode: %s. Supported modes are 'dev' and 'prod'.\n", appMode)
	}

	// init db
	db := storage.NewDatabase(os.Getenv("POSTGRES_DSN"))

	conn, err := grpc.Dial(os.Getenv("REC_SVC_ADDRESS"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Fatalf("Failed to connect to Rec Service: %s", err)
	}
	defer conn.Close()

	// gRPC client for the Recommendation Service
	client := pb_rec.NewRecServiceClient(conn)

	// repositories
	courseRepo := repository.NewCourseRepository(db)
	chapterRepo := repository.NewChapterRepository(db)
	recRepo := repository.NewRecommendationRepository(client)
	_ = repository.NewChapterRepository(db)

	// services
	courseSvc := service.NewCourseService(courseRepo, recRepo)
	chapterSvc := service.NewChapterService(chapterRepo, courseRepo)

	listener, err := net.Listen("tcp", os.Getenv("GRPC_ADDRESS"))
	if err != nil {
		slog.Fatal(err)
	}

	srv, err := api.NewGRPCServer(courseSvc, chapterSvc)
	if err != nil {
		slog.Fatalf("Could not init gRPC Server: %s", err)
	}

	if err := srv.Serve(listener); err != nil {
		slog.Fatalf("Could not serve RPC: %s", err)
	}
}
