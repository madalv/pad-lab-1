package model

type Category struct {
	Base
	Title string `gorm:"not null;"`
}
