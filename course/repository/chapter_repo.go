package repository

import (
	"course/model"

	"github.com/gookit/slog"
	"gorm.io/gorm"
)

type ChapterRepository struct {
	db *gorm.DB
}

func NewChapterRepository(db *gorm.DB) *ChapterRepository {
	slog.Info("Initializing Chapter Repository")

	return &ChapterRepository{
		db: db,
	}
}

func (r *ChapterRepository) Create(chapter *model.Chapter, categoryIDs []string) error {
	return r.db.Create(chapter).Error
}

func (r *ChapterRepository) GetByID(id string) (model.Chapter, error) {
	var chapter model.Chapter
	err := r.db.First(&chapter, "id = ?", id).Error
	return chapter, err
}
