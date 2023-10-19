package service

import (
	"course/model"
	"fmt"
	"github.com/gookit/slog"
)

type chapterRepo interface {
	Create(chapter *model.Chapter) error
	GetByID(id string) (model.Chapter, error)
}

type ChapterService struct {
	chapterRepo chapterRepo
	courseRepo  courseRepo
}

func NewChapterService(chr chapterRepo, cr courseRepo) *ChapterService {
	slog.Info("Initializing Course Service")

	return &ChapterService{
		chapterRepo: chr,
		courseRepo:  cr,
	}
}

func (svc *ChapterService) GetByID(id string) (model.Chapter, error) {
	slog.Infof("Getting chapter by id %s", id)
	chapter, err := svc.chapterRepo.GetByID(id)
	if err != nil {
		slog.Error(err)
		return model.Chapter{}, err
	}

	return chapter, nil
}

func (svc *ChapterService) Create(chapter model.Chapter) (string, error) {
	slog.Infof("Creating chapter %s", chapter.Title)

	// check if the course exists first
	_, err := svc.courseRepo.GetByIDWithChapters(chapter.CourseID)
	if err != nil {
		slog.Error(err)
		return "", fmt.Errorf("could not retrieve the course by id %s", chapter.CourseID)
	}

	if err := svc.chapterRepo.Create(&chapter); err != nil {
		slog.Error(err)
		return "", err
	}
	return chapter.ID, nil
}
