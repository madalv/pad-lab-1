package api

import (
	"context"

	"course/model"
	"course/pb"
	"course/util"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gookit/slog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	STATUS_SERVING = "SERVING"
)

type courseService interface {
	GetByID(id string) (model.Course, error)
	Create(course model.Course, categoryIDs []string) (string, error)
	GetAll(pagination util.Pagination) []model.Course
	GetCourseIDsForUser(userID string) []string
}

type chapterService interface {
	GetByID(id string) (model.Chapter, error)
	Create(chapter model.Chapter) (string, error)
}

type grpcServer struct {
	pb.CourseServiceServer
	courseSvc  courseService
	chapterSvc chapterService
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	slog.Printf("Received request for method %s with request: %v", info.FullMethod, req)
	return handler(ctx, req)
}

func NewGRPCServer(courseSvc courseService, chapterSvc chapterService) (*grpc.Server, error) {
	slog.Info("Initializing gRPC Server...")

	s := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	pb.RegisterCourseServiceServer(s, &grpcServer{
		courseSvc:  courseSvc,
		chapterSvc: chapterSvc,
	})

	return s, nil
}

// CreateChapter implements pb.CourseServiceServer.
func (srv *grpcServer) CreateChapter(_ context.Context, req *pb.CreateChapterRequest) (*pb.ChapterId, error) {
	chapter := model.Chapter{
		Title:    req.Title,
		Body:     req.Body,
		CourseID: req.CourseId,
	}

	id, err := srv.chapterSvc.Create(chapter)
	if err != nil {
		return nil, err
	}

	return &pb.ChapterId{Id: id}, nil
}

// GetChapter implements pb.CourseServiceServer.
func (srv *grpcServer) GetChapter(_ context.Context, req *pb.ChapterId) (*pb.Chapter, error) {
	chapter, err := srv.chapterSvc.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Chapter{
		Id:        chapter.ID,
		CourseId:  chapter.CourseID,
		Title:     chapter.Title,
		Body:      chapter.Body,
		CreatedAt: timestamppb.New(chapter.CreatedAt),
		UpdatedAt: timestamppb.New(chapter.UpdatedAt),
	}, nil
}

// CreateCourse implements pb.CourseServiceServer.
func (srv *grpcServer) CreateCourse(_ context.Context, req *pb.CreateCourseRequest) (*pb.CourseId, error) {
	course := model.Course{
		Title:       req.Title,
		Description: req.Description,
		AuthorID:    req.AuthorId,
	}

	id, err := srv.courseSvc.Create(course, req.CategoryIds)
	if err != nil {
		return nil, err
	}

	return &pb.CourseId{Id: id}, nil
}

// GetCourse implements pb.CourseServiceServer.
func (srv *grpcServer) GetCourse(_ context.Context, req *pb.CourseId) (*pb.CourseWithChapters, error) {
	course, err := srv.courseSvc.GetByID(req.Id)
	if err != nil {
		return nil, err
	}

	// transform category list into pb
	pbCategories := make([]*pb.Category, len(course.Categories))
	for i := range pbCategories {
		category := course.Categories[i]
		pbCategories[i] = &pb.Category{
			Id:    category.ID,
			Title: category.Title,
		}
	}

	// transform chapter list into pb
	pbChapters := make([]*pb.ChapterTitle, len(course.Chapters))
	for i := range pbChapters {
		chapter := course.Chapters[i]
		pbChapters[i] = &pb.ChapterTitle{
			Id:    chapter.ID,
			Title: chapter.Title,
		}
	}

	resp := pb.CourseWithChapters{
		Course: &pb.Course{
			Id:          course.ID,
			Description: course.Description,
			Title:       course.Title,
			AuthorId:    course.AuthorID,
			CreatedAt:   timestamppb.New(course.CreatedAt),
			UpdatedAt:   timestamppb.New(course.UpdatedAt),
			Categories:  pbCategories,
		},
		Chapters: pbChapters,
	}

	return &resp, nil
}

// GetCourseIdsForUser implements pb.CourseServiceServer.
func (srv *grpcServer) GetCourseIdsForUser(_ context.Context, req *pb.UserId) (*pb.CourseIds, error) {
	ids := srv.courseSvc.GetCourseIDsForUser(req.Id)
	return &pb.CourseIds{
		Ids: ids,
	}, nil
}

// GetCourses implements pb.CourseServiceServer.
func (srv *grpcServer) GetCourses(_ context.Context, req *pb.PaginationQuery) (*pb.Courses, error) {
	pag := util.Pagination{
		Page:  int(req.Page),
		Limit: int(req.Limit),
	}

	// transform list of courses into pb
	courses := srv.courseSvc.GetAll(pag)
	pbCourses := make([]*pb.Course, len(courses))
	for i := range pbCourses {
		course := courses[i]

		pbCategories := make([]*pb.Category, len(course.Categories))
		for i := range pbCategories {
			category := course.Categories[i]
			pbCategories[i] = &pb.Category{
				Id:    category.ID,
				Title: category.Title,
			}
		}

		pbCourses[i] = &pb.Course{
			Id:          course.ID,
			Description: course.Description,
			Title:       course.Title,
			AuthorId:    course.AuthorID,
			CreatedAt:   timestamppb.New(course.CreatedAt),
			UpdatedAt:   timestamppb.New(course.UpdatedAt),
			Categories:  pbCategories,
		}
	}

	return &pb.Courses{
		Courses: pbCourses,
	}, nil
}

func (srv *grpcServer) GetServerStatus(context.Context, *empty.Empty) (*pb.ServerStatus, error) {
	return &pb.ServerStatus{
		Status: STATUS_SERVING,
	}, nil
}
