package bluesky

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"post-pilot/apps/api/internal/social"
	"strings"
	"time"
)

// Ensure Publisher satisfies social.Publisher at compile time.
var _ social.Publisher = (*Publisher)(nil)

// Publisher is the Bluesky (AT Protocol) adapter.
type Publisher struct {
	client *Client
}

func NewPublisher(client *Client) *Publisher {
	return &Publisher{client: client}
}

func (p *Publisher) Platform() string { return "bluesky" }

func (p *Publisher) Publish(ctx context.Context, req social.PublishRequest) (*social.PublishResult, error) {
	pdsURL := strings.TrimSpace(p.client.pdsURL)
	if v, ok := req.Metadata["pds_url"].(string); ok && strings.TrimSpace(v) != "" {
		pdsURL = strings.TrimSpace(v)
	}
	if pdsURL == "" {
		pdsURL = "https://bsky.social"
	}

	handle := strings.TrimSpace(p.client.handle)
	if v, ok := req.Metadata["handle"].(string); ok && strings.TrimSpace(v) != "" {
		handle = strings.TrimSpace(v)
	}
	if handle == "" {
		return nil, fmt.Errorf("bluesky: missing handle in metadata")
	}

	password := strings.TrimSpace(req.AccessToken)
	if password == "" {
		return nil, fmt.Errorf("bluesky: missing app password")
	}

	sessionReqBody, _ := json.Marshal(map[string]any{
		"identifier": handle,
		"password":   password,
	})
	sessionReq, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(pdsURL, "/")+"/xrpc/com.atproto.server.createSession", bytes.NewReader(sessionReqBody))
	if err != nil {
		return nil, fmt.Errorf("bluesky: build session request: %w", err)
	}
	sessionReq.Header.Set("Content-Type", "application/json")

	sessionResp, err := http.DefaultClient.Do(sessionReq)
	if err != nil {
		return nil, fmt.Errorf("bluesky: create session failed: %w", err)
	}
	defer sessionResp.Body.Close()

	if sessionResp.StatusCode < 200 || sessionResp.StatusCode >= 300 {
		return nil, fmt.Errorf("bluesky: create session failed with status %d", sessionResp.StatusCode)
	}

	var sessionData struct {
		AccessJwt string `json:"accessJwt"`
		DID       string `json:"did"`
	}
	if err := json.NewDecoder(sessionResp.Body).Decode(&sessionData); err != nil {
		return nil, fmt.Errorf("bluesky: decode session response: %w", err)
	}

	createRecordPayload, _ := json.Marshal(map[string]any{
		"repo":       sessionData.DID,
		"collection": "app.bsky.feed.post",
		"record": map[string]any{
			"text":      strings.TrimSpace(req.Content),
			"createdAt": time.Now().UTC().Format(time.RFC3339Nano),
		},
	})

	createReq, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(pdsURL, "/")+"/xrpc/com.atproto.repo.createRecord", bytes.NewReader(createRecordPayload))
	if err != nil {
		return nil, fmt.Errorf("bluesky: build create record request: %w", err)
	}
	createReq.Header.Set("Authorization", "Bearer "+strings.TrimSpace(sessionData.AccessJwt))
	createReq.Header.Set("Content-Type", "application/json")

	createResp, err := http.DefaultClient.Do(createReq)
	if err != nil {
		return nil, fmt.Errorf("bluesky: create post request failed: %w", err)
	}
	defer createResp.Body.Close()

	if createResp.StatusCode < 200 || createResp.StatusCode >= 300 {
		return nil, fmt.Errorf("bluesky: create post failed with status %d", createResp.StatusCode)
	}

	var created struct {
		URI string `json:"uri"`
		CID string `json:"cid"`
	}
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		return nil, fmt.Errorf("bluesky: decode create post response: %w", err)
	}

	postURL := ""
	parts := strings.Split(strings.TrimSpace(created.URI), "/")
	if len(parts) > 0 {
		rkey := parts[len(parts)-1]
		if rkey != "" {
			postURL = "https://bsky.app/profile/" + handle + "/post/" + rkey
		}
	}

	return &social.PublishResult{ExternalID: strings.TrimSpace(created.URI), URL: postURL}, nil
}
