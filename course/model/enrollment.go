package model

type Enrollment struct {
	UserID   string `gorm:"primaryKey"`
	CourseID string `gorm:"primaryKey"`
	Course   Course `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
