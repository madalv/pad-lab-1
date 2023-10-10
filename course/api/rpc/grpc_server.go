package api

import (
	"context"
	"time"

	"course/model"
	"course/pb"
	"course/util"

	"github.com/gookit/slog"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type courseService interface {
	GetByID(id string) (model.Course, error)
	Create(course model.Course, categoryIDs []string) (string, error)
	GetAll(pagination util.Pagination) []model.Course
	GetCourseIDsForUser(userID string) []string
}

type grpcServer struct {
	pb.CourseServiceServer
	courseSvc courseService
}

func NewGRPCServer(courseSvc courseService) (*grpc.Server, error) {
	slog.Info("Initializing gRPC Server...")

	s := grpc.NewServer()

	pb.RegisterCourseServiceServer(s, &grpcServer{
		courseSvc: courseSvc,
	})

	return s, nil
}

// CreateChapter implements pb.CourseServiceServer.
func (srv *grpcServer) CreateChapter(context.Context, *pb.CreateChapterRequest) (*pb.ChapterId, error) {
	panic("unimplemented")
}

// CreateCourse implements pb.CourseServiceServer.
func (srv *grpcServer) CreateCourse(_ context.Context, req *pb.CreateCourseRequest) (*pb.CourseId, error) {
	start := time.Now()
	slog.Debug("Time now ", start)
	course := model.Course{
		Title:       req.Title,
		Description: req.Description,
		AuthorID:    req.AuthorId,
	}

	id, err := srv.courseSvc.Create(course, req.CategoryIds)
	if err != nil {
		return nil, err
	}

	elapsed := time.Since(start)
	slog.Printf("Req took %s", elapsed)

	return &pb.CourseId{Id: id}, nil
}

// GetChapter implements pb.CourseServiceServer.
func (*grpcServer) GetChapter(context.Context, *pb.ChapterId) (*pb.Chapter, error) {
	panic("unimplemented")
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
