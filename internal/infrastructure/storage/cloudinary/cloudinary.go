package cloudinary

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	storage "github.com/cureerel/cserver/internal/infrastructure/storage"
)

type Client struct {
	cld *cloudinary.Cloudinary
}

func New(cloudName, apiKey, apiSecret string) (storage.Provider, error) {
	url := fmt.Sprintf("cloudinary://%s:%s@%s", apiKey, apiSecret, cloudName)
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cloudinary: %w", err)
	}
	return &Client{cld: cld}, nil
}

func (c *Client) Upload(ctx context.Context, in storage.UploadInput) (storage.UploadResult, error) {
	params := uploader.UploadParams{
		PublicID:     in.Key,
		Folder:       in.Folder,
		ResourceType: "auto",
	}
	resp, err := c.cld.Upload.Upload(ctx, in.Body, params)
	if err != nil {
		return storage.UploadResult{}, fmt.Errorf("cloudinary upload failed: %w", err)
	}

	// ← debug: print full response to find why URL/Key are empty
	fmt.Printf("[Cloudinary] SecureURL=%q PublicID=%q AssetID=%q Error=%+v\n",
		resp.SecureURL, resp.PublicID, resp.AssetID, resp.Error)

	return storage.UploadResult{
		URL: resp.SecureURL,
		Key: resp.PublicID,
	}, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	resp, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: key})
	if err != nil {
		return fmt.Errorf("cloudinary delete failed: %w", err)
	}
	if resp.Result != "ok" && resp.Result != "not found" {
		return fmt.Errorf("cloudinary destroy result: %s", resp.Result)
	}
	return nil
}

func (c *Client) SignedURL(_ context.Context, key string, _ time.Duration) (string, error) {
	return "", fmt.Errorf("signed_url not implemented")
}