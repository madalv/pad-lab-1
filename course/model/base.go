package model

import "time"

type Base struct {
	ID        string    `gorm:"primaryKey;default:gen_random_uuid();"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;"`
}
