package service

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"blog/internal/database"
	"blog/internal/model"

	"github.com/gosimple/slug"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(html.WithUnsafe()),
)

func renderMarkdown(content string) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func makeSummary(content string) string {
	runes := []rune(content)
	if len(runes) > 200 {
		return string(runes[:200]) + "..."
	}
	return content
}

func uniqueSlug(base string) string {
	s := slug.Make(base)
	if s == "" {
		s = "post"
	}
	candidate := s
	for i := 1; ; i++ {
		var count int64
		database.DB.Model(&model.Post{}).Where("slug = ?", candidate).Count(&count)
		if count == 0 {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", s, i)
	}
}

type CreatePostInput struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	CoverURL   string `json:"cover_url"`
	CategoryID *uint  `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
	Published  bool   `json:"published"`
}

type UpdatePostInput struct {
	Title      *string `json:"title"`
	Content    *string `json:"content"`
	CoverURL   *string `json:"cover_url"`
	CategoryID *uint   `json:"category_id"`
	TagIDs     []uint  `json:"tag_ids"`
	Published  *bool   `json:"published"`
}

type PostFilter struct {
	CategorySlug string
	TagSlug      string
	Search       string
	Page         int
	PageSize     int
}

type AdminPostFilter struct {
	Search    string
	Published string // "true", "false", or "" (all)
	Page      int
	PageSize  int
}

func CreatePost(userID uint, in CreatePostInput) (*model.Post, error) {
	html, err := renderMarkdown(in.Content)
	if err != nil {
		return nil, err
	}

	post := model.Post{
		Title:       in.Title,
		Slug:        uniqueSlug(in.Title),
		Content:     in.Content,
		ContentHTML: html,
		Summary:     makeSummary(strings.TrimSpace(in.Content)),
		CoverURL:    in.CoverURL,
		CategoryID:  in.CategoryID,
		UserID:      userID,
		Published:   in.Published,
	}

	if len(in.TagIDs) > 0 {
		var tags []model.Tag
		if err := database.DB.Find(&tags, in.TagIDs).Error; err != nil {
			return nil, err
		}
		post.Tags = tags
	}

	if err := database.DB.Create(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func GetPost(postSlug string) (*model.Post, error) {
	var post model.Post
	err := database.DB.Preload("Category").Preload("Tags").Preload("User").
		Where("slug = ? AND published = ?", postSlug, true).
		First(&post).Error
	if err != nil {
		return nil, errors.New("post not found")
	}
	return &post, nil
}

func ListPosts(f PostFilter) ([]model.Post, int64, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 50 {
		f.PageSize = 10
	}

	q := database.DB.Model(&model.Post{}).Preload("Category").Preload("Tags").
		Where("published = ?", true)

	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("title LIKE ? OR content LIKE ?", like, like)
	}

	if f.CategorySlug != "" {
		var cat model.Category
		if err := database.DB.Where("slug = ?", f.CategorySlug).First(&cat).Error; err == nil {
			q = q.Where("category_id = ?", cat.ID)
		}
	}

	if f.TagSlug != "" {
		var tag model.Tag
		if err := database.DB.Where("slug = ?", f.TagSlug).First(&tag).Error; err == nil {
			q = q.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
				Where("post_tags.tag_id = ?", tag.ID)
		}
	}

	var total int64
	q.Count(&total)

	var posts []model.Post
	err := q.Order("created_at DESC").
		Offset((f.Page - 1) * f.PageSize).
		Limit(f.PageSize).
		Find(&posts).Error
	return posts, total, err
}

func UpdatePost(postSlug string, in UpdatePostInput) (*model.Post, error) {
	var post model.Post
	if err := database.DB.Preload("Tags").Where("slug = ?", postSlug).First(&post).Error; err != nil {
		return nil, errors.New("post not found")
	}

	if in.Title != nil {
		post.Title = *in.Title
	}
	if in.Content != nil {
		html, err := renderMarkdown(*in.Content)
		if err != nil {
			return nil, err
		}
		post.Content = *in.Content
		post.ContentHTML = html
		post.Summary = makeSummary(strings.TrimSpace(*in.Content))
	}
	if in.CoverURL != nil {
		post.CoverURL = *in.CoverURL
	}
	if in.CategoryID != nil {
		post.CategoryID = in.CategoryID
	}
	if in.Published != nil {
		post.Published = *in.Published
	}

	if in.TagIDs != nil {
		var tags []model.Tag
		if len(in.TagIDs) > 0 {
			if err := database.DB.Find(&tags, in.TagIDs).Error; err != nil {
				return nil, err
			}
		}
		if err := database.DB.Model(&post).Association("Tags").Replace(tags); err != nil {
			return nil, err
		}
	}

	if err := database.DB.Save(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func DeletePost(postSlug string) error {
	result := database.DB.Where("slug = ?", postSlug).Delete(&model.Post{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("post not found")
	}
	return nil
}

func ListAllPosts(f AdminPostFilter) ([]model.Post, int64, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 || f.PageSize > 100 {
		f.PageSize = 20
	}

	q := database.DB.Model(&model.Post{}).Preload("Category").Preload("Tags")

	switch f.Published {
	case "true":
		q = q.Where("published = ?", true)
	case "false":
		q = q.Where("published = ?", false)
	}

	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("title LIKE ? OR content LIKE ?", like, like)
	}

	var total int64
	q.Count(&total)

	var posts []model.Post
	err := q.Order("created_at DESC").
		Offset((f.Page - 1) * f.PageSize).
		Limit(f.PageSize).
		Find(&posts).Error
	return posts, total, err
}

func GetAnyPost(postSlug string) (*model.Post, error) {
	var post model.Post
	err := database.DB.Preload("Category").Preload("Tags").Preload("User").
		Where("slug = ?", postSlug).
		First(&post).Error
	if err != nil {
		return nil, errors.New("post not found")
	}
	return &post, nil
}
