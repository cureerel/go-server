// internal/infrastructure/storage/provider.go
package storage

import (
	"context"
	"io"
	"time"
)


type UploadInput struct {

	Key string


	Body io.Reader


	ContentType string

	Folder string
}


type UploadResult struct {

	URL string


	Key string
}


type Provider interface {

	Upload(ctx context.Context, input UploadInput) (UploadResult, error)


	Delete(ctx context.Context, key string) error

	SignedURL(ctx context.Context, key string, expires time.Duration) (string, error)
}