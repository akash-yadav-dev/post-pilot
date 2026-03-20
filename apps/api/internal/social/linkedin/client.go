package linkedin

// Client holds connection configuration for the LinkedIn API.
type Client struct {
	clientID     string
	clientSecret string
}

func NewClient(clientID, clientSecret string) *Client {
	return &Client{clientID: clientID, clientSecret: clientSecret}
}
