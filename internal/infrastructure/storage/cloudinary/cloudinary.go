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

	storage "github.com/cureerel/gotemplate/internal/infrastructure/storage"
)

type Client struct {
	cloudName string
	apiKey    string
	apiSecret string
	http      *http.Client
}

func New(cloudName, apiKey, apiSecret string) storage.Provider {
	return &Client{
		cloudName: cloudName,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		http:      &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) Upload(ctx context.Context, in storage.UploadInput) (storage.UploadResult, error) {
	endpoint := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", c.cloudName)

	fileBytes, err := io.ReadAll(in.Body)
	if err != nil {
		return storage.UploadResult{}, fmt.Errorf("cloudinary: read body: %w", err)
	}

	publicID := in.Key
	if in.Folder != "" && !strings.HasPrefix(in.Key, in.Folder+"/") {
		publicID = in.Folder + "/" + in.Key
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sig := c.sign(map[string]string{
		"public_id": publicID,
		"timestamp": timestamp,
	})

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("api_key", c.apiKey)
	_ = w.WriteField("timestamp", timestamp)
	_ = w.WriteField("signature", sig)
	_ = w.WriteField("public_id", publicID)
	fw, err := w.CreateFormFile("file", in.Key)
	if err != nil {
		return storage.UploadResult{}, fmt.Errorf("cloudinary: form file: %w", err)
	}
	if _, err := fw.Write(fileBytes); err != nil {
		return storage.UploadResult{}, fmt.Errorf("cloudinary: write file: %w", err)
	}
	w.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &buf)
	if err != nil {
		return storage.UploadResult{}, fmt.Errorf("cloudinary: build request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.http.Do(req)
	if err != nil {
		return storage.UploadResult{}, fmt.Errorf("cloudinary: http: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return storage.UploadResult{}, fmt.Errorf("cloudinary: status %d: %s", resp.StatusCode, respBody)
	}

	var result struct {
		SecureURL string `json:"secure_url"`
		PublicID  string `json:"public_id"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return storage.UploadResult{}, fmt.Errorf("cloudinary: decode: %w", err)
	}
	return storage.UploadResult{URL: result.SecureURL, Key: result.PublicID}, nil
}

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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
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

func (c *Client) SignedURL(_ context.Context, key string, _ time.Duration) (string, error) {
	return fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/%s", c.cloudName, key), nil
}

func (c *Client) sign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
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
	h := sha1.New()
	h.Write([]byte(toSign))
	return fmt.Sprintf("%x", h.Sum(nil))
}