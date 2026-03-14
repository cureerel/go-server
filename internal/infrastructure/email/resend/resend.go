// internal/infrastructure/email/resend/resend.go
package resend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cureerel/gotemplate/internal/infrastructure/email"
)

const baseURL = "https://api.resend.com"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send implements email.Provider.
// To swap to SES later: implement email.Provider in a new ses/ package,
// return that from your factory, delete nothing here.
func (c *Client) Send(ctx context.Context, e email.Email) error {
	body := map[string]any{
		"from":    e.From,
		"to":      e.To,
		"subject": e.Subject,
		"html":    e.HTML,
	}
	if e.Text != "" {
		body["text"] = e.Text
	}

	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("resend: marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/emails", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("resend: new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("resend: http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errBody map[string]any
		json.NewDecoder(resp.Body).Decode(&errBody)
		return fmt.Errorf("resend: status %d: %v", resp.StatusCode, errBody)
	}

	return nil
}