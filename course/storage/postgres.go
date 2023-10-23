package storage

import (
	"course/model"
	"course/util"

	"github.com/gookit/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(dsn string) *gorm.DB {
	db := connect(dsn)
	err := db.AutoMigrate(&model.Category{})

	if err != nil {
		slog.Fatal(err)
	}
	err = db.AutoMigrate(&model.Chapter{})
	if err != nil {
		slog.Fatal(err)
	}
	err = db.AutoMigrate(&model.Course{})
	if err != nil {
		slog.Fatal(err)
	}
	err = db.AutoMigrate(&model.Enrollment{})
	if err != nil {
		slog.Fatal(err)
	}

	return db
}

func connect(dsn string) *gorm.DB {
	var err error

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		slog.Fatalf("Could not connect to Postgres: %s", err)
	} else {
		slog.Info("Successfully connected to the Postgres")
	}

	return database
}

func Paginate(pagination util.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}