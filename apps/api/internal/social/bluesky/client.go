package bluesky

// Client holds connection configuration for the Bluesky AT Protocol API.
type Client struct {
	pdsURL   string
	handle   string
	password string
}

func NewClient(pdsURL, handle, password string) *Client {
	return &Client{pdsURL: pdsURL, handle: handle, password: password}
}
