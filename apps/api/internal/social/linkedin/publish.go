package linkedin

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

// Publisher is the LinkedIn adapter.
type Publisher struct {
	client *Client
}

func NewPublisher(client *Client) *Publisher {
	return &Publisher{client: client}
}

func (p *Publisher) Platform() string { return "linkedin" }

func (p *Publisher) Publish(ctx context.Context, req social.PublishRequest) (*social.PublishResult, error) {
	authorURN := strings.TrimSpace(req.ExternalAccountID)
	if authorURN == "" {
		return nil, fmt.Errorf("linkedin: missing external account id")
	}
	if !strings.HasPrefix(authorURN, "urn:") {
		authorURN = "urn:li:person:" + authorURN
	}

	payload := map[string]any{
		"author":         authorURN,
		"lifecycleState": "PUBLISHED",
		"specificContent": map[string]any{
			"com.linkedin.ugc.ShareContent": map[string]any{
				"shareCommentary":    map[string]any{"text": req.Content},
				"shareMediaCategory": "NONE",
			},
		},
		"visibility": map[string]any{
			"com.linkedin.ugc.MemberNetworkVisibility": "PUBLIC",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("linkedin: marshal payload: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.linkedin.com/v2/ugcPosts", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("linkedin: build request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+strings.TrimSpace(req.AccessToken))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("linkedin: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("linkedin: publish failed with status %d", resp.StatusCode)
	}

	postURN := strings.TrimSpace(resp.Header.Get("X-RestLi-Id"))
	if postURN == "" {
		postURN = strings.TrimSpace(resp.Header.Get("x-restli-id"))
	}

	return &social.PublishResult{
		ExternalID: postURN,
		URL:        "",
	}, nil
}
