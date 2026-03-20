package twitter

// Client holds connection configuration for the Twitter API.
type Client struct {
	apiKey            string
	apiKeySecret      string
	accessToken       string
	accessTokenSecret string
	bearerToken       string
}

func NewClient(apiKey, apiKeySecret, accessToken, accessTokenSecret, bearerToken string) *Client {
	return &Client{
		apiKey:            apiKey,
		apiKeySecret:      apiKeySecret,
		accessToken:       accessToken,
		accessTokenSecret: accessTokenSecret,
		bearerToken:       bearerToken,
	}
}
