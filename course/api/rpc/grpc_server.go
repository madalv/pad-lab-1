package api

import (
	"context"
	"encoding/json"
	"net/http"

	"course/cache"
	"course/model"
	"course/pb"
	"course/util"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gookit/slog"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
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
	EnrollUser(enr model.Enrollment) error
}

type chapterService interface {
	GetByID(id string) (model.Chapter, error)
	Create(chapter model.Chapter) (string, error)
}

type grpcServer struct {
	pb.CourseServiceServer
	courseSvc  courseService
	chapterSvc chapterService
	cache      cache.Cache
}

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	slog.Printf("Received request for method %s with request: %v", info.FullMethod, req)
	return handler(ctx, req)
}

// newServerMetrics initializes Prometheus metrics with gRPC method interceptors
// and sets up an HTTP endpoint for Prometheus to collect the metrics
func newServerMetrics() *grpcprom.ServerMetrics {
	slog.Info("Initiating server metrics...")
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	go func() {
		if err := http.ListenAndServe(":40052", nil); err != nil {
			slog.Fatal(err)
		}
	}()
	return srvMetrics
}

func NewGRPCServer(courseSvc courseService, chapterSvc chapterService, cache cache.Cache) (*grpc.Server, error) {
	slog.Info("Initializing gRPC Server...")

	srvMetrics := newServerMetrics()
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			srvMetrics.UnaryServerInterceptor(),
		),
	)

	pb.RegisterCourseServiceServer(s, &grpcServer{
		courseSvc:  courseSvc,
		chapterSvc: chapterSvc,
		cache:      cache,
	})

	return s, nil
}

func (srv *grpcServer) EnrollUser(_ context.Context, req *pb.EnrollRequest) (*emptypb.Empty, error) {
	enr := model.Enrollment{
		UserID:   req.UserId,
		CourseID: req.CourseId,
	}

	err := srv.courseSvc.EnrollUser(enr)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
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
	chapter := model.Chapter{}
	jsonVal, err := srv.cache.Get(req.Id)

	if err != nil {
		slog.Error(err)

		chapter, err = srv.chapterSvc.GetByID(req.Id)
		if err != nil {
			return nil, err
		}

		encoded, err := json.Marshal(chapter)
		if err != nil {
			slog.Error(err)
		}
		srv.cache.Set(chapter.ID, string(encoded), 30)
		slog.Info("Set data to cache")
	} else {
		slog.Info("Taking data from cache")
		err = json.Unmarshal([]byte(jsonVal), &chapter)
		if err != nil {
			slog.Error(err)
		}
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
