package model

type Chapter struct {
	Base
	Title    string `gorm:"not null;"`
	Body     string `gorm:"not null;"`
	CourseID string `gorm:"not null"`
}
