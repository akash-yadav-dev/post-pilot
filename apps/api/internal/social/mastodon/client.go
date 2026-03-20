package mastodon

// Client holds connection configuration for the Mastodon API.
type Client struct {
	serverURL   string
	accessToken string
}

func NewClient(serverURL, accessToken string) *Client {
	return &Client{serverURL: serverURL, accessToken: accessToken}
}
