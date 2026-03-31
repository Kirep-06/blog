package storage

import (
	"context"
	"io"
)

// StorageProvider abstracts image storage backends.
type StorageProvider interface {
	Upload(ctx context.Context, key string, r io.Reader, size int64, mimeType string) (url string, err error)
	Delete(ctx context.Context, key string) error
}
