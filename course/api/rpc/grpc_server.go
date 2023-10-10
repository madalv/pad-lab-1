package api

import (
	"context"

	"course/pb"

	"github.com/gookit/slog"
	"google.golang.org/grpc"
)

type CourseService interface {
}

type grpcServer struct {
	pb.CourseServiceServer
	courseService CourseService
}

func NewGRPCServer(courseSvc CourseService) (*grpc.Server, error) {
	slog.Info("Initializing gRPC Server...")

	s := grpc.NewServer()

	pb.RegisterCourseServiceServer(s, &grpcServer{
		courseService: courseSvc,
	})

	return s, nil
}

// CreateChapter implements pb.CourseServiceServer.
func (*grpcServer) CreateChapter(context.Context, *pb.CreateChapterRequest) (*pb.ChapterId, error) {
	panic("unimplemented")
}

// CreateCourse implements pb.CourseServiceServer.
func (*grpcServer) CreateCourse(context.Context, *pb.CreateCourseRequest) (*pb.CourseId, error) {
	panic("unimplemented")
}

// GetChapter implements pb.CourseServiceServer.
func (*grpcServer) GetChapter(context.Context, *pb.ChapterId) (*pb.Chapter, error) {
	panic("unimplemented")
}

// GetCourse implements pb.CourseServiceServer.
func (*grpcServer) GetCourse(context.Context, *pb.CourseId) (*pb.CourseWithChapters, error) {
	panic("unimplemented")
}

// GetCourseIdsForUser implements pb.CourseServiceServer.
func (*grpcServer) GetCourseIdsForUser(context.Context, *pb.UserId) (*pb.CourseIds, error) {
	panic("unimplemented")
}

// GetCourses implements pb.CourseServiceServer.
func (*grpcServer) GetCourses(context.Context, *pb.PaginationQuery) (*pb.Courses, error) {
	panic("unimplemented")
}
