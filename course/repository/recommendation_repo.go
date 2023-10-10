package repository

import (
	"context"
	"course/model"
	"course/pb/pb_rec"
	"time"

	"github.com/gookit/slog"
)

type RecommendationRepository struct {
	rpcClient pb_rec.RecServiceClient
}

func NewRecommendationRepository(client pb_rec.RecServiceClient) *RecommendationRepository {
	slog.Info("Initializing Recommendation Repository")

	return &RecommendationRepository{
		rpcClient: client,
	}
}

func (repo *RecommendationRepository) AddCourse(course model.Course) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cateogyTitles := make([]string, len(course.Categories))
	for i := range cateogyTitles {
		cat := course.Categories[i]
		cateogyTitles[i] = cat.Title
	}

	_, err := repo.rpcClient.AddCourse(ctx, &pb_rec.Course{
		Id: course.ID,
		Title: course.Title,
		Description: course.Description,
		Author: course.AuthorID,
		Categories: cateogyTitles,
	})

	return err
}
