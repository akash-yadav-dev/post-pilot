package twitter

import (
	"context"
	"fmt"

	"post-pilot/apps/api/internal/social"
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
	// TODO: implement OAuth 1.0a tweet creation via Twitter API v2.
	return nil, fmt.Errorf("twitter: Publish not yet implemented")
}
