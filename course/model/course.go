package model

type Course struct {
	Base
	Title       string      `gorm:"not null;"`
	Description string      `gorm:"not null;"`
	AuthorID    string      `gorm:"not null;"`
	Chapters    []Chapter   `gorm:"foreignKey:CourseID"`
	Categories  []*Category `gorm:"many2many:course_category;"`
}
