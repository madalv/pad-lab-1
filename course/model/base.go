package model

import "time"

type Base struct {
	ID        string    `gorm:"primaryKey;default:gen_random_uuid();"`
	CreatedAt time.Time ``
	UpdatedAt time.Time ``
}
