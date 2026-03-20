package social

import "context"

// PublishRequest carries the payload sent to a social platform.
type PublishRequest struct {
	// Content is the text body of the post.
	Content string

	// MediaURLs are optional pre-uploaded media attachments.
	MediaURLs []string

	// ExternalAccountID is the platform-specific account identifier
	// (e.g. Twitter user ID, Mastodon account handle).
	ExternalAccountID string

	// AccessToken is the OAuth access token for the account.
	AccessToken string

	// AccessTokenSecret is required by OAuth 1.0a platforms (e.g. Twitter v1.1).
	AccessTokenSecret string

	// Metadata carries platform-specific context (e.g. mastodon server URL, bluesky handle).
	Metadata map[string]any
}

// PublishResult is the outcome of a successful publish operation.
type PublishResult struct {
	// ExternalID is the platform-assigned ID for the created post.
	ExternalID string

	// URL is the canonical URL of the created post, when available.
	URL string
}

// Publisher is the interface every social-platform adapter must satisfy.
// Each provider (twitter, linkedin, mastodon, bluesky) implements this.
type Publisher interface {
	// Platform returns the canonical platform identifier (e.g. "twitter").
	Platform() string

	// Publish sends a post to the platform and returns the result.
	Publish(ctx context.Context, req PublishRequest) (*PublishResult, error)
}
