package riot

import (
	"net/http"
	"os"
)

type Client struct {
	APIKey string
	HTTP   *http.Client
}

func NewClient() *Client {
	return &Client{
		APIKey: os.Getenv("RIOT_API_KEY"),
		HTTP:   &http.Client{},
	}
}
