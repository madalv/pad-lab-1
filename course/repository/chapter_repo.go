package repository

import (
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
