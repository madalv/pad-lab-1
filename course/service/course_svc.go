package service

import (
	"course/model"
	"course/util"
	"fmt"

	"github.com/gookit/slog"
)

type courseRepo interface {
	GetCategoryList(categoryIDs []string) ([]*model.Category, error)
	Create(course *model.Course) error
	GetByIDWithChapters(id string) (model.Course, error)
	GetAll(pagination util.Pagination) []model.Course
	GetCourseIDsForUser(userID string) []string
	CreateEnrollment(enr *model.Enrollment) error
	Delete(id string) error
}

type recRepo interface {
	AddCourse(course model.Course) error
}

type CourseService struct {
	courseRepo courseRepo
	recRepo    recRepo
}

func NewCourseService(cr courseRepo, rr recRepo) *CourseService {
	slog.Info("Initializing Course Service")

	return &CourseService{
		courseRepo: cr,
		recRepo:    rr,
	}
}

func (svc *CourseService) EnrollUser(enr model.Enrollment) error {
	slog.Infof("Enrolling user %s in course %s", enr.UserID, enr.CourseID)

	// check if the course exists first
	_, err := svc.courseRepo.GetByIDWithChapters(enr.CourseID)
	if err != nil {
		slog.Error(err)
		return fmt.Errorf("could not retrieve the course by id %s", enr.CourseID)
	}

	if err := svc.courseRepo.CreateEnrollment(&enr); err != nil {
		slog.Error(err)
		return err
	}

	return nil
}

func (svc *CourseService) GetByID(id string) (model.Course, error) {
	slog.Infof("Getting course by id %s", id)
	course, err := svc.courseRepo.GetByIDWithChapters(id)
	if err != nil {
		slog.Error(err)
		return model.Course{}, err
	}

	return course, nil
}

func (svc *CourseService) Create(course model.Course, categoryIDs []string) (string, error) {
	slog.Infof("Creating course %s", course.Title)
	categories, err := svc.courseRepo.GetCategoryList(categoryIDs)
	if err != nil {
		slog.Error(err)
		return "", err
	}

	course.Categories = categories
	// create course locally
	if err := svc.courseRepo.Create(&course); err != nil {
		slog.Errorf("Course could not be created, aborting: %v", err)
		return "", err
	}

	// create course in the rec svc
	if err := svc.recRepo.AddCourse(course); err != nil {
		slog.Errorf("Course %v creation failed in Rec svc, rolling back: %v", course.ID, err)
		// if course could not be created in rec svc, "rollback" locally
		svc.courseRepo.Delete(course.ID)
		return "", err
	}

	return course.ID, nil
}

func (svc *CourseService) GetAll(pagination util.Pagination) []model.Course {
	slog.Info("Getting all courses")
	return svc.courseRepo.GetAll(pagination)
}

func (svc *CourseService) GetCourseIDsForUser(userID string) []string {
	slog.Infof("Getting a list of course IDs for user %s", userID)
	return svc.courseRepo.GetCourseIDsForUser(userID)
}
