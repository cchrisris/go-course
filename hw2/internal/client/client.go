package client

type Client struct {
	Url string
}

func NewClient(url string) *Client {
	return &Client{url}
}
