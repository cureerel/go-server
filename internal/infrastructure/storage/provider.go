// internal/infrastructure/storage/provider.go
package storage

import (
	"context"
	"io"
	"time"
)

// UploadInput is everything needed to upload a file.
type UploadInput struct {
	// Key is the path/filename in the storage bucket.
	// e.g. "blog/covers/my-post.jpg"
	Key string

	// Body is the file content.
	Body io.Reader

	// ContentType e.g. "image/jpeg", "image/webp"
	ContentType string

	// Folder is an optional provider-level folder/prefix.
	// Cloudinary uses this as the upload preset folder.
	Folder string
}

// UploadResult is what every adapter returns after a successful upload.
type UploadResult struct {
	// URL is the public CDN URL of the uploaded file.
	URL string

	// Key is the provider's identifier for the file (use for deletes).
	Key string
}

// Provider is the single interface every storage adapter must satisfy.
// To switch from Cloudinary to S3/R2/GCS — implement this interface
// and swap the wire-up in main.go. Nothing else changes.
type Provider interface {
	// Upload stores the file and returns its public URL + key.
	Upload(ctx context.Context, input UploadInput) (UploadResult, error)

	// Delete removes a file by its key (returned from Upload).
	Delete(ctx context.Context, key string) error

	// SignedURL returns a temporary URL for private files.
	// For public buckets this can just return the permanent URL.
	SignedURL(ctx context.Context, key string, expires time.Duration) (string, error)
}