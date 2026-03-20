package linkedin

import (
	"context"
	"fmt"

	"post-pilot/apps/api/internal/social"
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
	// TODO: implement UGC Post creation via LinkedIn API v2.
	return nil, fmt.Errorf("linkedin: Publish not yet implemented")
}
