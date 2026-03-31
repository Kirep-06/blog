package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title       string    `gorm:"size:255;not null"             json:"title"`
	Slug        string    `gorm:"uniqueIndex;size:255;not null" json:"slug"`
	Content     string    `gorm:"type:longtext"                 json:"content"`
	ContentHTML string    `gorm:"type:longtext"                 json:"content_html"`
	Summary     string    `gorm:"size:500"                      json:"summary"`
	CoverURL    string    `gorm:"size:500"                      json:"cover_url"`
	Published   bool      `gorm:"default:false"                 json:"published"`
	UserID      uint      `gorm:"not null"                      json:"user_id"`
	User        User      `gorm:"foreignKey:UserID"             json:"user,omitempty"`
	CategoryID  *uint     `                                     json:"category_id"`
	Category    *Category `gorm:"foreignKey:CategoryID"         json:"category,omitempty"`
	Tags        []Tag     `gorm:"many2many:post_tags"           json:"tags"`
}
