package mastodon

import (
	"context"
	"fmt"

	"post-pilot/apps/api/internal/social"
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
	// TODO: implement POST /api/v1/statuses via Mastodon API.
	return nil, fmt.Errorf("mastodon: Publish not yet implemented")
}
