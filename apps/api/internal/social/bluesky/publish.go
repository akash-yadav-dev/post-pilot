package bluesky

import (
	"context"
	"fmt"

	"post-pilot/apps/api/internal/social"
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
	// TODO: implement com.atproto.repo.createRecord via AT Protocol.
	return nil, fmt.Errorf("bluesky: Publish not yet implemented")
}
