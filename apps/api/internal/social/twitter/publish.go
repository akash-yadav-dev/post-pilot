package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"post-pilot/apps/api/internal/social"
	"strings"
)

// Ensure Publisher satisfies social.Publisher at compile time.
var _ social.Publisher = (*Publisher)(nil)

// Publisher is the Twitter/X adapter.
type Publisher struct {
	client *Client
}

func NewPublisher(client *Client) *Publisher {
	return &Publisher{client: client}
}

func (p *Publisher) Platform() string { return "twitter" }

func (p *Publisher) Publish(ctx context.Context, req social.PublishRequest) (*social.PublishResult, error) {
	payload := map[string]any{"text": strings.TrimSpace(req.Content)}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("twitter: marshal payload: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.twitter.com/2/tweets", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("twitter: build request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+strings.TrimSpace(req.AccessToken))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("twitter: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("twitter: publish failed with status %d", resp.StatusCode)
	}

	var data struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("twitter: decode response: %w", err)
	}

	result := &social.PublishResult{ExternalID: strings.TrimSpace(data.Data.ID)}
	if result.ExternalID != "" {
		result.URL = "https://twitter.com/i/web/status/" + result.ExternalID
	}

	return result, nil
}
