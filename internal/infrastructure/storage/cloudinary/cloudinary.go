// internal/infrastructure/storage/cloudinary/cloudinary.go
package cloudinary

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	storageinfra "github.com/cureerel/gotemplate/internal/infrastructure/storage"
)

type Client struct {
	cloudName    string
	apiKey       string
	apiSecret    string
	uploadPreset string // set on Cloudinary dashboard — "unsigned" preset or leave empty for signed
	http         *http.Client
}

func New(cloudName, apiKey, apiSecret string) storageinfra.Provider {
	return &Client{
		cloudName:    cloudName,
		apiKey:       apiKey,
		apiSecret:    apiSecret,
		uploadPreset: "", // leave blank to use signed uploads
		http:         &http.Client{Timeout: 30 * time.Second},
	}
}

func NewWithPreset(cloudName, apiKey, apiSecret, uploadPreset string) storageinfra.Provider {
	c := New(cloudName, apiKey, apiSecret).(*Client)
	c.uploadPreset = uploadPreset
	return c
}

// Upload sends the file to Cloudinary.
// If uploadPreset is set, uses unsigned upload. Otherwise uses signed upload
// (recommended for server-side use — no preset needed, API key + secret suffice).
func (c *Client) Upload(ctx context.Context, in storageinfra.UploadInput) (storageinfra.UploadResult, error) {
	endpoint := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", c.cloudName)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	// Folder determines the public_id prefix in Cloudinary
	publicID := in.Filename
	if in.Folder != "" {
		publicID = in.Folder + "/" + in.Filename
	}

	if c.uploadPreset != "" {
		// Unsigned upload — simpler, requires a preset configured in the dashboard
		_ = w.WriteField("upload_preset", c.uploadPreset)
		_ = w.WriteField("public_id", publicID)
	} else {
		// Signed upload — more secure, no preset required
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		sig := c.sign(map[string]string{
			"public_id": publicID,
			"timestamp": timestamp,
		})
		_ = w.WriteField("api_key", c.apiKey)
		_ = w.WriteField("timestamp", timestamp)
		_ = w.WriteField("signature", sig)
		_ = w.WriteField("public_id", publicID)
	}

	fw, err := w.CreateFormFile("file", in.Filename)
	if err != nil {
		return storageinfra.UploadResult{}, fmt.Errorf("cloudinary: form file: %w", err)
	}
	if _, err := fw.Write(in.File); err != nil {
		return storageinfra.UploadResult{}, fmt.Errorf("cloudinary: write file: %w", err)
	}
	w.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &buf)
	if err != nil {
		return storageinfra.UploadResult{}, fmt.Errorf("cloudinary: build request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.http.Do(req)
	if err != nil {
		return storageinfra.UploadResult{}, fmt.Errorf("cloudinary: http: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return storageinfra.UploadResult{}, fmt.Errorf("cloudinary: status %d: %s", resp.StatusCode, body)
	}

	var result struct {
		SecureURL string `json:"secure_url"`
		PublicID  string `json:"public_id"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return storageinfra.UploadResult{}, fmt.Errorf("cloudinary: decode: %w", err)
	}
	return storageinfra.UploadResult{URL: result.SecureURL, Key: result.PublicID}, nil
}

// Delete removes an asset by public_id using a signed destroy request.
func (c *Client) Delete(ctx context.Context, key string) error {
	endpoint := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/destroy", c.cloudName)

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sig := c.sign(map[string]string{
		"public_id": key,
		"timestamp": timestamp,
	})

	form := url.Values{}
	form.Set("public_id", key)
	form.Set("api_key", c.apiKey)
	form.Set("timestamp", timestamp)
	form.Set("signature", sig)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint,
		strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("cloudinary: delete request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("cloudinary: delete http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("cloudinary: delete status %d: %s", resp.StatusCode, b)
	}
	return nil
}

// SignedURL returns a plain HTTPS URL for public assets.
// Cloudinary does not require signed URLs for publicly accessible images;
// for private resources you'd need their SDK — this is the standard pattern.
func (c *Client) SignedURL(_ context.Context, key string, _ time.Duration) (string, error) {
	return fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/%s", c.cloudName, key), nil
}

// ── helpers ───────────────────────────────────────────────────

// sign produces a Cloudinary API signature.
// Format: SHA1( "key1=val1&key2=val2" + apiSecret )
// Parameters must be sorted alphabetically — this is required by Cloudinary.
func (c *Client) sign(params map[string]string) string {
	// Build sorted param string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	// Simple insertion sort (small map, no import needed)
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	var parts []string
	for _, k := range keys {
		parts = append(parts, k+"="+params[k])
	}
	toSign := strings.Join(parts, "&") + c.apiSecret

	h := sha1.New() //nolint:gosec — Cloudinary's own spec uses SHA-1
	h.Write([]byte(toSign))
	return fmt.Sprintf("%x", h.Sum(nil))
}