package model

import "time"

type Image struct {
	ID        uint      `gorm:"primaryKey"      json:"id"`
	Filename  string    `gorm:"size:255;not null" json:"filename"`
	URL       string    `gorm:"size:500;not null" json:"url"`
	Size      int64     `gorm:"not null"         json:"size"`
	MimeType  string    `gorm:"size:64"          json:"mime_type"`
	UserID    uint      `gorm:"not null"         json:"user_id"`
	CreatedAt time.Time `                         json:"created_at"`
}
