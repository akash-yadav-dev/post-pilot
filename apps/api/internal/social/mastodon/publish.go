package mastodon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"post-pilot/apps/api/internal/social"
	"strings"
)

// Ensure Publisher satisfies social.Publisher at compile time.
var _ social.Publisher = (*Publisher)(nil)

// Publisher is the Mastodon adapter.
type Publisher struct {
	client *Client
}

func NewPublisher(client *Client) *Publisher {
	return &Publisher{client: client}
}

func (p *Publisher) Platform() string { return "mastodon" }

func (p *Publisher) Publish(ctx context.Context, req social.PublishRequest) (*social.PublishResult, error) {
	serverURL := strings.TrimSpace(p.client.serverURL)
	if v, ok := req.Metadata["server_url"].(string); ok && strings.TrimSpace(v) != "" {
		serverURL = strings.TrimSpace(v)
	}
	if serverURL == "" {
		serverURL = "https://mastodon.social"
	}

	form := url.Values{}
	form.Set("status", strings.TrimSpace(req.Content))

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(serverURL, "/")+"/api/v1/statuses", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("mastodon: build request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+strings.TrimSpace(req.AccessToken))
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("mastodon: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("mastodon: publish failed with status %d", resp.StatusCode)
	}

	var out struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("mastodon: decode response: %w", err)
	}

	return &social.PublishResult{ExternalID: strings.TrimSpace(out.ID), URL: strings.TrimSpace(out.URL)}, nil
}
