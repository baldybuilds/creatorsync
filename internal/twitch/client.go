package twitch

import (
	"net/http"
	"time"
)

const (
	twitchAPIBaseURL = "https://api.twitch.tv/helix"
	twitchAuthURL    = "https://id.twitch.tv/oauth2/token"
)

type Client struct {
	clientID   string
	httpClient *http.Client
}

func NewClient() (*Client, error) {
	return &Client{
		clientID: "",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}
