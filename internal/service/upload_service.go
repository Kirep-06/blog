package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"blog/config"
	"blog/internal/database"
	"blog/internal/model"
	"blog/internal/storage"

	"github.com/google/uuid"
)

func UploadImage(ctx context.Context, provider storage.StorageProvider, fh *multipart.FileHeader, userID uint) (*model.Image, error) {
	cfg := config.C.Image

	// validate size
	maxBytes := int64(cfg.MaxSizeMB) * 1024 * 1024
	if fh.Size > maxBytes {
		return nil, fmt.Errorf("file too large: max %dMB", cfg.MaxSizeMB)
	}

	// validate MIME type
	mimeType := fh.Header.Get("Content-Type")
	allowed := false
	for _, t := range cfg.AllowedTypes {
		if t == mimeType {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, errors.New("unsupported image type")
	}

	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ext := extensionFromMIME(mimeType)
	key := uuid.NewString() + ext

	url, err := provider.Upload(ctx, key, f.(io.Reader), fh.Size, mimeType)
	if err != nil {
		return nil, err
	}

	img := model.Image{
		Filename: key,
		URL:      url,
		Size:     fh.Size,
		MimeType: mimeType,
		UserID:   userID,
	}
	if err := database.DB.Create(&img).Error; err != nil {
		return nil, err
	}
	return &img, nil
}

func extensionFromMIME(mime string) string {
	switch mime {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}
