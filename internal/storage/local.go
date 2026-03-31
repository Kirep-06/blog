package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	uploadDir string
	baseURL   string // e.g. http://localhost:8080/uploads
}

func NewLocalStorage(uploadDir, baseURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("create upload dir: %w", err)
	}
	return &LocalStorage{uploadDir: uploadDir, baseURL: baseURL}, nil
}

func (l *LocalStorage) Upload(_ context.Context, key string, r io.Reader, _ int64, _ string) (string, error) {
	dst := filepath.Join(l.uploadDir, key)
	f, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err = io.Copy(f, r); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}
	return l.baseURL + "/" + key, nil
}

func (l *LocalStorage) Delete(_ context.Context, key string) error {
	return os.Remove(filepath.Join(l.uploadDir, key))
}
